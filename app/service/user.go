package service

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/library/esign/realname"
	"battery/library/esign/realname/beans"
	"battery/library/snowflake"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"math/rand"
	"time"
)

var UserService = userService{}

type userService struct {
}

// CheckRegisterMobile 检测手机号是否可注册
func (s *userService) CheckRegisterMobile(mobile string) bool {
	cnt, err := dao.User.Where(dao.User.Columns.Mobile, mobile).FindCount()
	return cnt == 0 && err == nil
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req model.UserRegisterReq) (uint64, error) {
	if !s.CheckRegisterMobile(req.Mobile) {
		return 0, errors.New("手机号码已经注册，请登录")
	}
	if !SmsServer.Verify(ctx, model.SmsVerifyReq{Mobile: req.Mobile, Code: req.Sms}) {
		return 0, errors.New("验证码错误")
	}
	return s.Create(ctx, req.Mobile, model.UserTypePersonal, 0)
}

// Create 添加用户
func (s *userService) Create(ctx context.Context, mobile string, userType, groupId uint) (uint64, error) {
	salt := dao.User.GenerateSalt()
	id, err := dao.User.Ctx(ctx).InsertAndGetId(g.Map{
		dao.User.Columns.Salt:        salt,
		dao.User.Columns.Type:        userType,
		dao.User.Columns.GroupId:     groupId,
		dao.User.Columns.Mobile:      mobile,
		dao.User.Columns.Qr:          "sgjdriver://" + snowflake.Service().Generate().String(),
		dao.User.Columns.AccessToken: dao.User.GenerateAccessToken(uint64(rand.New(rand.NewSource(time.Now().Unix())).Intn(1000000)), salt),
	})
	return uint64(id), err
}

// CreateGroupUsers 批量添加团签用户
func (*userService) CreateGroupUsers(ctx context.Context, users []model.User) error {
	data := g.List{}
	for key, user := range users {
		salt := dao.User.GenerateSalt()
		data = append(data, g.Map{
			dao.User.Columns.Type:        user.Type,
			dao.User.Columns.RealName:    user.RealName,
			dao.User.Columns.Mobile:      user.Mobile,
			dao.User.Columns.GroupId:     user.GroupId,
			dao.User.Columns.Qr:          snowflake.Service().Generate(),
			dao.User.Columns.AccessToken: dao.User.GenerateAccessToken(uint64(key), salt),
			dao.User.Columns.Salt:        salt,
		})
	}
	_, err := dao.User.Ctx(ctx).Data(data).Insert()
	return err
}

// GroupUserCnt 批量获取团签用户数量
func (*userService) GroupUserCnt(ctx context.Context, groupIds []uint) map[uint]uint {
	var res []struct {
		GroupId uint
		Cnt     uint
	}
	_ = dao.User.Ctx(ctx).
		Fields(dao.User.Columns.GroupId, "count(*) as cnt").
		WhereIn(dao.User.Columns.GroupId, groupIds).
		Group(dao.User.Columns.GroupId).Scan(&res)
	cnt := make(map[uint]uint, len(groupIds))
	for _, groupId := range groupIds {
		cnt[groupId] = 0
	}
	for _, row := range res {
		cnt[row.GroupId] = row.Cnt
	}
	return cnt
}

// SetUsersGroupId 设置用户团体ID
func (*userService) SetUsersGroupId(ctx context.Context, userIds []uint64, groupId uint) error {
	res, err := dao.User.Ctx(ctx).
		WhereIn(dao.User.Columns.BatteryState, []int{model.BatteryStateDefault, model.BatteryStateExit}).
		WhereIn(dao.User.Columns.Id, userIds).Update(g.Map{
		dao.User.Columns.Type:    model.UserTypeGroupRider,
		dao.User.Columns.GroupId: groupId,
	})
	if err != nil {
		return err
	}
	cnt, _err := res.RowsAffected()
	if _err != nil {
		return _err
	}
	if int64(len(userIds)) != cnt {
		return errors.New("团签用户设置失败")
	}
	return nil
}

