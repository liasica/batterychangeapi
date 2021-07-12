package service

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/library/snowflake"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
)

var ShopService = shopService{}

type shopService struct {
}

// ListUser 骑手获取店铺列表
func (s *shopService) ListUser(ctx context.Context, req model.ShopListUserReq) (rep []model.ShopListUserRep) {
	m := dao.Shop.Ctx(ctx).Fields(
		dao.Shop.Columns.Id,
		dao.Shop.Columns.Name,
		dao.Shop.Columns.Lat,
		dao.Shop.Columns.Lng,
		dao.Shop.Columns.Mobile,
		dao.Shop.Columns.Address,
		dao.Shop.Columns.State,
		fmt.Sprintf("ST_DISTANCE_SPHERE(POINT(lng, lat), POINT(%f, %f)) as distance", req.Lng, req.Lat)).
		Where(dao.Shop.Columns.CityId, req.CityId).
		OrderAsc("distance").
		Page(req.PageIndex, req.PageLimit)
	if req.Name != "" {
		m = m.WhereLike(dao.Shop.Columns.Name, fmt.Sprintf("%%%s%%", req.Name))
	}
	_ = m.Scan(&rep)
	return
}

// Detail 获取店铺详情
func (s *shopService) Detail(ctx context.Context, id uint) (rep model.Shop, err error) {
	err = dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Id, id).Limit(1).Scan(&rep)
	return
}

// DetailByQr 获取店铺详情
func (s *shopService) DetailByQr(ctx context.Context, qr string) (rep model.Shop, err error) {
	err = dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Qr, qr).Limit(1).Scan(&rep)
	return
}

// State 修改店铺状态
func (s *shopService) State(ctx context.Context, id uint, state uint) error {
	_, err := dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Id, id).Update(g.Map{dao.Shop.Columns.State: state})
	return err
}

// BatteryIn 电池入店
func (s *shopService) BatteryIn(ctx context.Context, shopId, batterType, num uint) error {
	if batterType == 60 {
		_, err := dao.Shop.Ctx(ctx).WherePri(shopId).Increment(dao.Shop.Columns.BatteryInCnt60, float64(num))
		return err
	}
	if batterType == 72 {
		_, err := dao.Shop.Ctx(ctx).WherePri(shopId).Increment(dao.Shop.Columns.BatteryInCnt72, float64(num))
		return err
	}
	return errors.New("未知的电池型号")
}

// BatteryOut 电池出店
func (s *shopService) BatteryOut(ctx context.Context, shopId, batterType, num uint) error {
	if batterType == 60 {
		_, err := dao.Shop.Ctx(ctx).WherePri(shopId).Increment(dao.Shop.Columns.BatteryOutCnt60, float64(num))
		return err
	}
	if batterType == 72 {
		_, err := dao.Shop.Ctx(ctx).WherePri(shopId).Increment(dao.Shop.Columns.BatteryOutCnt72, float64(num))
		return err
	}
	return errors.New("未知的电池型号")
}

// BatteryRenewal 换电池
func (s *shopService) BatteryRenewal(ctx context.Context, req model.BizProfileRep) error {
	_, err := UserBizService.Create(ctx, model.UserBiz{
		UserId: req.Id,
		//GoroupId:   req.GroupId,
		Type: model.UserBizBatteryRenewal,
		//PackagesId: req.PackagesId,
	})
	return err
}

// MapIdName 获取店名IDMap
func (*shopService) MapIdName(ctx context.Context, ids []uint) map[uint]string {
	var list []model.Shop
	rep := map[uint]string{}
	_ = dao.Shop.Ctx(ctx).WhereIn(dao.Shop.Columns.Id, ids).Fields(dao.Shop.Columns.Id, dao.Shop.Columns.Name).Scan(&list)
	for _, shop := range list {
		rep[shop.Id] = shop.Name
	}
	return rep
}

// CheckMobile 检测电话是否可用
func (*shopService) CheckMobile(ctx context.Context, shopId uint, mobile string) bool {
	if shopId == 0 {
		cnt, _ := dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Mobile, mobile).Count()
		return cnt == 0
	}
	var shop model.Shop
	_ = dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Mobile, mobile).Scan(&shop)
	return shop.Id == 0 || shop.Id == shopId
}

// CheckName 检测名称是否可用
func (*shopService) CheckName(ctx context.Context, shopId uint, name string) bool {
	if shopId == 0 {
		cnt, _ := dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Name, name).Count()
		return cnt == 0
	}
	var shop model.Shop
	_ = dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Mobile, name).Scan(&shop)
	return shop.Id == 0 || shop.Id == shopId
}

// Create 创建店铺
func (*shopService) Create(ctx context.Context, shop model.Shop) (uint, error) {
	shop.Qr = fmt.Sprintf("%d", snowflake.Service().Generate())
	id, err := dao.Shop.Ctx(ctx).Fields(
		dao.Shop.Columns.Name,
		dao.Shop.Columns.ManagerName,
		dao.Shop.Columns.Mobile,
		dao.Shop.Columns.ProvinceId,
		dao.Shop.Columns.CityId,
		dao.Shop.Columns.DistrictId,
		dao.Shop.Columns.Address,
		dao.Shop.Columns.Lng,
		dao.Shop.Columns.Lat,
		dao.Shop.Columns.BatteryInCnt60,
		dao.Shop.Columns.BatteryInCnt72,
		dao.Shop.Columns.State,
		dao.Shop.Columns.Qr,
	).InsertAndGetId(shop)
	return uint(id), err
}

// Edit 编辑店铺
func (*shopService) Edit(ctx context.Context, shop model.Shop) error {
	_, err := dao.Shop.Ctx(ctx).WherePri(shop.Id).Fields(
		dao.Shop.Columns.Name,
		dao.Shop.Columns.ManagerName,
		dao.Shop.Columns.Mobile,
		dao.Shop.Columns.ProvinceId,
		dao.Shop.Columns.CityId,
		dao.Shop.Columns.DistrictId,
		dao.Shop.Columns.Address,
		dao.Shop.Columns.Lng,
		dao.Shop.Columns.Lat,
		dao.Shop.Columns.State,
	).Update(shop)
	return err
}

// ListAdmin 管理员获取店铺列表
func (s *shopService) ListAdmin(ctx context.Context, req model.ShopListAdminReq) (total int, items []model.Shop) {
	m := dao.Shop.Ctx(ctx).Page(req.PageIndex, req.PageLimit)
	if req.Name != "" {
		m = m.WhereLike(dao.Shop.Columns.Name, fmt.Sprintf("%%%s%%", req.Name))
	}
	total, _ = m.Count()
	if total > 0 {
		_ = m.Scan(&items)
	}
	return
}
