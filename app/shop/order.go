package shop

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
	"context"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/net/ghttp"
	"strings"
)

var OrderApi = orderApi{}

type orderApi struct {
}

// Total 订单月份统计
// @summary 店长-订单月份统计
// @tags    店长-订单
// @Accept  json
// @Produce  json
// @param 	month query integer  true "月份 如：202106"
// @router  /sapi/order_total [GET]
// @success 200 {object} response.JsonResponse{data=model.ShopOrderTotalRep} "返回结果"
func (*orderApi) Total(r *ghttp.Request) {
	var req model.ShopOrderTotalReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	res := service.PackagesOrderService.ShopMonthTotal(r.Context(), req.Month, r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId)
	response.JsonOkExit(r, res)
}

// List 订单列表
// @summary 店长-订单列表
// @tags    店长-订单
// @Accept  json
// @Produce  json
// @param 	pageIndex query integer  true "当前页码"
// @param 	pageLimit query integer  true "每页行数"
// @param 	month query integer  true "月份 如：202106"
// @param 	keywords query string  false "搜索关键字"
// @router  /sapi/order [GET]
// @success 200 {object} response.JsonResponse{data=[]model.ShopOrderListItem} "返回结果"
func (*orderApi) List(r *ghttp.Request) {
	var req model.ShopOrderListReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	orderList := service.PackagesOrderService.ShopMonthList(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId, req)
	if len(orderList) > 0 {
		userIds := make([]uint64, 0)
		packagesIds := make([]uint, 0)
		for _, order := range orderList {
			userIds = append(userIds, order.UserId)
			packagesIds = append(packagesIds, order.PackageId)
		}
		userList := service.UserService.GetByIds(r.Context(), userIds)
		userIdList := make(map[uint64]model.User, len(userList))
		for _, user := range userList {
			userIdList[user.Id] = user
		}
		res := make([]model.ShopOrderListItem, len(orderList))
		packagesList := service.PackagesService.GetByIds(r.Context(), packagesIds)
		packagesIdList := make(map[uint]model.Packages, len(packagesList))
		for _, packages := range packagesList {
			packagesIdList[packages.Id] = packages
		}
		for key, order := range orderList {
			res[key] = model.ShopOrderListItem{
				Id:          order.Id,
				OrderNo:     order.No,
				Amount:      order.Amount,
				Type:        order.Type,
				UserName:    userList[order.UserId].RealName,
				UserMobile:  userList[order.UserId].Mobile,
				PayAt:       order.PayAt,
				PackageName: packagesIdList[order.PackageId].Name,
			}
		}
		response.JsonOkExit(r, res)
	}
	response.JsonOkExit(r, make([]model.ShopOrderListItem, 0))
}

// ListDetail
// @summary 店长-订单列表获取订单详情
// @tags    店长-订单
// @Accept  json
// @Produce  json
// @param 	code path integer  true "订单记录获取订单详情"
// @router  /sapi/order/:id [GET]
// @success 200 {object} response.JsonResponse{data=model.ShopManagerPackagesOrderListDetailRep} "返回结果"
func (*orderApi) ListDetail(r *ghttp.Request) {
	var req model.IdReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	order, _ := service.PackagesOrderService.Detail(r.Context(), req.Id)
	if order.ShopId != r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId {
		response.Json(r, response.RespCodeArgs, "无权查看")
	}
	packages, _ := service.PackagesService.Detail(r.Context(), order.PackageId)
	response.JsonOkExit(r, model.ShopManagerPackagesOrderListDetailRep{
		BatteryType: packages.BatteryType,
		OrderNo:     order.No,
		Amount:      order.Amount,
		Earnest:     order.Earnest,
		PayType:     order.PayType,
		PayAt:       order.PayAt,
	})
}

// ScanDetail
// @summary 店长-二维码获取订单详情
// @tags    店长-订单
// @Accept  json
// @Produce  json
// @param 	code path string  true "订单二维码扫码获取的code"
// @router  /sapi/order_scan/:code [GET]
// @success 200 {object} response.JsonResponse{data=model.ShopManagerPackagesOrderScanDetailRep} "返回结果"
func (*orderApi) ScanDetail(r *ghttp.Request) {
	var req model.BizNewCdeReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	codes := strings.Split(req.Code, "-")
	if len(codes) == 3 {
		groupId := codes[0]
		userQr := codes[1]
		batteryType := codes[2]
		user := service.UserService.DetailByQr(r.Context(), userQr)
		if user.GroupId == 0 || fmt.Sprintf("%d", user.GroupId) != groupId || fmt.Sprintf("%d", user.BatteryType) != batteryType {
			response.Json(r, response.RespCodeArgs, "参数错误")
		}
		if user.BatteryState == model.BatteryStateDefault {
			response.Json(r, response.RespCodeArgs, "骑手未选择电池类型")
		}
		group := service.GroupService.Detail(r.Context(), user.GroupId)
		rep := model.ShopManagerPackagesOrderScanDetailRep{
			UserName:    user.RealName,
			UserMobile:  user.Mobile,
			UserType:    user.Type,
			BatteryType: user.BatteryType,
			GroupName:   group.Name,
			ClaimState:  1,
		}
		if user.BatteryState > model.BatteryStateNew {
			rep.ClaimState = 2
		}
		response.JsonOkExit(r, rep)
	} else {
		order, err := service.PackagesOrderService.DetailByNo(r.Context(), req.Code)
		if err != nil {
			response.JsonErrExit(r)
		}
		user := service.UserService.Detail(r.Context(), order.UserId)
		packages, _ := service.PackagesService.Detail(r.Context(), order.PackageId)
		rep := model.ShopManagerPackagesOrderScanDetailRep{
			UserType:     user.Type,
			UserName:     user.RealName,
			UserMobile:   user.Mobile,
			PackagesName: packages.Name,
			BatteryType:  packages.BatteryType,
			Amount:       order.Amount,
			Earnest:      order.Earnest,
			PayType:      order.PayType,
			OrderNo:      order.No,
			PayAt:        order.PayAt,
			ClaimState:   1,
		}
		if order.ShopId > 0 {
			rep.ClaimState = 2
		}
		response.JsonOkExit(r, rep)
	}
}

