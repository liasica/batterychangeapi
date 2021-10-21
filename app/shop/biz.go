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
// @Summary 店长-业务-办理获取用户信息
// @Tags    店长-业务
// @Accept  json
// @Produce  json
// @Param 	code path integer  true "用户二维码"
// @Router  /sapi/user_biz_profile/:code [GET]
// @Success 200 {object} response.JsonResponse{data=model.BizProfileRep}  "返回结果"
func (*bizApi) Profile(r *ghttp.Request) {
    var req model.BizNewCdeReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    profile := service.UserService.BizProfile(r.Context(), req.Code)
    if profile.Id == 0 {
        response.Json(r, response.RespCodeArgs, "二维码错误")
    }
    response.JsonOkExit(r, profile)
}

// Post
// @Summary 店长-业务-业务办理提交
// @Tags    店长-业务
// @Accept  json
// @Produce  json
// @Param   entity  body model.UserBizReq true "请求数据"
// @Router  /sapi/user_biz [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
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
    // 换电
    case model.UserBizBatteryRenewal:
        if user.BatteryState != model.BatteryStateUse {
            response.Json(r, response.RespCodeArgs, "没有正在租借中的电池，不能办理换电")
        }
        err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
            if _, err := service.UserBizService.Create(ctx, model.UserBiz{
                ShopId:      shop.Id,
                CityId:      shop.CityId,
                UserId:      user.Id,
                GoroupId:    user.GroupId,
                Type:        model.UserBizBatteryRenewal,
                ComboId:     user.ComboId,
                BatteryType: user.BatteryType,
            }); err != nil {
                return err
            }
            if err := service.UserService.IncrBizBatteryRenewalCnt(ctx, user.Id, 1); err != nil {
                return err
            }
            return nil
        })

    // 寄存
    case model.UserBizBatterySave:
        if user.BatteryState != model.BatteryStateUse {
            response.Json(r, response.RespCodeArgs, "用户不是租借中状态，不能办理寄存")
        }
        if user.GroupId > 0 {
            response.Json(r, response.RespCodeArgs, "团签用户，不能办理寄存")
        }
        err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
            // 用户状态
            if err := service.UserService.BizBatterySave(ctx, user); err != nil {
                return nil
            }
            // 寄存记录
            bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
                CityId:       shop.CityId,
                ShopId:       shop.Id,
                UserId:       user.Id,
                GoroupId:     0,
                GoroupUserId: 0,
                Type:         model.UserBizBatterySave,
                ComboId:      user.ComboId,
                BatteryType:  user.BatteryType,
            })
            if err != nil {
                return err
            }
            // 电池入库
            if err := service.ShopService.BatteryIn(ctx, shop.Id, user.BatteryType, 1); err != nil {
                return err
            }
            if err := service.ShopBatteryRecordService.DriverBiz(ctx,
                model.ShopBatteryRecordTypeIn,
                model.UserBizBatterySave,
                shop.Id,
                bizId,
                user); err != nil {
                return err
            }
            return nil
        })

    // 恢复计费
    case model.UserBizBatteryUnSave:
        if user.BatteryState != model.BatteryStateSave {
            response.Json(r, response.RespCodeArgs, "用户不是寄存中状态，不能办理恢复计费")
        }
        if user.BatteryState == model.BatteryStateExpired {
            response.Json(r, response.RespCodeArgs, "用户套餐已过期，不能办理恢复计费")
        }
        if user.GroupId > 0 {
            response.Json(r, response.RespCodeArgs, "团签用户，不能办理恢复计费")
        }
        err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
            // 用户状态
            if err := service.UserService.BizBatteryUnSave(ctx, user); err != nil {
                return nil
            }
            // 取电记录
            bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
                CityId:       shop.CityId,
                ShopId:       shop.Id,
                UserId:       user.Id,
                GoroupId:     0,
                GoroupUserId: 0,
                Type:         model.UserBizBatteryUnSave,
                ComboId:      user.ComboId,
                BatteryType:  user.BatteryType,
            })
            if err != nil {
                return err
            }
            // 电池出库
            if err := service.ShopService.BatteryOut(ctx, shop.Id, user.BatteryType, 1); err != nil {
                return err
            }
            if err := service.ShopBatteryRecordService.DriverBiz(ctx,
                model.ShopBatteryRecordTypeIn,
                model.UserBizBatteryUnSave,
                shop.Id,
                bizId,
                user); err != nil {
                return err
            }
            return nil
        })

    // 退租
    case model.UserBizClose:
        if user.BatteryState != model.BatteryStateUse &&
            user.BatteryState != model.BatteryStateSave &&
            user.BatteryState != model.BatteryStateExpired {
            response.Json(r, response.RespCodeArgs, "用户未开通或已退租，不能办理退租")
        }

        var refundId uint64
        var refundNo string

        groupUser := service.GroupUserService.GetBuyUserId(r.Context(), user.Id)

        err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
            // 用户状态
            if err := service.UserService.BizBatteryExit(ctx, user); err != nil {
                return nil
            }
            // 业务记录
            bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
                CityId:       shop.CityId,
                ShopId:       shop.Id,
                UserId:       user.Id,
                GoroupId:     user.GroupId,
                GoroupUserId: groupUser.Id,
                Type:         model.UserBizClose,
                ComboId:      user.ComboId,
                BatteryType:  user.BatteryType,
            })
            if err != nil {
                return err
            }
            if user.BatteryState == model.BatteryStateUse {
                // 电池入库
                if err := service.ShopService.BatteryIn(ctx, shop.Id, user.BatteryType, 1); err != nil {
                    return err
                }
                if err := service.ShopBatteryRecordService.DriverBiz(ctx,
                    model.ShopBatteryRecordTypeIn,
                    model.UserBizClose,
                    shop.Id,
                    bizId,
                    user); err != nil {
                    return err
                }
            }
            if user.GroupId > 0 {
                if err := service.GroupSettlementDetailService.Cancel(ctx, user); err != nil {
                    return err
                }
            } else {
                comboOrder, err := service.ComboOrderService.Detail(ctx, user.ComboOrderId)
                if err != nil {
                    return err
                }
                if comboOrder.Deposit > 0 {
                    refundNo = service.RefundService.No()
                    refundId, err = service.RefundService.Create(ctx, user.Id, comboOrder.Id, model.RefundRelationTypeComboOrder, refundNo, "骑手套餐退租", comboOrder.Deposit)
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
// @Summary 店长-业务-换电记录列表
// @Tags    店长-业务
// @Accept  json
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Param 	month 	  query integer  true "月份数字，如 202106"
// @Param 	bizType   query integer  false "业务类型 2 换电 3 寄存(仅个签可用)，5 退租"
// @Param 	userType  query integer  true  "用户类型 1 个签  2 团签"
// @Router  /sapi/biz_record [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.UserBizShopRecordRep} "返回结果"
func (*bizApi) RecordUser(r *ghttp.Request) {
    var req model.UserBizShopRecordReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    list := service.UserBizService.ListShop(r.Context(), req)
    if len(list) > 0 {
        userIds := make([]uint64, len(list))
        comboIds := make([]uint, len(list))
        groupIds := make([]uint, len(list))
        for _, record := range list {
            userIds = append(userIds, record.UserId)
            comboIds = append(comboIds, record.ComboId)
            groupIds = append(groupIds, record.GoroupId)
        }
        userList := service.UserService.GetByIds(r.Context(), userIds)
        comboList := service.ComboService.GetByIds(r.Context(), comboIds)
        groupList := service.GroupService.GetByIds(r.Context(), groupIds)
        userIdList := make(map[uint64]model.User, len(userList))
        comboIdList := make(map[uint]model.Combo, len(comboList))
        groupIdList := make(map[uint]model.Group, len(groupList))
        for _, user := range userList {
            userIdList[user.Id] = user
        }
        for _, combo := range comboList {
            comboIdList[combo.Id] = combo
        }
        for _, group := range groupList {
            groupIdList[group.Id] = group
        }
        res := make([]model.UserBizShopRecordRep, len(list))
        for key, record := range list {
            res[key] = model.UserBizShopRecordRep{
                UserName:   userIdList[record.UserId].RealName,
                UserMobile: userIdList[record.UserId].Mobile,
                ComboName:  comboIdList[record.ComboId].Name,
                BizType:    record.Type,
                At:         record.CreatedAt,
            }
            if record.GoroupId > 0 {
                res[key].GroupName = groupIdList[record.GoroupId].Name
            }
        }
        response.JsonOkExit(r, res)
    }
    response.JsonOkExit(r, make([]model.UserBizShopRecordRep, 0))
}

// RecordUserTotal
// @Summary 店长-业务-记录统计
// @Tags    店长-业务
// @Accept  json
// @Produce  json
// @Param 	month 	query integer  true "月份数字，如 202106"
// @Param 	userType  query integer  true  "业务类型 1 个签  2 团签"
// @Param 	bizType   query integer  false "业务类型 2 换电 3 寄存(仅个签可用)，5 退租"
// @Router  /sapi/biz_record_total [GET]
// @Success 200 {object} response.JsonResponse{data=model.UserBizShopRecordMonthTotalRep} "返回结果"
func (*bizApi) RecordUserTotal(r *ghttp.Request) {
    var req model.UserBizShopRecordMonthTotalReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r, service.UserBizService.ListShopMonthTotal(r.Context(), req))
}
