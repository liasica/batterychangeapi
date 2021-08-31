package service

import (
	"context"
	"fmt"
	"time"

	"battery/app/dao"
	"battery/app/model"
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
		WhereIn(dao.UserBiz.Columns.Type, []int{model.UserBizNew, model.UserBizBatteryRenewal}).
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
		OrderDesc(dao.UserBiz.Columns.Id).
		Page(req.PageIndex, req.PageLimit)
	if req.BizType == 0 {
		m = m.WhereIn(dao.UserBiz.Columns.Type, []uint{model.UserBizBatteryRenewal, model.UserBizBatterySave, model.UserBizClose})
	} else {
		m = m.Where(dao.UserBiz.Columns.Type, req.BizType)
	}

	if req.Keywords != "" {
		users := []struct {
			Id uint64
		}{}
		_ = dao.User.Where(dao.User.Columns.Mobile, req.Keywords).
			WhereOrLike(dao.User.Columns.RealName, fmt.Sprintf("%%%s%%", req.Keywords)).
			Fields(dao.User.Columns.Id).Scan(&users)

		if len(users) == 0 {
			return []model.UserBiz{}
		}
		userIds := make([]uint64, len(users))
		for i, u := range users {
			userIds[i] = u.Id
		}
		m = m.WhereIn(dao.UserBiz.Columns.UserId, userIds)
	}

	if req.Month > 0 {
		year := int(req.Month / 100)
		month := time.Month(req.Month % 100)
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
		m = m.WhereGTE(dao.UserBiz.Columns.CreatedAt, firstOfMonth)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		et, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%d-%02d-%02d 23:59:59", lastOfMonth.Year(), lastOfMonth.Month(), lastOfMonth.Day()))
		m = m.WhereLTE(dao.UserBiz.Columns.CreatedAt, et)
	}
	if req.UserType == 1 {
		m = m.Where(dao.UserBiz.Columns.GoroupId, 0)
	} else if req.UserType == 2 {
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
	} else if req.UserType == 2 {
		m = m.WhereGT(dao.UserBiz.Columns.GoroupId, 0)
	}
	if req.Month > 0 {
		year := int(req.Month / 100)
		month := time.Month(req.Month % 100)
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
		m = m.WhereGTE(dao.UserBiz.Columns.CreatedAt, firstOfMonth)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		et, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%d-%02d-%02d 23:59:59", lastOfMonth.Year(), lastOfMonth.Month(), lastOfMonth.Day()))
		m = m.WhereLTE(dao.UserBiz.Columns.CreatedAt, et)
	}
	if req.BizType == 0 {
		m = m.WhereIn(dao.UserBiz.Columns.Type, []uint{model.UserBizBatteryRenewal, model.UserBizBatterySave, model.UserBizClose})
	} else {
		m = m.Where(dao.UserBiz.Columns.Type, req.BizType)
	}
	rep.Cnt, _ = m.Count()
	return
}
