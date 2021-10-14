package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/library/amap"
    "context"
    "fmt"
)

var DistrictsService = districtsService{}

type districtsService struct {
}

type DistrictEl struct {
    Id       uint          `json:"id"`   // id
    Name     string        `json:"name"` // 名称
    Children []*DistrictEl `json:"children,omitempty"`
}

// Child 获取下级地区
func (s *districtsService) Child(parentId uint64) []model.DistrictsChildRep {
    var rep []model.DistrictsChildRep
    _ = dao.Districts.Where(dao.Districts.Columns.ParentId, parentId).Fields(model.DistrictsChildRep{}).Scan(&rep)
    return rep
}

// CurrentCity 获取当前城市
func (s *districtsService) CurrentCity(ctx context.Context, req model.DistrictsCurrentCityReq) (rep model.DistrictsCurrentCityRep, err error) {
    var geo = new(amap.Regeo)
    geo, err = amap.GetRegeo(fmt.Sprintf("%f", req.Lng), fmt.Sprintf("%f", req.Lat))
    if err != nil {
        return
    }
    cityCode := geo.Regeocode.AddressComponent.Citycode
    err = dao.Districts.Ctx(ctx).Where(dao.Districts.Columns.CityCode, cityCode).Where(dao.Districts.Columns.Level, "city").Fields(rep).Limit(1).Scan(&rep)
    return
}

func (s *districtsService) Detail(ctx context.Context, id uint) (d model.Districts, err error) {
    err = dao.Districts.Ctx(ctx).WherePri(id).Scan(&d)
    return
}

// MapIdName 获取地区名IDMap
func (s *districtsService) MapIdName(ctx context.Context, ids []uint) map[uint]string {
    var list []model.Shop
    rep := map[uint]string{}
    _ = dao.Districts.Ctx(ctx).WhereIn(dao.Shop.Columns.Id, ids).Fields(dao.Shop.Columns.Id, dao.Shop.Columns.Name).Scan(&list)
    for _, d := range list {
        rep[d.Id] = d.Name
    }
    return rep
}

// GetByIds 获取地区名IDMap
func (s *districtsService) GetByIds(ctx context.Context, ids []uint) (list []model.Districts) {
    _ = dao.Districts.Ctx(ctx).WhereIn(dao.Shop.Columns.Id, ids).Scan(&list)
    return
}

func (s *districtsService) ListCityTree() (result []*DistrictEl) {
    var items []model.Districts
    field := dao.Districts.Columns.Level
    order := fmt.Sprintf(`FIELD(%s, '%s', '%s', '%s')`, field, model.DistrictLevelDistrict, model.DistrictLevelCity, model.DistrictLevelProvince)
    _ = dao.Districts.WhereIn(field, []string{model.DistrictLevelProvince, model.DistrictLevelCity, model.DistrictLevelDistrict}).Order(order).Scan(&items)
    ds := make(map[uint][]*DistrictEl)
    cs := make(map[uint][]*DistrictEl)
    for _, item := range items {
        switch item.Level {
        case model.DistrictLevelProvince:
            result = append(result, &DistrictEl{
                Id:       item.Id,
                Name:     item.Name,
                Children: cs[item.Id],
            })
        case model.DistrictLevelCity:
            cs[item.ParentId] = append(cs[item.ParentId], &DistrictEl{Id: item.Id, Name: item.Name, Children: ds[item.Id]})
        case model.DistrictLevelDistrict:
            ds[item.ParentId] = append(ds[item.ParentId], &DistrictEl{Id: item.Id, Name: item.Name})
        }
    }
    return result
}
