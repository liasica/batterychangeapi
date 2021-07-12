package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
)

var PackagesService = packagesService{}

type packagesService struct {
}

// ListUser 用户套餐列表
func (s *packagesService) ListUser(ctx context.Context, req model.PackagesListUserReq) model.PackagesListUserRep {
	var rep model.PackagesListUserRep
	_ = dao.Packages.Ctx(ctx).
		Where(dao.Packages.Columns.CityId, req.CityId).
		Page(req.PageIndex, req.PageLimit).
		Scan(&rep)
	return rep
}

// ListAdmin 管理套餐列表
func (s *packagesService) ListAdmin(ctx context.Context, req model.Page) (total int, items []model.Packages) {
	m := dao.Packages.Ctx(ctx).Page(req.PageIndex, req.PageLimit)
	total, _ = m.Count()
	if total > 0 {
		_ = m.Scan(&items)
	}
	return
}

// Detail 套餐详情
func (s *packagesService) Detail(ctx context.Context, id uint) (rep model.Packages, err error) {
	err = dao.Packages.Ctx(ctx).WherePri(id).Scan(&rep)
	return
}

// PenaltyAmount 获取套餐违约金
func (s *packagesService) PenaltyAmount(ctx context.Context, id, days uint) (amount float64, err error) {
	packages, err := s.Detail(ctx, id)
	if err == nil {
		//TODO 高精度计算
		amount = ((packages.Amount - packages.Earnest) / 30) * float64(days)
	}
	return
}

// GetByIds 套餐详情
func (s *packagesService) GetByIds(ctx context.Context, ids []uint) (rep []model.Packages) {
	_ = dao.Packages.Ctx(ctx).WhereIn(dao.Packages.Columns.Id, ids).Scan(&rep)
	return
}

// GetCityIds 所有城市ID
func (s *packagesService) GetCityIds(ctx context.Context) (rep []uint) {
	var res []struct {
		CityId uint
	}
	_ = dao.Packages.Ctx(ctx).Group(dao.Packages.Columns.CityId).Scan(&res)
	for _, c := range res {
		rep = append(rep, c.CityId)
	}
	return rep
}

// Create 创建套餐
func (s *packagesService) Create(ctx context.Context, req model.Packages) (id int64, err error) {
	id, err = dao.Packages.Ctx(ctx).InsertAndGetId(req)
	return
}
