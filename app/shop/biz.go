package shop

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/net/ghttp"

	"battery/app/dao"
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
)

var UserBizApi = bizApi{}

type bizApi struct {
}

// Profile
// @summary 店长-业务-办理获取用户信息
// @tags    店长-业务
// @Accept  json
// @Produce  json
// @param 	code path integer  true "用户二维码"
// @router  /sapi/user_biz_profile/:code [GET]
// @success 200 {object} response.JsonResponse{data=model.BizProfileRep}  "返回结果"
func (*bizApi) Profile(r *ghttp.Request) {
	var req model.BizNewCdeReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	profile := service.UserService.BizProfile(r.Context(), req.Code)
	if profile.Id == 0 {
		response.Json(r, response.RespCodeArgs, "code错误")
	}
	response.JsonOkExit(r, profile)
}

// Post
// @summary 店长-业务-业务办理提交
// @tags    店长-业务
// @Accept  json
// @Produce  json
// @param   entity  body model.UserBizReq true "请求数据"
// @router  /sapi/user_biz [POST]
// @success 200 {object} response.JsonResponse "返回结果"
func (*bizApi) Post(r *ghttp.Request) {
	var req model.UserBizReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	profile := service.UserService.BizProfile(r.Context(), req.Code)

	if profile.BatteryState == model.BatteryStateOverdue {
		response.Json(r, response.RespCodeArgs, "用户已经逾期，请提醒用户先缴纳违约金")
	}
	shop, _ := service.ShopService.Detail(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId)
	user := service.UserService.Detail(r.Context(), profile.Id)
	var err error
	switch req.Type {
	case model.UserBizBatterySave: //寄存
		if user.BatteryState != model.BatteryStateUse {
			response.Json(r, response.RespCodeArgs, "用户不是租借中状态，不能办理寄存")
		}
		if user.GroupId > 0 {
			response.Json(r, response.RespCodeArgs, "团签用户，不能办理寄存")
		}
		err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
			//用户状态
			if err := service.UserService.BizBatterySave(ctx, user); err != nil {
				return nil
			}
			//寄存记录
			bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
				CityId:       shop.CityId,
				ShopId:       shop.Id,
				UserId:       user.Id,
				GoroupId:     0,
				GoroupUserId: 0,
				Type:         model.UserBizBatterySave,
				PackagesId:   user.PackagesId,
				BatteryType:  user.BatteryType,
			})
			if err != nil {
				return err
			}
			//电池入库
			if err := service.ShopService.BatteryIn(ctx, shop.Id, user.BatteryType, 1); err != nil {
				return err
			}
			if err := service.ShopBatteryRecordService.User(ctx,
				model.ShopBatteryRecordTypeIn,
				model.UserBizBatterySave,
				shop.Id,
				bizId,
				user.RealName,
				user.BatteryType); err != nil {
				return err
			}
			return nil
		})

	case model.UserBizBatteryUnSave: //恢复计费
		if user.BatteryState != model.BatteryStateSave {
			response.Json(r, response.RespCodeArgs, "用户不是寄存中状态，不能办理恢复计费")
		}
		if user.GroupId > 0 {
			response.Json(r, response.RespCodeArgs, "团签用户，不能办理恢复计费")
		}
		err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
			//用户状态
			if err := service.UserService.BizBatteryUnSave(ctx, user); err != nil {
				return nil
			}
			//取电记录
			bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
				CityId:       shop.CityId,
				ShopId:       shop.Id,
				UserId:       user.Id,
				GoroupId:     0,
				GoroupUserId: 0,
				Type:         model.UserBizBatteryUnSave,
				PackagesId:   user.PackagesId,
				BatteryType:  user.BatteryType,
			})
			if err != nil {
				return err
			}
			//电池出库
			if err := service.ShopService.BatteryOut(ctx, shop.Id, user.BatteryType, 1); err != nil {
				return err
			}
			if err := service.ShopBatteryRecordService.User(ctx,
				model.ShopBatteryRecordTypeIn,
				model.UserBizBatterySave,
				shop.Id,
				bizId,
				user.RealName,
				user.BatteryType); err != nil {
				return err
			}
			return nil
		})

	case model.UserBizClose: //退租
		if user.BatteryState != model.BatteryStateUse && user.BatteryState != model.BatteryStateSave {
			response.Json(r, response.RespCodeArgs, "用户未开通或已退租，不能办理退租")
		}

		var refundId uint64
		var refundNo string

		err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
			//用户状态
			if err := service.UserService.BizBatteryExit(ctx, user); err != nil {
				return nil
			}
			//业务记录
			bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
				CityId:       shop.CityId,
				ShopId:       shop.Id,
				UserId:       user.Id,
				GoroupId:     0,
				GoroupUserId: 0,
				Type:         model.UserBizClose,
				PackagesId:   user.PackagesId,
				BatteryType:  user.BatteryType,
			})
			if err != nil {
				return err
			}
			if user.BatteryState == model.BatteryStateUse {
				//电池入库
				if err := service.ShopService.BatteryIn(ctx, shop.Id, user.BatteryType, 1); err != nil {
					return err
				}
				if err := service.ShopBatteryRecordService.User(ctx,
					model.ShopBatteryRecordTypeIn,
					model.UserBizBatterySave,
					shop.Id,
					bizId,
					user.RealName,
					user.BatteryType); err != nil {
					return err
				}
			}
			if user.GroupId > 0 {
				if err := service.GroupDailyStatService.RiderBizExit(ctx, user.GroupId, user.BatteryType, user.Id); err != nil {
					return err
				}
			} else {
				packagesOrder, err := service.PackagesOrderService.Detail(ctx, user.PackagesOrderId)
				if err != nil {
					return err
				}
				if packagesOrder.Earnest > 0 {
					refundNo = service.RefundService.No()
					refundId, err = service.RefundService.Create(ctx, user.Id, packagesOrder.Id, model.RefundRelationTypePackagesOrder, refundNo, "骑手套餐退租", packagesOrder.Earnest)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	}
	if err == nil {
		response.JsonOkExit(r)
	}
	response.JsonErrExit(r)
}

