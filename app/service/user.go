package service

import (
    "battery/app/model/user"
    "battery/library/mq"
    "context"
    "errors"
    "fmt"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/os/gtime"
    "github.com/golang-module/carbon"
    "math"
    "math/rand"
    "time"

    "battery/app/dao"
    "battery/app/model"
    "battery/library/esign/realname"
    "battery/library/esign/realname/beans"
    "battery/library/snowflake"
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
    c := dao.User.Columns
    salt := dao.User.GenerateSalt()
    id, err := dao.User.Ctx(ctx).InsertAndGetId(g.Map{
        c.Salt:        salt,
        c.Type:        userType,
        c.GroupId:     groupId,
        c.Mobile:      mobile,
        c.Qr:          snowflake.Service().Generate().String(),
        c.AccessToken: dao.User.GenerateAccessToken(uint64(rand.New(rand.NewSource(time.Now().Unix())).Intn(1000000)), salt),
    })
    return uint64(id), err
}

// AddGroupUsers 批量添加团签用户
func (*userService) AddGroupUsers(ctx context.Context, users []model.User) error {
    c := dao.User.Columns
    for k, _ := range users {
        users[k].Qr = snowflake.Service().Generate().String()
        users[k].Salt = dao.User.GenerateSalt()
    }
    _, err := dao.User.Ctx(ctx).Data(users).OnDuplicateEx(c.Salt, c.AccessToken, c.Qr).Save()
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
        dao.User.Columns.Type:    model.UserTypeGroupMember,
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

// AddOrSetGroupUser 添加或更新用户为团签管理员
func (*userService) AddOrSetGroupUser(ctx context.Context, user model.User) error {
    c := dao.User.Columns
    // 查询是否有冲突用户
    if cnt, err := dao.User.Ctx(ctx).
        Where(c.Mobile, user.Mobile).
        Where(
            fmt.Sprintf("%s NOT IN (?) OR %s > 0", c.BatteryState, c.GroupId),
            g.Slice{model.BatteryStateDefault, model.BatteryStateExit}).
        Count(); err != nil {
        return err
    } else if cnt > 0 {
        return errors.New("已有用户冲突")
    }

    // 新增或修改为管理员
    user.Qr = snowflake.Service().Generate().String()
    user.Salt = dao.User.GenerateSalt()
    _, err := dao.User.Ctx(ctx).Data(user).OnDuplicateEx(c.Salt, c.AccessToken, c.Qr).Save()
    return err
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
    if cnt, err := res.RowsAffected(); err != nil {
        return err
    } else if cnt != 1 {
        return errors.New("团签管理员设置失败")
    }
    return nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req model.UserLoginReq) (rep model.UserLoginRep, err error) {
    var user model.User
    err = dao.User.Ctx(ctx).Where(dao.User.Columns.Mobile, req.Mobile).Scan(&user)
    if err != nil || user.Id == 0 {
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
    } else {
        if !SmsServer.Verify(ctx, model.SmsVerifyReq{
            Mobile: req.Mobile,
            Code:   req.Sms,
        }) {
            err = errors.New("手机号或验证码错误登录失败")
            return
        }
    }
    token := dao.User.GenerateAccessToken(user.Id, user.Salt)
    _, err = dao.User.Where(dao.User.Columns.Id, user.Id).Update(g.Map{
        dao.User.Columns.AccessToken: token,
    })
    if err == nil {
        rep.AccessToken = token
        rep.Type = user.Type
        rep.AuthState = user.AuthState
    }
    return
}

// GetUserByIdCardNo 使用证件号码获取用户
func (s *userService) GetUserByIdCardNo(ctx context.Context, idCardNo string) (model.User, error) {
    var user model.User
    err := dao.User.Ctx(ctx).Where(dao.User.Columns.IdCardNo, idCardNo).Limit(1).Scan(&user)
    return user, err
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
        dao.User.Columns.AuthState:      model.AuthStateVerifyPending,
        dao.User.Columns.EsignAccountId: accountId,
    })
    if err != nil {
        return rep, err
    }

    resWeb, err := realname.Service().WebIndivIdentityUrl(beans.WebIndivIdentityUrlInfo{
        AuthType:           "PSN_FACEAUTH_BYURL",
        AvailableAuthTypes: []string{"PSN_FACEAUTH_BYURL"},
        ContextInfo: beans.ContextInfo{
            NotifyUrl:   g.Cfg().GetString("api.host") + "/esign/callback/real_name",
            RedirectUrl: "sgjdriver://driverapp.shiguangjv.com?path=/webview&from=/verify&mode=off&data=success",
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
    _, err := dao.User.Ctx(ctx).Where(dao.User.Columns.EsignAccountId, eSignAccountId).
        Where(dao.User.Columns.AuthState, model.AuthStateVerifyPending).
        Update(g.Map{
            dao.User.Columns.AuthState: req.AuthState,
        })
    return err
}

// Profile 用户信息
func (s *userService) Profile(ctx context.Context) (rep model.UserProfileRep) {
    u := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
    var mu model.User
    _ = dao.User.WherePri(u.Id).Scan(&mu)
    rep.Name = mu.RealName
    rep.Mobile = mu.Mobile
    rep.Type = mu.Type
    rep.Qr = u.Qr
    rep.AuthState = mu.AuthState
    if mu.Type == model.UserTypePersonal {
        combo, _ := ComboService.Detail(ctx, mu.ComboId)
        rep.User.ComboId = mu.ComboId
        rep.User.ComboName = combo.Name
        rep.User.BatteryState = mu.BatteryState
        rep.User.BatteryReturnAt = mu.BatteryReturnAt
        // 违约
        if mu.BatteryState == model.BatteryStateUse {
            if mu.BatteryReturnAt.Timestamp() < gtime.Now().Timestamp() {
                rep.User.BatteryState = model.BatteryStateOverdue
            }
        }
        // 过期
        if mu.BatteryState == model.BatteryStateSave {
            if mu.BatteryReturnAt.Timestamp() < gtime.Now().Timestamp() {
                rep.User.BatteryState = model.BatteryStateExpired
            }
        }
    }
    if mu.Type == model.UserTypeGroupMember {
        rep.GroupUser.BatteryState = mu.BatteryState
        rep.GroupUser.BatteryType = u.BatteryType
    }
    if mu.Type == model.UserTypeGroupBoss {
        rep.GroupBoss.BillDays, _ = GroupSettlementDetailService.GetDays(ctx, mu.GroupId)
        rep.GroupBoss.MemberCnt = GroupUserService.UserCnt(ctx, mu.GroupId)
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

// MyCombo 用户获取当前套餐信息
func (*userService) MyCombo(ctx context.Context) (rep model.UserCurrentComboOrder, err error) {
    u := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
    if u.ComboId > 0 {
        order, err := ComboOrderService.Detail(ctx, u.ComboOrderId)
        if err == nil {
            rep.Amount = order.Amount
            rep.Deposit = order.Deposit
            rep.OrderNo = order.No
            rep.PayAt = order.PayAt
            rep.PayType = order.PayType
            rep.StartUseAt = order.FirstUseAt
        }
        combo, err := ComboService.Detail(ctx, u.ComboId)
        if err == nil {
            rep.ComboName = combo.Name
            rep.ComboAmount = combo.Price
            city, _ := DistrictsService.Detail(ctx, combo.CityId)
            rep.CityName = city.Name
        }
        rep.ExpirationAt = u.BatteryReturnAt
    } else {
        err = errors.New("还未购买套餐")
    }
    return
}

// BizProfile 用户办理业务获取用户信息
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
        if user.ComboId > 0 {
            if combo, err := ComboService.Detail(ctx, user.ComboId); err == nil {
                rep.ComboName = combo.Name
            }
        }
        // 已逾期
        if user.BatteryState == model.BatteryStateUse {
            if user.BatteryReturnAt.Timestamp() < gtime.Now().Timestamp() {
                rep.BatteryState = model.BatteryStateOverdue
            }
        }
        // 已过期
        if user.BatteryState == model.BatteryStateSave {
            if user.BatteryReturnAt.Timestamp() < gtime.Now().Timestamp() {
                rep.BatteryState = model.BatteryStateExpired
            }
        }
    }
    return rep
}

// BizBatterySave 个签用户寄存电池
func (s *userService) BizBatterySave(ctx context.Context, user model.User) error {
    days := carbon.Parse(user.BizBatteryRenewalDaysStartAt.String()).DiffInDays(carbon.Parse(gtime.Now().String()))
    res, err := dao.User.Ctx(ctx).WherePri(user.Id).
        Where(dao.User.Columns.GroupId, 0).
        Where(dao.User.Columns.BatteryState, model.BatteryStateUse).
        Update(g.Map{
            dao.User.Columns.BatteryState:                 model.BatteryStateSave,
            dao.User.Columns.BatterySaveAt:                gtime.Now(),
            dao.User.Columns.BizBatteryRenewalDays:        user.BizBatteryRenewalDays + uint(days),
            dao.User.Columns.BizBatteryRenewalDaysStartAt: nil,
        })
    if err == nil {
        if rows, err := res.RowsAffected(); rows > 0 && err == nil {
            return nil
        }
        err = errors.New("寄存失败")
    }
    return err
}

// BizBatteryUnSave 个签用户恢复计费
func (s *userService) BizBatteryUnSave(ctx context.Context, user model.User) error {
    now := gtime.Now()
    days := carbon.Parse(user.BatterySaveAt.String()).DiffInDays(carbon.Parse(now.String()))
    var returnAt *gtime.Time
    if days == 0 {
        // 同一天寄取，归还电池时间不变
        returnAt = user.BatteryReturnAt
    } else {
        returnAt = user.BatteryReturnAt.Add(time.Hour * 24 * time.Duration(days))
    }
    res, err := dao.User.Ctx(ctx).WherePri(user.Id).
        Where(dao.User.Columns.GroupId, 0).
        Where(dao.User.Columns.BatteryState, model.BatteryStateSave).
        Update(g.Map{
            dao.User.Columns.BatteryReturnAt:              returnAt,
            dao.User.Columns.BatteryState:                 model.BatteryStateUse,
            dao.User.Columns.BatterySaveAt:                nil,
            dao.User.Columns.BizBatteryRenewalCnt:         user.BizBatteryRenewalCnt + 1,
            dao.User.Columns.BizBatteryRenewalDaysStartAt: gtime.Now(),
        })
    if err == nil {
        if rows, err := res.RowsAffected(); rows > 0 && err == nil {
            return nil
        }
        err = errors.New("恢复计费失败")
    }
    return err
}

// BizBatteryExit 用户退租
func (s *userService) BizBatteryExit(ctx context.Context, user model.User) error {
    var days int64
    if !user.BizBatteryRenewalDaysStartAt.IsZero() {
        days = carbon.Parse(user.BizBatteryRenewalDaysStartAt.String()).DiffInDays(carbon.Parse(gtime.Now().String()))
    }
    _, err := dao.User.Ctx(ctx).WherePri(user.Id).
        WhereIn(dao.User.Columns.BatteryState, []int{model.BatteryStateUse, model.BatteryStateSave}).
        Update(g.Map{
            dao.User.Columns.BatteryState:                 model.BatteryStateExit,
            dao.User.Columns.BatterySaveAt:                nil,
            dao.User.Columns.BizBatteryRenewalDaysStartAt: nil,
            dao.User.Columns.BizBatteryRenewalDays:        user.BizBatteryRenewalDays + uint(days),
        })
    return err
}

// BuyComboSuccess 用户成功新购套餐
func (*userService) BuyComboSuccess(ctx context.Context, order model.ComboOrder) error {
    combo, err := ComboService.Detail(ctx, order.ComboId)
    if err != nil {
        return err
    }
    _, err = dao.User.Ctx(ctx).WherePri(order.UserId).Update(g.Map{
        dao.User.Columns.ComboOrderId:    order.Id,
        dao.User.Columns.ComboId:         order.ComboId,
        dao.User.Columns.BatteryType:     combo.BatteryType,
        dao.User.Columns.BatteryState:    model.BatteryStateNew,
        dao.User.Columns.BatteryReturnAt: nil,
    })
    return err
}

// RenewalComboSuccess 用户成功续购套餐
func (*userService) RenewalComboSuccess(ctx context.Context, order model.ComboOrder) error {
    combo, err := ComboService.Detail(ctx, order.ComboId)
    if err != nil {
        return err
    }
    var user model.User
    err = dao.User.Ctx(ctx).WherePri(order.UserId).Scan(&user)
    if err == nil {
        _, err = dao.User.Ctx(ctx).WherePri(order.UserId).Update(g.Map{
            dao.User.Columns.BatteryReturnAt: user.BatteryReturnAt.Add(time.Duration(combo.Days) * 24 * time.Hour),
        })
    }
    return err
}

// PenaltyComboSuccess 用户支付违约金
func (*userService) PenaltyComboSuccess(ctx context.Context, order model.ComboOrder) error {
    var user model.User
    err := dao.User.Ctx(ctx).WherePri(order.UserId).Scan(&user)
    if err == nil {
        _, err = dao.User.Ctx(ctx).WherePri(order.UserId).Update(g.Map{
            // 更新用户的换电时间到当日凌晨
            dao.User.Columns.BatteryReturnAt: gtime.NewFromStr(fmt.Sprintf("%d-%02d-%02d 23:59:59", order.CreatedAt.Year(), order.CreatedAt.Month(), order.CreatedAt.Day())),
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

// GroupUserSignDone 团签用户选择电池型号
func (s *userService) GroupUserSignDone(ctx context.Context, sign model.Sign) error {
    res, err := dao.User.Ctx(ctx).WherePri(sign.UserId).
        Where(dao.User.Columns.GroupId, sign.GroupId).
        WhereIn(dao.User.Columns.BatteryState, []int{model.BatteryStateDefault, model.BatteryStateExit}).
        Update(g.Map{
            dao.User.Columns.BatteryState: model.BatteryStateNew,
            dao.User.Columns.BatteryType:  sign.BatteryType,
        })
    if err != nil {
        return err
    }
    cnt, _err := res.RowsAffected()
    if _err == nil && cnt > 0 {
        return nil
    }
    return errors.New("选择失败")
}

// IncrBizBatteryRenewalCnt 增加用户换电次数
func (s *userService) IncrBizBatteryRenewalCnt(ctx context.Context, userId uint64, cnt uint) error {
    _, err := dao.User.Ctx(ctx).WherePri(userId).Increment(dao.User.Columns.BizBatteryRenewalCnt, cnt)
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

// ComboStartUse 个签用户新购套餐首次领取
func (s *userService) ComboStartUse(ctx context.Context, order model.ComboOrder) error {
    combo, err := ComboService.Detail(ctx, order.ComboId)
    if err != nil {
        return err
    }
    user := s.Detail(ctx, order.UserId)
    // 使用时间按自然天计算
    now := gtime.Now()
    y, m, d := now.Add(time.Duration(combo.Days-1) * 24 * time.Hour).Date()
    returnAt := gtime.NewFromStr(fmt.Sprintf("%d-%d-%d 23:59:59", y, m, d))
    res, err := dao.User.Ctx(ctx).WherePri(order.UserId).
        Where(dao.User.Columns.BatteryState, model.BatteryStateNew).
        Update(g.Map{
            dao.User.Columns.BatteryReturnAt:              returnAt,
            dao.User.Columns.BatteryState:                 model.BatteryStateUse,
            dao.User.Columns.BizBatteryRenewalDaysStartAt: gtime.Now(),
            dao.User.Columns.BizBatteryRenewalCnt:         user.BizBatteryRenewalCnt + 1,
            dao.User.Columns.BizBatteryRenewalDays:        user.BizBatteryRenewalDays + 1,
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
func (s *userService) GroupUserStartUse(ctx context.Context, userId uint64) error {
    c := dao.User.Columns
    u := s.Detail(ctx, userId)
    res, err := dao.User.Ctx(ctx).WherePri(userId).
        Where(c.BatteryState, model.BatteryStateNew).
        Update(g.Map{
            c.BatteryState:                 model.BatteryStateUse,
            c.BizBatteryRenewalDaysStartAt: gtime.Now(),
            c.BizBatteryRenewalDays:        u.BizBatteryRenewalDays + 1,
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

// ListVerifyItems 认证列表
func (s *userService) ListVerifyItems(ctx context.Context, req *model.UserVerifyReq) (total int, items []model.UserVerifyListItem) {
    query := dao.User.Ctx(ctx)
    c := dao.User.Columns
    if req.RealName != "" {
        query = query.WhereLike(c.RealName, "%"+req.RealName+"%")
    }
    if req.Type > 0 {
        query = query.Where(c.Type, req.Type)
    }
    if req.Mobile != "" {
        query = query.Where(c.Mobile, req.Mobile)
    }
    total, _ = query.Count()
    _ = query.OrderDesc(c.CreatedAt).Page(req.PageIndex, req.PageLimit).Scan(&items)
    return
}

func (s *userService) ListPersonalItems(ctx context.Context, req *model.UserListReq) (total int, items []model.UserListItem) {
    query := dao.User.Ctx(ctx)

    c := dao.User.Columns
    layout := "Y-m-d"
    params := mq.ParseStructToQuery(*req, "GroupId")
    query = query.Where(params).Where(c.GroupId, req.GroupId)

    if !req.StartDate.IsZero() {
        query = query.WhereGTE(c.CreatedAt, req.StartDate.Format(layout))
    }
    if !req.EndDate.IsZero() {
        query = query.WhereLTE(c.CreatedAt, req.EndDate.Format(layout))
    }

    fields := mq.FieldsWithTable(user.Table, c)

    _ = query.WithAll().
        Page(req.PageIndex, req.PageLimit).
        OrderDesc(c.CreatedAt).
        Fields(fields).
        Scan(&items)

    now := gtime.Now()
    for k, item := range items {
        // 计算剩余天数
        if !item.BatteryReturnAt.IsZero() && now.Before(item.BatteryReturnAt) {
            items[k].Days = uint(math.Floor(item.BatteryReturnAt.Sub(now).Hours() / 24.0))
        }

        if item.ComboDetail != nil {
            items[k].ComboName = item.ComboDetail.Name
            items[k].ComboType = item.ComboDetail.Type
        }

        if item.Group != nil {
            items[k].GroupName = item.Group.Name
        }

        // // 计算换电次数
        // if len(item.BizItems) > 0 {
        //     for _, bizItem := range item.BizItems {
        //         if bizItem.Type == model.UserBizBatteryRenewal {
        //             items[k].ChangeTimes++
        //         }
        //     }
        // }
    }

    total, _ = query.Count()
    return
}
