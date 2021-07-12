package payment

import (
	"battery/app/dao"
	"battery/app/service"
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
)

// packageOrderNewSuccess 新购买套餐支付成功处理
func packageOrderNewSuccess(ctx context.Context, patAt *gtime.Time, no, payPlatformNo string) error {
	return dao.PackagesOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if err := service.PackagesOrderService.PaySuccess(ctx, patAt, no, payPlatformNo); err != nil {
			return err
		}
		order, err := service.PackagesOrderService.DetailByNo(ctx, no)
		if err != nil {
			return err
		}
		if err := service.UserService.BuyPackagesSuccess(ctx, order); err != nil {
			return err
		}
		return nil
	})
}

// packageOrderRenewalSuccess 续购买套餐支付成功处理
func packageOrderRenewalSuccess(ctx context.Context, patAt *gtime.Time, no, payPlatformNo string) error {
	return dao.PackagesOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if err := service.PackagesOrderService.PaySuccess(ctx, patAt, no, payPlatformNo); err != nil {
			return err
		}
		order, err := service.PackagesOrderService.DetailByNo(ctx, no)
		if err != nil {
			return err
		}
		if err := service.UserService.RenewalPackagesSuccess(ctx, order); err != nil {
			return err
		}
		return nil
	})
}

// packageOrderPenaltySuccess 违约金支付成功处理
func packageOrderPenaltySuccess(ctx context.Context, patAt *gtime.Time, no, payPlatformNo string) error {
	return dao.PackagesOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if err := service.PackagesOrderService.PaySuccess(ctx, patAt, no, payPlatformNo); err != nil {
			return err
		}
		order, err := service.PackagesOrderService.DetailByNo(ctx, no)
		if err != nil {
			return err
		}
		if err := service.UserService.BuyPackagesSuccess(ctx, order); err != nil {
			return err
		}
		return nil
	})
}