// SetUserTypeGroupBoss 设置用户为团签管理员
func (*userService) SetUserTypeGroupBoss(ctx context.Context, mobile string, groupId uint) error {
	res, err := dao.User.Ctx(ctx).
		Where(dao.User.Columns.Mobile, mobile).
		WhereIn(dao.User.Columns.BatteryState, []int{model.BatteryStateDefault, model.BatteryStateExit}).
		Where(dao.User.Columns.GroupId, 0).
		Update(g.Map{
			dao.User.Columns.Type:    model.UserTypeGroupBoss,
			dao.User.Columns.GroupId: groupId,
		})
	if err != nil {
		return err
	}
	cnt, _err := res.RowsAffected()
	if _err != nil {
		return _err
	}
	if cnt != 1 {
		return errors.New("团签管理员设置失败")
	}
	return nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req model.UserLoginReq) (rep model.UserLoginRep, err error) {
	var user model.User
	err = dao.User.Ctx(ctx).Where(dao.User.Columns.Mobile, req.Mobile).Scan(&user)
	if err != nil {
		return
	}
	if user.Id == 0 {
		userId, err := s.Register(ctx, model.UserRegisterReq{
			Mobile: req.Mobile,
			Sms:    req.Sms,
		})
		if err != nil {
			return rep, err
		}
		err = dao.User.Ctx(ctx).WherePri(userId).Scan(&user)
		if err != nil {
			return rep, err
		}
		goto GenerateAccessToken
	}
	if !SmsServer.Verify(ctx, model.SmsVerifyReq{
		Mobile: req.Mobile,
		Code:   req.Sms,
	}) {
		err = errors.New("手机号或验证码错误登录失败")
		return
	}
GenerateAccessToken:
	token := dao.User.GenerateAccessToken(user.Id, user.Salt)
	_, err = dao.User.Where(dao.User.Columns.Id, user.Id).Update(g.Map{
		dao.User.Columns.AccessToken: token,
		dao.User.Columns.Qr:          snowflake.Service().Generate(),
	})
	if err == nil {
		rep.AccessToken = token
		rep.Type = user.Type
		rep.AuthState = user.AuthState
	}
	return
}

// RealNameAuthSubmit 骑手实名认证提交
func (s *userService) RealNameAuthSubmit(ctx context.Context, req model.UserRealNameAuthReq) (rep model.UserRealNameAuthRep, err error) {
	user := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
	accountId := user.EsignAccountId
	if accountId == "" {
		idType := "CRED_PSN_CH_IDCARD"
		if req.IdType != "" {
			idType = req.IdType
		}
		res, err := realname.Service().CreatePersonByThirdPartyUserId(beans.CreatePersonByThirdPartyUserIdInfo{
			ThirdPartyUserId: fmt.Sprintf("%s-%d", req.IdCardNo, time.Now().UnixNano()),
			Mobile:           user.Mobile,
			Name:             req.RealName,
			IdNumber:         req.IdCardNo,
			IdType:           idType,
		})
		if err != nil || res.Code != 0 {
			return rep, err
		}
		accountId = res.Data.AccountId
	}
	_, err = dao.User.Ctx(ctx).Where(dao.User.Columns.Id, user.Id).Update(g.Map{
		dao.User.Columns.RealName:       req.RealName,
		dao.User.Columns.IdCardNo:       req.IdCardNo,
		dao.User.Columns.IdCardImg1:     req.IdCardImg1,
		dao.User.Columns.IdCardImg2:     req.IdCardImg2,
		dao.User.Columns.IdCardImg3:     req.IdCardImg3,
		dao.User.Columns.AuthState:      model.AuthStateVerifyWait,
		dao.User.Columns.EsignAccountId: accountId,
	})
	if err != nil {
		return rep, err
	}

	resWeb, err := realname.Service().WebIndivIdentityUrl(beans.WebIndivIdentityUrlInfo{
		AuthType: "PSN_FACEAUTH_BYURL",
		ContextInfo: beans.ContextInfo{
			NotifyUrl: g.Cfg().GetString("app.host") + "/esign/callback/real_name",
		},
		ConfigParams: beans.ConfigParams{
			IndivUneditableInfo: []string{"name", "certNo", "mobileNo"},
		},
	}, accountId)
	if err == nil && resWeb.Code == 0 {
		rep.FlowId = resWeb.Data.FlowId
		rep.ShortLink = resWeb.Data.ShortLink
		rep.Url = resWeb.Data.Url
	}
	return
}

