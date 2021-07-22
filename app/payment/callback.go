package payment

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/app/service"
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
)

// packageOrderNewSuccess 新购买套餐支付成功处理
func packageOrderNewSuccess(ctx context.Context, payAt *gtime.Time, no, payPlatformNo string, payType uint) error {
	return dao.PackagesOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		order, err := service.PackagesOrderService.DetailByNo(ctx, no)
		if err != nil {
			return err
		}
		if order.PayState == model.PayStateSuccess {
			return nil
		}
		if err := service.PackagesOrderService.PaySuccess(ctx, payAt, no, payPlatformNo, payType); err != nil {
			return err
		}
		if err := service.UserService.BuyPackagesSuccess(ctx, order); err != nil {
			return err
		}
		packages, _ := service.PackagesService.Detail(ctx, order.PackageId)
		city, _ := service.DistrictsService.Detail(ctx, packages.CityId)

		_, _ = service.MessageService.Create(ctx,
			order.UserId,
			model.MessageTypeUserBizNewSuccess,
			"支付成功",
			"开通套餐支付成功",
			model.MessageDetail{
				Type:           "开通套餐",
				PackagesName:   packages.Name,
				CityName:       city.Name,
				BatteryType:    packages.BatteryType,
				PayType:        payType,
				PayAt:          payAt,
				Amount:         order.Amount,
				Earnest:        order.Earnest,
				PackageOrderNo: no,
			})
		return nil
	})
}

// packageOrderRenewalSuccess 续购买套餐支付成功处理
func packageOrderRenewalSuccess(ctx context.Context, payAt *gtime.Time, no, payPlatformNo string, payType uint) error {
	return dao.PackagesOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		order, err := service.PackagesOrderService.DetailByNo(ctx, no)
		if err != nil {
			return err
		}
		if order.PayState == model.PayStateSuccess {
			return nil
		}
		if err := service.PackagesOrderService.PaySuccess(ctx, payAt, no, payPlatformNo, payType); err != nil {
			return err
		}
		if err := service.UserService.RenewalPackagesSuccess(ctx, order); err != nil {
			return err
		}
		_, _ = service.MessageService.Create(ctx,
			order.UserId,
			model.MessageTypeUserBizRenewalSuccess,
			"支付成功",
			"开通套餐支付成功",
			model.MessageDetail{
				Type:           "续费套餐",
				PayType:        payType,
				PayAt:          payAt,
				Amount:         order.Amount,
				PackageOrderNo: no,
			})
		return nil
	})
}

// packageOrderPenaltySuccess 违约金支付成功处理
func packageOrderPenaltySuccess(ctx context.Context, payAt *gtime.Time, no, payPlatformNo string, payType uint) error {
	return dao.PackagesOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		order, err := service.PackagesOrderService.DetailByNo(ctx, no)
		if err != nil {
			return err
		}
		if order.PayState == model.PayStateSuccess {
			return nil
		}
		if err := service.PackagesOrderService.PaySuccess(ctx, payAt, no, payPlatformNo, payType); err != nil {
			return err
		}
		if err := service.UserService.PenaltyPackagesSuccess(ctx, order); err != nil {
			return err
		}
		return nil
	})
}
