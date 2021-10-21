package admin

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
)

var ShopApi = shopApi{}

type shopApi struct {
}

// List
// @Summary 门店列表
// @Tags    管理
// @Accept  json
// @Produce json
// @Param   entity body model.ShopListAdminReq true "门店列表请求"
// @Router  /admin/shop [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.ShopListItem}}  "返回结果"
func (*shopApi) List(r *ghttp.Request) {
    var req model.ShopListAdminReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }

    total, items := service.ShopService.ListAdmin(r.Context(), req)
    response.JsonOkExit(r, g.Map{
        "total": total,
        "items": items,
    })
}

// Create
// @Summary 创建门店
// @Tags    管理
// @Accept  json
// @Param   entity body model.ShopDetail true "门店详情"
// @Produce json
// @Router  /admin/shop [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*shopApi) Create(r *ghttp.Request) {
    var req model.ShopDetail
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if !service.ShopService.CheckMobile(r.Context(), 0, req.Mobile) {
        response.Json(r, response.RespCodeArgs, "手机号码已被使用")
    }
    if !service.ShopService.CheckName(r.Context(), 0, req.Name) {
        response.Json(r, response.RespCodeArgs, "门店名称已被使用")
    }
    if dao.Shop.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
        shop := &model.Shop{
            Name:        req.Name,
            Mobile:      req.Mobile,
            State:       req.State,
            ProvinceId:  req.ProvinceId,
            CityId:      req.CityId,
            DistrictId:  req.DistrictId,
            Address:     req.Address,
            Lng:         req.Lng,
            Lat:         req.Lat,
            ManagerName: req.ManagerName,
        }
        err := service.ShopService.Create(ctx, shop)
        if err != nil {
            return err
        }
        if _, _err := service.ShopManagerService.Create(ctx, model.ShopManager{
            Name:   req.ManagerName,
            Mobile: req.Mobile,
            ShopId: shop.Id,
        }); _err != nil {
            return _err
        }
        // 电池入库记录
        sysUser := ctx.Value(model.ContextAdminKey).(*model.ContextAdmin)
        if err := service.ShopBatteryRecordService.Transfer(ctx, model.ShopBatteryRecord{
            ShopId:      shop.Id,
            Type:        model.ShopBatteryRecordTypeIn,
            SysUserId:   sysUser.Id,
            SysUserName: sysUser.Username,
            BatteryType: model.BatteryType60,
            Num:         req.V60,
        }, shop); err != nil {
            return err
        }

        // 电池入库
        return service.ShopBatteryRecordService.Transfer(ctx, model.ShopBatteryRecord{
            ShopId:      shop.Id,
            Type:        model.ShopBatteryRecordTypeIn,
            SysUserId:   sysUser.Id,
            SysUserName: sysUser.Username,
            BatteryType: model.BatteryType72,
            Num:         req.V72,
        }, shop)

    }) != nil {
        response.JsonErrExit(r)
    }
    response.JsonOkExit(r)
}

// Edit
// @Summary 编辑门店
// @Tags    管理
// @Accept  json
// @Param   id path int true "门店ID"
// @Param   entity body model.ModifyShopReq true "门店详情"
// @Produce json
// @Router  /admin/shop/{id} [PUT]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*shopApi) Edit(r *ghttp.Request) {
    req := new(model.ModifyShopReq)
    _ = request.ParseRequest(r, req)
    id := r.GetUint("id")
    if !service.ShopService.CheckMobile(r.Context(), id, req.Mobile) {
        response.Json(r, response.RespCodeArgs, "手机号码已被使用")
    }
    if !service.ShopService.CheckName(r.Context(), id, req.Name) {
        response.Json(r, response.RespCodeArgs, "门店名称已被使用")
    }
    shop, err := service.ShopService.GetShop(r.Context(), id)
    if err != nil {
        response.Json(r, response.RespCodeArgs, "未找到门店")
    }
    // 查找店铺
    if dao.Shop.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
        if shop.Mobile != req.Mobile {
            if _, err := service.ShopManagerService.Create(ctx, model.ShopManager{
                Name:   req.ManagerName,
                Mobile: req.Mobile,
                ShopId: shop.Id,
            }); err != nil {
                return err
            }
            if err := service.ShopManagerService.Delete(ctx, shop.Mobile); err != nil {
                return err
            }
        }
        _, err = dao.Shop.Ctx(ctx).Data(req).Save()
        return err
    }) != nil {
        response.JsonErrExit(r)
    }
    response.JsonOkExit(r)
}

// Detail
// @Summary 门店详情
// @Tags    管理
// @Accept  json
// @Param   id path int true "门店ID"
// @Produce json
// @Router  /admin/shop/{id} [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopDetail} "返回结果"
func (*shopApi) Detail(r *ghttp.Request) {
    var req model.IdReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    var shop *model.ShopDetail
    err := dao.Shop.Ctx(r.Context()).Where(dao.Shop.Columns.Id, req.Id).Limit(1).Scan(&shop)
    if err != nil {
        response.JsonErrExit(r, response.RespCodeNotFound)
    }
    response.JsonOkExit(r, shop)
}

// ListIdName
// @Summary 门店选择列表(id, name)
// @Tags    管理
// @Accept  json
// @Produce json
// @Router  /admin/shop/idname [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopIdNameList}  "返回结果"
func (*shopApi) ListIdName(r *ghttp.Request) {
    c := dao.Shop.Columns
    var items []model.ShopIdNameList
    _ = dao.Shop.Ctx(r.Context()).OrderAsc(c.CreatedAt).Fields(c.Id, c.Name).Scan(&items)
    response.JsonOkExit(r, items)
}
