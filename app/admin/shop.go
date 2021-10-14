package admin

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/net/ghttp"
)

var ShopApi = shopApi{}

type shopApi struct {
}

// List
// @summary 门店列表
// @tags    管理,门店
// @Accept  json
// @Produce  json
// @router  /admin/shop [GET]
// @success 200 {object} response.JsonResponse{data=model.ShopListResp}  "返回结果"
func (*shopApi) List(r *ghttp.Request) {
    var req model.ShopListAdminReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }

    var rep model.ShopListResp
    total, items := service.ShopService.ListAdmin(r.Context(), req)
    rep.Total = total
    if rep.Total > 0 {
        cityIds := make([]uint, len(items))
        for key, item := range items {
            cityIds[key] = item.CityId
        }
        cityIdName := service.DistrictsService.MapIdName(r.Context(), cityIds)
        rep.Items = make([]model.ShopListItem, len(items))
        for key, item := range items {
            rep.Items[key] = model.ShopListItem{
                Id:              item.Id,
                Name:            item.Name,
                State:           item.State,
                Mobile:          item.Mobile,
                ManagerName:     item.ManagerName,
                CityName:        cityIdName[item.CityId],
                BatteryInCnt60:  item.BatteryInCnt60,
                BatteryInCnt72:  item.BatteryInCnt72,
                BatteryOutCnt60: item.BatteryOutCnt60,
                BatteryOutCnt72: item.BatteryOutCnt72,
                BatteryCnt60:    item.BatteryCnt60,
                BatteryCnt72:    item.BatteryCnt72,
            }
        }
    }
    response.JsonOkExit(r, rep)
}

// Create
// @summary 创建门店
// @tags    管理,门店
// @Accept  json
// @param   entity body model.CreateShopResp true "门店详情"
// @Produce  json
// @router  /admin/shop [POST]
// @success 200 {object} response.JsonResponse "返回结果"
func (*shopApi) Create(r *ghttp.Request) {
    var req model.CreateShopResp
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
        shopId, err := service.ShopService.Create(ctx, model.Shop{
            Name:           req.Name,
            Mobile:         req.Mobile,
            State:          req.State,
            ProvinceId:     req.ProvinceId,
            CityId:         req.CityId,
            DistrictId:     req.DistrictId,
            Address:        req.Address,
            Lng:            req.Lng,
            Lat:            req.Lat,
            BatteryInCnt60: req.BatteryInCnt60,
            BatteryInCnt72: req.BatteryInCnt72,
            ManagerName:    req.ManagerName,
        })
        if err != nil {
            return err
        }
        if _, _err := service.ShopManagerService.Create(ctx, model.ShopManager{
            Name:   req.ManagerName,
            Mobile: req.Mobile,
            ShopId: shopId,
        }); _err != nil {
            return _err
        }
        // 电池入库记录
        if _err := service.ShopBatteryRecordService.Platform(ctx, model.ShopBatteryRecordTypeIn, shopId, req.BatteryInCnt60, 60); _err != nil {
            return _err
        }
        if _err := service.ShopBatteryRecordService.Platform(ctx, model.ShopBatteryRecordTypeIn, shopId, req.BatteryInCnt72, 72); _err != nil {
            return _err
        }
        return nil
    }) != nil {
        response.JsonErrExit(r)
    }
    response.JsonOkExit(r)
}

type editReq struct {
    Id          uint    `json:"id" v:"required|integer|min:1"`
    Name        string  `json:"name"  v:"required"`
    State       uint    `json:"state" v:"required|in:1,2"`
    ManagerName string  `json:"managerName" v:"required"`
    Mobile      string  `json:"mobile" v:"required|phone-loose"`
    ProvinceId  uint    `json:"provinceId" v:"required|integer|min:1"`
    CityId      uint    `json:"cityId" v:"required|integer|min:1"`
    DistrictId  uint    `json:"districtId" v:"required|integer|min:1"`
    Address     string  `json:"address" v:"required"`
    Lng         float64 `json:"lng" v:"required"`
    Lat         float64 `json:"lat" v:"required"`
}

func (*shopApi) Edit(r *ghttp.Request) {
    var req editReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if !service.ShopService.CheckMobile(r.Context(), req.Id, req.Mobile) {
        response.Json(r, response.RespCodeArgs, "手机号码已被使用")
    }
    if !service.ShopService.CheckName(r.Context(), req.Id, req.Name) {
        response.Json(r, response.RespCodeArgs, "门店名称已被使用")
    }
    if dao.Shop.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
        shop, err := service.ShopService.Detail(ctx, req.Id)
        if err != nil {
            return err
        }
        if shop.Mobile != req.Mobile {
            if _, _err := service.ShopManagerService.Create(ctx, model.ShopManager{
                Name:   req.ManagerName,
                Mobile: req.Mobile,
                ShopId: shop.Id,
            }); _err != nil {
                return _err
            }
            if _err := service.ShopManagerService.Delete(ctx, shop.Mobile); _err != nil {
                return _err
            }
        }
        if _err := service.ShopService.Edit(ctx, model.Shop{
            Id:          req.Id,
            Name:        req.Name,
            ManagerName: req.ManagerName,
            Mobile:      req.Mobile,
            ProvinceId:  req.ProvinceId,
            CityId:      req.CityId,
            DistrictId:  req.DistrictId,
            Address:     req.Address,
            Lng:         req.Lng,
            Lat:         req.Lat,
            State:       req.State,
        }); _err != nil {
            return _err
        }
        return nil
    }) != nil {
        response.JsonErrExit(r)
    }
    response.JsonOkExit(r)
}

func (*shopApi) Detail(r *ghttp.Request) {
    var req model.IdReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    shop, err := service.ShopService.Detail(r.Context(), uint(req.Id))
    if err != nil || shop.Id == 0 {
        response.JsonErrExit(r, response.RespCodeNotFound)
    }
    response.JsonOkExit(r, model.CreateShopResp{
        Name:           shop.Name,
        State:          shop.State,
        ManagerName:    shop.ManagerName,
        Mobile:         shop.Mobile,
        ProvinceId:     shop.ProvinceId,
        BatteryInCnt60: uint(shop.BatteryCnt60),
        BatteryInCnt72: uint(shop.BatteryCnt72),
        CityId:         shop.CityId,
        DistrictId:     shop.DistrictId,
        Address:        shop.Address,
        Lng:            shop.Lng,
        Lat:            shop.Lat,
    })
}