// RealNameAuthVerifyCallBack 骑手实名认证回调通知结果
func (s *userService) RealNameAuthVerifyCallBack(ctx context.Context, eSignAccountId string, req model.RealNameAuthVerifyReq) error {
	_, err := dao.User.Ctx(ctx).Where(dao.User.Columns.EsignAccountId, eSignAccountId).Where(dao.User.Columns.AuthState, req.AuthState).Update(g.Map{
		dao.User.Columns.AuthState: req.AuthState,
	})
	return err
}

// Profile 用户信息
func (s *userService) Profile(ctx context.Context) (rep model.UserProfileRep) {
	u := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
	var user model.User
	_ = dao.User.WherePri(u.Id).Scan(&user)
	rep.Name = user.RealName
	rep.Mobile = user.Mobile
	rep.Type = user.Type
	rep.Qr = u.Qr
	rep.AuthState = user.AuthState
	if user.Type == model.UserTypePersonal {
		packages, _ := PackagesService.Detail(ctx, user.PackagesId)
		rep.User.PackagesId = user.PackagesId
		rep.User.PackagesName = packages.Name
		rep.User.BatteryState = user.BatteryState
		rep.User.BatteryReturnAt = user.BatteryReturnAt
	}
	if user.Type == model.UserTypeGroupRider {
		rep.GroupUser.BatteryState = user.BatteryState
		rep.GroupUser.BatteryType = u.BatteryType
	}
	if user.Type == model.UserTypeGroupBoss {
		rep.GroupBoos.UserCnt = 0 //TODO
		rep.GroupBoos.Days = 0    //TODO
	}
	return
}

// PushToken 用户修改推送token
func (*userService) PushToken(ctx context.Context, req model.PushTokenReq) error {
	u := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
	_, err := dao.User.WherePri(u.Id).Update(g.Map{
		dao.User.Columns.DeviceType:  req.DeviceType,
		dao.User.Columns.DeviceToken: req.DeviceToken,
	})
	return err
}

// MyPackage 用户获取当前套餐信息
func (*userService) MyPackage(ctx context.Context) (rep model.UserCurrentPackageOrder, err error) {
	u := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
	if u.PackagesId > 0 {
		order, err := PackagesOrderService.Detail(ctx, u.PackagesOrderId)
		if err == nil {
			rep.Amount = order.Amount
			rep.Earnest = order.Earnest
			rep.OrderNo = order.No
			rep.PayAt = order.PayAt
			rep.PayType = order.PayType
			rep.StartUseAt = order.FirstUserAt
		}
		packages, err := PackagesService.Detail(ctx, u.PackagesId)
		if err == nil {
			rep.PackageName = packages.Name
			city, _ := DistrictsService.Detail(ctx, packages.CityId)
			rep.CityName = city.Name
		}
		rep.ExpirationAt = u.ExpirationAt
	} else {
		err = errors.New("还未购买套餐")
	}
	return
}

//BizProfile 用户办理业务获取用户信息
func (s *userService) BizProfile(ctx context.Context, qr string) model.BizProfileRep {
	user := s.DetailByQr(ctx, qr)
	rep := model.BizProfileRep{
		Id:           user.Id,
		RealName:     user.RealName,
		Mobile:       user.Mobile,
		IdCardNo:     user.IdCardNo,
		AuthState:    user.AuthState,
		BatteryState: user.BatteryState,
		BatteryType:  user.BatteryType,
		GroupId:      user.GroupId,
	}
	if user.GroupId > 0 {
		group := GroupService.Detail(ctx, user.GroupId)
		rep.GroupName = group.Name
	} else {
		if user.BatteryState == model.BatteryStateUse {
			if user.BatteryReturnAt.Timestamp() < gtime.Now().Timestamp() {
				rep.BatteryState = model.BatteryStateOverdue
			}
		}
	}
	return rep
}

// BizBatterySave 用户寄存电池
func (s *userService) BizBatterySave(ctx context.Context, user model.User) error {
	var seconds int64 = 0
	if !user.BizBatterySecondsStartAt.IsZero() {
		seconds = gtime.Now().Timestamp() - user.BizBatterySecondsStartAt.Timestamp()
	}
	res, err := dao.User.Ctx(ctx).WherePri(user.Id).
		Where(dao.User.Columns.GroupId, 0).
		Where(dao.User.Columns.BatteryState, model.BatteryStateUse).
		Update(g.Map{
			dao.User.Columns.BatteryState:             model.BatteryStateSave,
			dao.User.Columns.BizBatterySecondsStartAt: nil,
			dao.User.Columns.BizBatteryRenewalSeconds: user.BizBatteryRenewalSeconds + uint(seconds),
		})
	if err == nil {
		if rows, err := res.RowsAffected(); rows > 0 && err == nil {
			return nil
		}
		err = errors.New("寄存失败")
	}
	return err
}