// RecordUser
// @summary 店长-业务-记录列表
// @tags    店长-业务
// @Accept  json
// @Produce  json
// @param 	pageIndex query integer  true "当前页码"
// @param 	pageLimit query integer  true "每页行数"
// @param 	month 	  query integer  true "月份数字，如 202106"
// @param 	bizType   query integer  false "业务类型 2 换电 3 寄存(仅个签可用)，5 退租"
// @param 	userType  query integer  true  "业务类型 1 个签  2 团签"
// @router  /sapi/biz_record [GET]
// @success 200 {object} response.JsonResponse{data=[]model.UserBizShopRecordRep} "返回结果"
func (*bizApi) RecordUser(r *ghttp.Request) {
	var req model.UserBizShopRecordReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	list := service.UserBizService.ListShop(r.Context(), req)
	if len(list) > 0 {
		userIds := make([]uint64, len(list))
		packagesIds := make([]uint, len(list))
		groupIds := make([]uint, len(list))
		for _, record := range list {
			userIds = append(userIds, record.UserId)
			packagesIds = append(packagesIds, record.PackagesId)
			groupIds = append(groupIds, record.GoroupId)
		}
		userList := service.UserService.GetByIds(r.Context(), userIds)
		packagesList := service.PackagesService.GetByIds(r.Context(), packagesIds)
		groupList := service.GroupService.GetByIds(r.Context(), packagesIds)
		userIdList := make(map[uint64]model.User, len(userList))
		packagesIdList := make(map[uint]model.Packages, len(packagesList))
		groupIdList := make(map[uint]model.Group, len(groupList))
		for _, user := range userList {
			userIdList[user.Id] = user
		}
		for _, packages := range packagesList {
			packagesIdList[packages.Id] = packages
		}
		for _, group := range groupList {
			groupIdList[group.Id] = group
		}
		res := make([]model.UserBizShopRecordRep, len(list))
		for key, record := range list {
			res[key] = model.UserBizShopRecordRep{
				UserName:     userIdList[record.UserId].RealName,
				UserMobile:   userIdList[record.UserId].Mobile,
				PackagesName: packagesIdList[record.PackagesId].Name,
				BizType:      record.Type,
				At:           record.CreatedAt,
			}
			if record.GoroupUserId > 0 {
				res[key].GroupName = groupIdList[record.GoroupId].Name
			}
		}
	}
	response.JsonOkExit(r, make([]model.UserBizShopRecordRep, 0))
}

// RecordUserTotal
// @summary 店长-业务-记录统计
// @tags    店长-业务
// @Accept  json
// @Produce  json
// @param 	month 	query integer  true "月份数字，如 202106"
// @param 	userType  query integer  true  "业务类型 1 个签  2 团签"
// @router  /sapi/biz_record_total [GET]
// @success 200 {object} response.JsonResponse{data=model.UserBizShopRecordMonthTotalRep} "返回结果"
func (*bizApi) RecordUserTotal(r *ghttp.Request) {
	var req model.UserBizShopRecordMonthTotalReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	response.JsonOkExit(r, service.UserBizService.ListShopMonthTotal(r.Context(), req))
}
