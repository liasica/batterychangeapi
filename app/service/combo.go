package service

import (
    "battery/app/dao"
    "battery/app/model"
    "context"
    "github.com/shopspring/decimal"
)

var ComboService = comboService{}

type comboService struct {
}

// ListUser 用户套餐列表
func (s *comboService) ListUser(ctx context.Context, req model.ComboListUserReq) model.ComboListUserRep {
    var rep model.ComboListUserRep
    _ = dao.Combo.Ctx(ctx).
        Where(dao.Combo.Columns.CityId, req.CityId).
        Page(req.PageIndex, req.PageLimit).
        Scan(&rep)
    if l := len(rep); l > 0 {
        city, _ := DistrictsService.Detail(ctx, uint(req.CityId))
        cnt, _ := dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.CityId, req.CityId).Count()
        for i := range rep {
            rep[i].UsableCityName = city.Name
            rep[i].UsableShopCnt = cnt
        }
    }
    return rep
}

// ListAdmin 管理套餐列表
func (s *comboService) ListAdmin(ctx context.Context, req model.Page) (total int, items []model.Combo) {
    m := dao.Combo.Ctx(ctx).Page(req.PageIndex, req.PageLimit)
    total, _ = m.Count()
    if total > 0 {
        _ = m.Scan(&items)
    }
    return
}

// Detail 套餐详情
func (s *comboService) Detail(ctx context.Context, id uint) (rep model.Combo, err error) {
    err = dao.Combo.Ctx(ctx).WherePri(id).Scan(&rep)
    return
}

// PenaltyAmount 获取套餐违约金
func (s *comboService) PenaltyAmount(ctx context.Context, id, days uint) (amount float64, err error) {
    combo, err := s.Detail(ctx, id)
    if err == nil {
        amount, _ = decimal.NewFromFloat(combo.Price).
            Div(decimal.NewFromInt(int64(combo.Days))).
            Mul(decimal.NewFromInt(int64(days))).
            Round(2).
            Float64()
    }
    return
}

// GetByIds 套餐详情
func (s *comboService) GetByIds(ctx context.Context, ids []uint) (rep []model.Combo) {
    _ = dao.Combo.Ctx(ctx).WhereIn(dao.Combo.Columns.Id, ids).Scan(&rep)
    return
}

// GetCityIds 所有城市ID
func (s *comboService) GetCityIds(ctx context.Context) (rep []uint) {
    var res []struct {
        CityId uint
    }
    _ = dao.Combo.Ctx(ctx).Group(dao.Combo.Columns.CityId).Scan(&res)
    for _, c := range res {
        rep = append(rep, c.CityId)
    }
    return rep
}

// Create 创建套餐
func (s *comboService) Create(ctx context.Context, req model.Combo) (id int64, err error) {
    id, err = dao.Combo.Ctx(ctx).InsertAndGetId(req)
    return
}
