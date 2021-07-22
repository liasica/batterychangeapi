package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
)

var UserBizService = userBizService{}

type userBizService struct {
}

// Create 添加记录
func (*userBizService) Create(ctx context.Context, req model.UserBiz) (uint64, error) {
	id, err := dao.UserBiz.Ctx(ctx).InsertAndGetId(req)
	return uint64(id), err
}

// ListUser 骑手获取换电记录
func (*userBizService) ListUser(ctx context.Context, req model.Page) (rep []model.UserBiz) {
	user := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
	_ = dao.UserBiz.Ctx(ctx).
		Where(dao.UserBiz.Columns.UserId, user.Id).
		Where(dao.UserBiz.Columns.Type, model.UserBizBatteryRenewal).
		OrderDesc(dao.UserBiz.Columns.Id).
		Page(req.PageIndex, req.PageLimit).
		Scan(&rep)
	return
}

// ListShop 店铺获取换电记录
func (*userBizService) ListShop(ctx context.Context, req model.UserBizShopRecordReq) (rep []model.UserBiz) {
	manager := ctx.Value(model.ContextShopManagerKey).(*model.ContextShopManager)
	m := dao.UserBiz.Ctx(ctx).
		Where(dao.UserBiz.Columns.ShopId, manager.ShopId).
		WhereIn(dao.UserBiz.Columns.Type, []uint{model.UserBizBatteryRenewal, model.UserBizBatterySave, model.UserBizClose}).
		OrderDesc(dao.UserBiz.Columns.Id).
		Page(req.PageIndex, req.PageLimit)
	if req.UserType == 1 {
		m = m.Where(dao.UserBiz.Columns.GoroupId, 0)
	} else {
		m = m.WhereGT(dao.UserBiz.Columns.GoroupId, 0)
	}
	_ = m.Scan(&rep)
	return
}

// ListShopMonthTotal 店铺获取换电记录按月统计
func (*userBizService) ListShopMonthTotal(ctx context.Context, req model.UserBizShopRecordMonthTotalReq) (rep model.UserBizShopRecordMonthTotalRep) {
	manager := ctx.Value(model.ContextShopManagerKey).(*model.ContextShopManager)
	m := dao.UserBiz.Ctx(ctx).
		Where(dao.UserBiz.Columns.ShopId, manager.ShopId).
		WhereIn(dao.UserBiz.Columns.Type, []uint{model.UserBizBatteryRenewal, model.UserBizBatterySave, model.UserBizClose}).
		OrderDesc(dao.UserBiz.Columns.Id)
	if req.UserType == 1 {
		m = m.Where(dao.UserBiz.Columns.GoroupId, 0)
	} else {
		m = m.WhereGT(dao.UserBiz.Columns.GoroupId, 0)
	}
	rep.Cnt, _ = m.Count()
	return
}

// UserLastSave 获取个签用户最后一次寄存记录
func (*userBizService) UserLastSave(ctx context.Context, userId uint64) (res model.UserBiz, err error) {
	err = dao.UserBiz.Ctx(ctx).
		Where(dao.UserBiz.Columns.UserId, userId).
		Where(dao.UserBiz.Columns.Type, model.UserBizBatterySave).
		OrderDesc(dao.UserBiz.Columns.Id).Limit(1).Scan(&res)
	return
}
