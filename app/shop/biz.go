package shop

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/net/ghttp"
)

var UserBizApi = bizApi{}

type bizApi struct {
}

// Profile
// @Summary 门店-业务-办理获取用户信息
// @Tags    门店-业务
// @Accept  json
// @Produce json
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
// @Summary 门店-业务-业务办理提交
// @Tags    门店-业务
// @Accept  json
// @Produce json
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
    shop, _ := service.ShopService.GetShop(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId)
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
                GroupId:     user.GroupId,
                BizType:     model.UserBizBatteryRenewal,
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
    case model.UserBizBatteryPause:
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
                CityId:      shop.CityId,
                ShopId:      shop.Id,
                UserId:      user.Id,
                GroupId:     0,
                BizType:     model.UserBizBatteryPause,
                ComboId:     user.ComboId,
                BatteryType: user.BatteryType,
            })
            if err != nil {
                return err
            }
            // 电池入库
            return service.ShopBatteryRecordService.Transfer(
                ctx,
                model.ShopBatteryRecord{
                    ShopId:      shop.Id,
                    BizId:       bizId,
                    BizType:     model.UserBizBatteryPause,
                    UserName:    user.RealName,
                    BatteryType: user.BatteryType,
                    Num:         1,
                    Type:        model.ShopBatteryRecordTypeIn,
                    UserId:      user.Id,
                },
                shop,
            )
        })

    // 恢复计费
    case model.UserBizBatteryRetrieval:
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
                CityId:      shop.CityId,
                ShopId:      shop.Id,
                UserId:      user.Id,
                BizType:     model.UserBizBatteryRetrieval,
                ComboId:     user.ComboId,
                BatteryType: user.BatteryType,
            })
            if err != nil {
                return err
            }
            // 电池出库
            return service.ShopBatteryRecordService.Transfer(
                ctx,
                model.ShopBatteryRecord{
                    ShopId:      shop.Id,
                    BizId:       bizId,
                    BizType:     model.UserBizBatteryRetrieval,
                    UserName:    user.RealName,
                    BatteryType: user.BatteryType,
                    Num:         1,
                    Type:        model.ShopBatteryRecordTypeOut,
                    UserId:      user.Id,
                },
                shop,
            )
        })

    // 退租
    case model.UserBizCancel:
        if user.BatteryState != model.BatteryStateUse &&
            user.BatteryState != model.BatteryStateSave &&
            user.BatteryState != model.BatteryStateExpired {
            response.Json(r, response.RespCodeArgs, "用户未开通或已退租，不能办理退租")
        }

        var refundId uint64
        var refundNo string

        err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
            // 用户状态
            if err := service.UserService.BizBatteryExit(ctx, user); err != nil {
                return nil
            }
            // 业务记录
            bizId, err := service.UserBizService.Create(ctx, model.UserBiz{
                CityId:      shop.CityId,
                ShopId:      shop.Id,
                UserId:      user.Id,
                GroupId:     user.GroupId,
                BizType:     model.UserBizCancel,
                ComboId:     user.ComboId,
                BatteryType: user.BatteryType,
            })
            if err != nil {
                return err
            }
            if user.BatteryState == model.BatteryStateUse {
                // 电池入库
                if err := service.ShopBatteryRecordService.Transfer(
                    ctx,
                    model.ShopBatteryRecord{
                        ShopId:      shop.Id,
                        BizId:       bizId,
                        BizType:     model.UserBizCancel,
                        UserName:    user.RealName,
                        BatteryType: user.BatteryType,
                        Num:         1,
                        Type:        model.ShopBatteryRecordTypeIn,
                        UserId:      user.Id,
                    },
                    shop,
                ); err != nil {
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

// Record
// @Summary 门店-门店业务记录
// @Tags    门店
// @Accept  json
// @Produce json
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页数量"
// @Param 	month 	  query string true "月份, eg: 2021-09"
// @Param 	bizType   query integer false "业务类型: 2 换电 3 寄存(仅个签可用)，5 退租"
// @Param 	userType  query integer false "用户类型: 1个签 2团签"
// @Param 	realName  query string false "筛选用户姓名"
// @Router  /sapi/biz/record [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.BizShopFilterResp} "返回结果"
func (*bizApi) Record(r *ghttp.Request) {
    req := new(model.BizShopFilterReq)
    _ = request.ParseRequest(r, req)
    req.ShopId = r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId
    response.JsonOkExit(r, service.UserBizService.ShopFilter(r.Context(), req))
}
