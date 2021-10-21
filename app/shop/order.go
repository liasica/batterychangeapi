package shop

import (
    "context"
    "fmt"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "strings"

    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
)

var OrderApi = orderApi{}

type orderApi struct {
}

// Total 订单月份统计
// @Summary 店长-订单月份统计
// @Tags    店长-订单
// @Accept  json
// @Produce  json
// @Param 	month query integer  true "月份 如：202106"
// @Param 	type query integer  false "订单类型 1 新签 2 续费"
// @Router  /sapi/order_total [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopOrderTotalRep} "返回结果"
func (*orderApi) Total(r *ghttp.Request) {
    var req model.ShopOrderTotalReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    res := service.ComboOrderService.ShopMonthTotal(r.Context(), req.Month, r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId, req.Type)
    response.JsonOkExit(r, res)
}

// List 订单列表
// @Summary 店长-订单列表
// @Tags    店长-订单
// @Accept  json
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Param 	month query integer  true "月份 如：202106"
// @Param 	type query integer  false "订单类型 1 新签 2 续费"
// @Param 	keywords query string  false "搜索关键字"
// @Router  /sapi/order [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.ShopOrderListItem} "返回结果"
func (*orderApi) List(r *ghttp.Request) {
    var req model.ShopOrderListReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    orderList := service.ComboOrderService.ShopMonthList(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId, req)
    if len(orderList) > 0 {
        userIds := make([]uint64, 0)
        comboIds := make([]uint, 0)
        for _, order := range orderList {
            userIds = append(userIds, order.UserId)
            comboIds = append(comboIds, order.ComboId)
        }
        userList := service.UserService.GetByIds(r.Context(), userIds)
        userIdList := make(map[uint64]model.User, len(userList))
        for _, user := range userList {
            userIdList[user.Id] = user
        }
        res := make([]model.ShopOrderListItem, len(orderList))
        comboList := service.ComboService.GetByIds(r.Context(), comboIds)
        comboIdList := make(map[uint]model.Combo, len(comboList))
        for _, combo := range comboList {
            comboIdList[combo.Id] = combo
        }
        for key, order := range orderList {
            res[key] = model.ShopOrderListItem{
                Id:         order.Id,
                OrderNo:    order.No,
                Amount:     order.Amount,
                Type:       order.Type,
                UserName:   userIdList[order.UserId].RealName,
                UserMobile: userIdList[order.UserId].Mobile,
                PayAt:      order.PayAt,
                ComboName:  comboIdList[order.ComboId].Name,
            }
        }
        response.JsonOkExit(r, res)
    }
    response.JsonOkExit(r, make([]model.ShopOrderListItem, 0))
}

// ListDetail
// @Summary 店长-订单列表获取订单详情
// @Tags    店长-订单
// @Accept  json
// @Produce  json
// @Param 	code path integer  true "订单记录获取订单详情"
// @Router  /sapi/order/:id [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopManagerComboOrderListDetailRep} "返回结果"
func (*orderApi) ListDetail(r *ghttp.Request) {
    var req model.IdReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    order, _ := service.ComboOrderService.Detail(r.Context(), req.Id)
    if order.ShopId != r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId {
        response.Json(r, response.RespCodeArgs, "无权查看")
    }

    user := service.UserService.Detail(r.Context(), order.UserId)

    combo, _ := service.ComboService.Detail(r.Context(), order.ComboId)
    response.JsonOkExit(r, model.ShopManagerComboOrderListDetailRep{
        UserMobile: user.Mobile,
        UserName:   user.RealName,
        ComboName:  combo.Name,

        BatteryType: combo.BatteryType,
        OrderNo:     order.No,
        Amount:      order.Amount,
        Deposit:     order.Deposit,
        PayType:     order.PayType,
        PayAt:       order.PayAt,
    })
}