// Claim
// @summary 店长-认领订单
// @tags    店长-订单
// @Accept  json
// @Produce  json
// @param   entity  body model.ShopManagerPackagesOrderClaimReq true "请求数据"
// @router  /sapi/order_claim [POST]
// @success 200 {object} response.JsonResponse "返回结果"
func (*orderApi) Claim(r *ghttp.Request) {
	var req model.ShopManagerPackagesOrderClaimReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	codes := strings.Split(req.Code, "-")
	if len(codes) == 3 {
		userQr := codes[1]
		user := service.UserService.DetailByQr(r.Context(), userQr)
		if user.Id == 0 {
			response.Json(r, response.RespCodeArgs, "未知用户")
		}
		if user.GroupId == 0 {
			response.Json(r, response.RespCodeArgs, "不是团签用户")
		}
		if user.BatteryState != model.BatteryStateNew {
			response.Json(r, response.RespCodeArgs, "没有选择电池类型, 或已领取")
		}
		groupUser := service.GroupUserService.GetBuyUserId(r.Context(), user.Id)
		shop, _ := service.ShopService.Detail(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId)
		if dao.PackagesOrder.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
			//领取记录
			bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
				CityId:       shop.CityId,
				ShopId:       shop.Id,
				UserId:       user.Id,
				GoroupId:     user.GroupId,
				GoroupUserId: groupUser.Id,
				Type:         model.UserBizNew,
				PackagesId:   0,
				BatteryType:  user.BatteryType,
			})
			if err != nil {
				return err
			}
			//用户状态
			if err := service.UserService.GroupUserStartUse(ctx, user.Id); err != nil {
				return err
			}
			//电池出库
			if err := service.ShopService.BatteryOut(ctx, shop.Id, user.BatteryType, 1); err != nil {
				return err
			}
			if err := service.ShopBatteryRecordService.User(ctx,
				model.ShopBatteryRecordTypeOut,
				model.UserBizNew,
				shop.Id,
				bizId,
				user.RealName,
				user.BatteryType); err != nil {
				return err
			}
			//人数统计
			if err := service.GroupDailyStatService.RiderBizNew(ctx, user.GroupId, user.BatteryType); err != nil {
				return err
			}
			return nil
		}) == nil {
			response.JsonOkExit(r)
		}
		response.JsonErrExit(r)
	} else {
		order, err := service.PackagesOrderService.DetailByNo(r.Context(), req.Code)
		if err != nil {
			response.JsonErrExit(r)
		}
		if order.ShopId > 0 {
			response.Json(r, response.RespCodeArgs, "订单已被认领，不能重复认领")
		}
		packages, _ := service.PackagesService.Detail(r.Context(), order.PackageId)
		shop, _ := service.ShopService.Detail(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId)
		if packages.CityId != shop.CityId {
			response.Json(r, response.RespCodeArgs, "订单和店铺不在同一城市，不能认领")
		}
		if dao.PackagesOrder.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
			//订单状态
			if err := service.PackagesOrderService.ShopClaim(ctx, order.No, shop.Id); err != nil {
				return err
			}
			//领取记录
			bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
				CityId:       shop.CityId,
				ShopId:       shop.Id,
				UserId:       order.UserId,
				GoroupId:     0,
				GoroupUserId: 0,
				Type:         model.UserBizNew,
				PackagesId:   packages.Id,
				BatteryType:  packages.BatteryType,
			})
			if err != nil {
				return err
			}
			//用户状态
			if err := service.UserService.PackagesStartUse(ctx, order); err != nil {
				return err
			}
			//电池出库
			if err := service.ShopService.BatteryOut(ctx, shop.Id, packages.BatteryType, 1); err != nil {
				return err
			}
			user := service.UserService.Detail(ctx, order.UserId)
			if err := service.ShopBatteryRecordService.User(ctx,
				model.ShopBatteryRecordTypeOut,
				model.UserBizNew,
				shop.Id,
				bizId,
				user.RealName,
				packages.BatteryType); err != nil {
				return err
			}
			return nil
		}) == nil {
			response.JsonOkExit(r)
		}
		response.JsonErrExit(r)
	}
}