// BizBatteryUnSave 用户恢复计费
func (s *userService) BizBatteryUnSave(ctx context.Context, user model.User) error {
	now := gtime.Now()
	res, err := dao.User.Ctx(ctx).WherePri(user.Id).
		Where(dao.User.Columns.GroupId, 0).
		Where(dao.User.Columns.BatteryState, model.BatteryStateSave).
		Update(g.Map{
			dao.User.Columns.BatteryReturnAt:          user.BatteryReturnAt.Add(time.Duration(now.Timestamp()-user.BatteryReturnAt.Timestamp()) * time.Second),
			dao.User.Columns.BatteryState:             model.BatteryStateUse,
			dao.User.Columns.BizBatterySecondsStartAt: now,
		})
	if err == nil {
		if rows, err := res.RowsAffected(); rows > 0 && err == nil {
			return nil
		}
		err = errors.New("寄存失败")
	}
	return err
}

// BizBatteryExit 用户退租
func (s *userService) BizBatteryExit(ctx context.Context, user model.User) error {
	var seconds int64 = 0
	if !user.BizBatterySecondsStartAt.IsZero() {
		seconds = gtime.Now().Timestamp() - user.BizBatterySecondsStartAt.Timestamp()
	}
	_, err := dao.User.Ctx(ctx).WherePri(user.Id).
		WhereIn(dao.User.Columns.BatteryState, []int{model.BatteryStateUse, model.BatteryStateSave}).
		Update(g.Map{
			dao.User.Columns.BatteryState:             model.BatteryStateExit,
			dao.User.Columns.BizBatterySecondsStartAt: nil,
			dao.User.Columns.BizBatteryRenewalSeconds: user.BizBatteryRenewalSeconds + uint(seconds),
		})
	return err
}

// BuyPackagesSuccess 用户成功新购套餐
func (*userService) BuyPackagesSuccess(ctx context.Context, order model.PackagesOrder) error {
	packages, err := PackagesService.Detail(ctx, order.PackageId)
	if err != nil {
		return err
	}
	_, err = dao.User.Ctx(ctx).WherePri(order.UserId).Update(g.Map{
		dao.User.Columns.PackagesOrderId: order.Id,
		dao.User.Columns.PackagesId:      order.PackageId,
		dao.User.Columns.BatteryType:     packages.BatteryType,
		dao.User.Columns.BatteryState:    model.BatteryStateNew,
	})
	return err
}

// RenewalPackagesSuccess 用户成功续购套餐
func (*userService) RenewalPackagesSuccess(ctx context.Context, order model.PackagesOrder) error {
	packages, err := PackagesService.Detail(ctx, order.PackageId)
	if err != nil {
		return err
	}
	var user model.User
	err = dao.User.Ctx(ctx).WherePri(order.UserId).Scan(&user)
	if err == nil {
		_, err = dao.User.Ctx(ctx).WherePri(order.UserId).Update(g.Map{
			dao.User.Columns.BatteryReturnAt: user.BatteryReturnAt.Add(time.Duration(packages.Days) * 24 * time.Hour),
		})
	}
	return err
}

// PenaltyPackagesSuccess 用户支付违约金
func (*userService) PenaltyPackagesSuccess(ctx context.Context, order model.PackagesOrder) error {
	var user model.User
	err := dao.User.Ctx(ctx).WherePri(order.UserId).Scan(&user)
	if err == nil {
		_, err = dao.User.Ctx(ctx).WherePri(order.UserId).Update(g.Map{
			dao.User.Columns.BatteryReturnAt: gtime.Now(),
		})
	}
	return err
}

// GetUserByAccessToken 使用accessToken获取用户信息
func (s *userService) GetUserByAccessToken(accessToken string) (user model.User, err error) {
	err = dao.User.Where(dao.User.Columns.AccessToken, accessToken).Scan(&user)
	if err == nil && user.Id > 0 {
		if user.BatteryState == model.BatteryStateUse {
			if user.GroupId == 0 && user.BatteryReturnAt.Timestamp() <= gtime.Now().Timestamp() {
				user.BatteryState = model.BatteryStateOverdue
			}
		}
	}
	return
}

