package service

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/library/snowflake"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

var PackagesOrderService = packagesOrderService{}

type packagesOrderService struct {
}

// GenerateOrderNo 获取新的订单号
func (s packagesOrderService) GenerateOrderNo() string {
	id := snowflake.Service().Generate()
	now := gtime.Now()
	return fmt.Sprintf("%d%d%d%012d", now.Year(), now.Month(), now.Day(), id%1000000000000)
}

// Detail 套餐订单详情
func (s *packagesOrderService) Detail(ctx context.Context, id uint64) (rep model.PackagesOrder, err error) {
	err = dao.PackagesOrder.Ctx(ctx).WherePri(id).Scan(&rep)
	return
}

// DetailByNo 套餐订单详情
func (s *packagesOrderService) DetailByNo(ctx context.Context, no string) (rep model.PackagesOrder, err error) {
	err = dao.PackagesOrder.Ctx(ctx).Where(dao.PackagesOrder.Columns.No, no).Scan(&rep)
	return
}

// New 新购订单
func (s *packagesOrderService) New(ctx context.Context, userId uint64, packages model.Packages) (order model.PackagesOrder, err error) {
	no := s.GenerateOrderNo()
	id, insertErr := dao.PackagesOrder.Ctx(ctx).InsertAndGetId(g.Map{
		dao.PackagesOrder.Columns.PackageId: packages.Id,
		dao.PackagesOrder.Columns.Amount:    packages.Amount,
		dao.PackagesOrder.Columns.Earnest:   packages.Earnest,
		dao.PackagesOrder.Columns.PayType:   0,
		dao.PackagesOrder.Columns.No:        no,
		dao.PackagesOrder.Columns.UserId:    userId,
		dao.PackagesOrder.Columns.Type:      model.PackageTypeNew,
	})
	if insertErr == nil {
		err = dao.PackagesOrder.Ctx(ctx).WherePri(id).Scan(&order)
	}
	return order, err
}

// Renewal 续购订单
func (s *packagesOrderService) Renewal(ctx context.Context, payType uint, firstOrder model.PackagesOrder) (order model.PackagesOrder, err error) {
	no := s.GenerateOrderNo()
	id, insertErr := dao.PackagesOrder.Ctx(ctx).InsertAndGetId(g.Map{
		dao.PackagesOrder.Columns.PackageId: firstOrder.PackageId,
		dao.PackagesOrder.Columns.Amount:    firstOrder.Amount - firstOrder.Earnest,
		dao.PackagesOrder.Columns.Earnest:   0,
		dao.PackagesOrder.Columns.PayType:   payType,
		dao.PackagesOrder.Columns.No:        no,
		dao.PackagesOrder.Columns.UserId:    firstOrder.UserId,
		dao.PackagesOrder.Columns.ShopId:    firstOrder.ShopId,
		dao.PackagesOrder.Columns.Type:      model.PackageTypeRenewal,
	})
	if insertErr == nil {
		err = dao.PackagesOrder.WherePri(id).Scan(&order)
	}
	return order, err
}

// Penalty 订单违约金
func (s *packagesOrderService) Penalty(ctx context.Context, payType uint, amount float64, firstOrder model.PackagesOrder) (order model.PackagesOrder, err error) {
	no := s.GenerateOrderNo()
	id, insertErr := dao.PackagesOrder.Ctx(ctx).InsertAndGetId(g.Map{
		dao.PackagesOrder.Columns.PackageId: firstOrder.PackageId,
		dao.PackagesOrder.Columns.Amount:    amount,
		dao.PackagesOrder.Columns.Earnest:   0,
		dao.PackagesOrder.Columns.PayType:   payType,
		dao.PackagesOrder.Columns.No:        no,
		dao.PackagesOrder.Columns.UserId:    firstOrder.UserId,
		dao.PackagesOrder.Columns.ShopId:    firstOrder.ShopId,
		dao.PackagesOrder.Columns.Type:      model.PackageTypePenalty,
	})
	if insertErr == nil {
		err = dao.PackagesOrder.WherePri(id).Scan(&order)
	}
	return order, err
}

// PaySuccess 订单支付成功处理
func (s *packagesOrderService) PaySuccess(ctx context.Context, payAt *gtime.Time, no, PayPlatformNo string, payType uint) error {
	_, err := dao.PackagesOrder.Ctx(ctx).Where(dao.PackagesOrder.Columns.No, no).Update(g.Map{
		dao.PackagesOrder.Columns.PayState:      model.PayStateSuccess,
		dao.PackagesOrder.Columns.PayPlatformNo: PayPlatformNo,
		dao.PackagesOrder.Columns.PayAt:         payAt,
		dao.PackagesOrder.Columns.PayType:       payType,
	})
	return err
}

// ShopClaim 店铺认领订单
func (s *packagesOrderService) ShopClaim(ctx context.Context, no string, shopId uint) error {
	res, err := dao.PackagesOrder.Ctx(ctx).Where(dao.PackagesOrder.Columns.No, no).Where(dao.PackagesOrder.Columns.ShopId, 0).Update(g.Map{
		dao.PackagesOrder.Columns.ShopId:      shopId,
		dao.PackagesOrder.Columns.FirstUserAt: gtime.Now(),
	})
	if err == nil {
		rows, err := res.RowsAffected()
		if rows > 0 && err == nil {
			return nil
		}
		err = errors.New("认领失败")
	}
	return err
}

// ShopMonthTotal 店铺订单月统计
func (s *packagesOrderService) ShopMonthTotal(ctx context.Context, month, shopId uint) (res model.ShopOrderTotalRep) {
	_ = dao.PackagesOrder.Ctx(ctx).
		Fields("sum(amount) as amount, count(*) as cnt").
		Where(dao.PackagesOrder.Columns.Month, month).
		Where(dao.PackagesOrder.Columns.PayState, model.PayStateSuccess).
		Where(dao.PackagesOrder.Columns.ShopId, shopId).Scan(&res)
	return
}

// ShopMonthList 店铺订单月度列表
func (s *packagesOrderService) ShopMonthList(ctx context.Context, shopId uint, filter model.ShopOrderListReq) (res []model.PackagesOrder) {
	_ = dao.PackagesOrder.Ctx(ctx).
		Where(dao.PackagesOrder.Columns.Month, filter.Month).
		Where(dao.PackagesOrder.Columns.ShopId, shopId).
		Where(dao.PackagesOrder.Columns.PayState, model.PayStateSuccess).
		Page(filter.PageIndex, filter.PageLimit).
		Scan(&res)
	return
}