// ScanDetail
// @Summary 店长-二维码获取订单详情
// @Tags    店长-订单
// @Accept  json
// @Produce  json
// @Param 	code path string  true "订单二维码扫码获取的code"
// @Router  /sapi/order_scan/:code [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopManagerComboOrderScanDetailRep} "返回结果"
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
        rep := model.ShopManagerComboOrderScanDetailRep{
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
        order, err := service.ComboOrderService.DetailByNo(r.Context(), req.Code)
        if err != nil {
            response.Json(r, response.RespCodeArgs, "二维码错误")
        }
        user := service.UserService.Detail(r.Context(), order.UserId)
        combo, _ := service.ComboService.Detail(r.Context(), order.ComboId)
        rep := model.ShopManagerComboOrderScanDetailRep{
            UserType:    user.Type,
            UserName:    user.RealName,
            UserMobile:  user.Mobile,
            ComboName:   combo.Name,
            ComboAmount: combo.Price,
            BatteryType: combo.BatteryType,
            Amount:      order.Amount,
            Deposit:     order.Deposit,
            PayType:     order.PayType,
            OrderNo:     order.No,
            PayAt:       order.PayAt,
            ClaimState:  1,
        }
        if order.ShopId > 0 {
            rep.ClaimState = 2
        }
        response.JsonOkExit(r, rep)
    }
}

// Claim
// @Summary 店长-认领订单
// @Tags    店长-订单
// @Accept  json
// @Produce  json
// @Param   entity  body model.ShopManagerComboOrderClaimReq true "请求数据"
// @Router  /sapi/order_claim [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*orderApi) Claim(r *ghttp.Request) {
    var req model.ShopManagerComboOrderClaimReq
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
        if err := dao.ComboOrder.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
            // 领取记录
            bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
                CityId:       shop.CityId,
                ShopId:       shop.Id,
                UserId:       user.Id,
                GoroupId:     user.GroupId,
                GoroupUserId: groupUser.Id,
                Type:         model.UserBizNew,
                ComboId:      0,
                BatteryType:  user.BatteryType,
            })
            if err != nil {
                return err
            }
            // 用户状态
            if err := service.UserService.GroupUserStartUse(ctx, user.Id); err != nil {
                return err
            }
            // 电池出库
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
            // 账单入账
            return service.GroupSettlementDetailService.Earning(ctx, user)

            // 人数统计
            // return service.GroupDailyStatService.RiderBizNew(ctx, user.GroupId, user.BatteryType, user.Id)
        }); err == nil {
            response.JsonOkExit(r)
        } else {
            g.Log().Error("店主订单认领错误：", err.Error())
            response.JsonErrExit(r)
        }
    } else {
        order, err := service.ComboOrderService.DetailByNo(r.Context(), req.Code)
        if err != nil {
            response.JsonErrExit(r)
        }
        if order.ShopId > 0 {
            response.Json(r, response.RespCodeArgs, "订单已被认领，不能重复认领")
        }
        combo, _ := service.ComboService.Detail(r.Context(), order.ComboId)
        shop, _ := service.ShopService.Detail(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId)
        if combo.CityId != shop.CityId {
            response.Json(r, response.RespCodeArgs, "订单和门店不在同一城市，不能认领")
        }
        if err := dao.ComboOrder.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
            // 订单状态
            if err := service.ComboOrderService.ShopClaim(ctx, order.No, shop.Id); err != nil {
                return err
            }
            // 领取记录
            bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
                CityId:       shop.CityId,
                ShopId:       shop.Id,
                UserId:       order.UserId,
                GoroupId:     0,
                GoroupUserId: 0,
                Type:         model.UserBizNew,
                ComboId:      combo.Id,
                BatteryType:  combo.BatteryType,
            })
            if err != nil {
                return err
            }
            // 用户状态
            if err := service.UserService.ComboStartUse(ctx, order); err != nil {
                return err
            }
            // 电池出库
            if err := service.ShopService.BatteryOut(ctx, shop.Id, combo.BatteryType, 1); err != nil {
                return err
            }
            user := service.UserService.Detail(ctx, order.UserId)
            if err := service.ShopBatteryRecordService.User(ctx,
                model.ShopBatteryRecordTypeOut,
                model.UserBizNew,
                shop.Id,
                bizId,
                user.RealName,
                combo.BatteryType); err != nil {
                return err
            }
            return nil
        }); err == nil {
            response.JsonOkExit(r)
        } else {
            g.Log().Error("店主订单认领错误：", err.Error())
            response.JsonErrExit(r)
        }
    }
}