// GroupUserSelectBattery 团签用户选择电池型号
func (s *userService) GroupUserSelectBattery(ctx context.Context, batteryType uint) error {
	user := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
	res, err := dao.User.WherePri(user.Id).
		WhereGT(dao.User.Columns.GroupId, 0).
		Where(dao.User.Columns.BatteryState, model.BatteryStateDefault).
		Update(g.Map{
			dao.User.Columns.BatteryState: model.BatteryStateNew,
			dao.User.Columns.BatteryType:  batteryType,
		})
	if err != nil {
		return err
	}
	if cnt, _err := res.RowsAffected(); _err == nil && cnt > 0 {
		return nil
	}
	return errors.New("选择失败")
}

// IncrBizBatteryRenewalCnt 增加用户换电次数
func (s *userService) IncrBizBatteryRenewalCnt(ctx context.Context, userId uint64, cnt uint) error {
	_, err := dao.User.Ctx(ctx).WherePri(userId).Increment(dao.User.Columns.BizBatteryRenewalCnt, cnt)
	return err
}

// IncrBizBatteryRenewalSeconds 增加用户使用电池时间
func (s *userService) IncrBizBatteryRenewalSeconds(ctx context.Context, userId uint64, seconds uint) error {
	_, err := dao.User.Ctx(ctx).WherePri(userId).Increment(dao.User.Columns.BizBatteryRenewalSeconds, seconds)
	return err
}

// GetByIds 使用ID获取用户
func (s *userService) GetByIds(ctx context.Context, userIds []uint64) (res []model.User) {
	_ = dao.User.Ctx(ctx).WhereIn(dao.User.Columns.Id, userIds).Scan(&res)
	return res
}

// GetByMobiles 使用mobile获取用户
func (s *userService) GetByMobiles(ctx context.Context, mobiles []string) (res []model.User) {
	_ = dao.User.Ctx(ctx).WhereIn(dao.User.Columns.Mobile, mobiles).Scan(&res)
	return res
}

// Detail 使用ID获取用户信息
func (s *userService) Detail(ctx context.Context, userId uint64) (res model.User) {
	_ = dao.User.Ctx(ctx).WherePri(userId).Scan(&res)
	return res
}

// DetailByQr 使用Qr获取用户信息
func (s *userService) DetailByQr(ctx context.Context, qr string) (res model.User) {
	_ = dao.User.Ctx(ctx).Where(dao.User.Columns.Qr, qr).Scan(&res)
	return res
}

// PackagesStartUse 个签用户新购套餐首次领取
func (*userService) PackagesStartUse(ctx context.Context, order model.PackagesOrder) error {
	packages, err := PackagesService.Detail(ctx, order.PackageId)
	if err != nil {
		return err
	}
	now := gtime.Now()
	res, err := dao.User.Ctx(ctx).WherePri(order.UserId).
		Where(dao.User.Columns.BatteryState, model.BatteryStateNew).
		Update(g.Map{
			dao.User.Columns.BatteryReturnAt:          now.Add(time.Duration(packages.Days) * 24 * time.Hour),
			dao.User.Columns.BatteryState:             model.BatteryStateUse,
			dao.User.Columns.BizBatterySecondsStartAt: now,
		})
	if err == nil {
		if rows, err := res.RowsAffected(); rows > 0 && err == nil {
			return nil
		} else {
			err = errors.New("用户领取失败")
		}
	}
	return err
}

// GroupUserStartUse 团签用户首次领取电池
func (*userService) GroupUserStartUse(ctx context.Context, userId uint64) error {

	res, err := dao.User.Ctx(ctx).WherePri(userId).
		Where(dao.User.Columns.BatteryState, model.BatteryStateNew).
		Update(g.Map{
			dao.User.Columns.BatteryState:             model.BatteryStateUse,
			dao.User.Columns.BizBatterySecondsStartAt: time.Now(),
		})
	if err == nil {
		if rows, err := res.RowsAffected(); rows > 0 && err == nil {
			return nil
		} else {
			err = errors.New("用户领取失败")
		}
	}
	return err
}
