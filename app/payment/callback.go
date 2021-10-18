package payment

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/os/gtime"
)

// comboOrderNewSuccess 新购买套餐支付成功处理
func comboOrderNewSuccess(ctx context.Context, payAt *gtime.Time, no, payPlatformNo string, payType uint) error {
    return dao.ComboOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
        order, err := service.ComboOrderService.DetailByNo(ctx, no)
        if err != nil {
            return err
        }
        if order.PayState == model.PayStateSuccess {
            return nil
        }
        if err := service.ComboOrderService.PaySuccess(ctx, payAt, no, payPlatformNo, payType); err != nil {
            return err
        }
        if err := service.UserService.BuyComboSuccess(ctx, order); err != nil {
            return err
        }
        combo, _ := service.ComboService.Detail(ctx, order.ComboId)
        city, _ := service.DistrictsService.Detail(ctx, combo.CityId)

        _, _ = service.MessageService.Create(ctx,
            order.UserId,
            model.MessageTypeUserBizNewSuccess,
            "支付成功",
            "开通套餐支付成功",
            model.MessageDetail{
                Type:         "开通套餐",
                ComboName:    combo.Name,
                CityName:     city.Name,
                BatteryType:  combo.BatteryType,
                PayType:      payType,
                PayAt:        payAt,
                Amount:       order.Amount,
                Deposit:      order.Deposit,
                ComboOrderNo: no,
            })
        return nil
    })
}

// comboOrderRenewalSuccess 续购买套餐支付成功处理
func comboOrderRenewalSuccess(ctx context.Context, payAt *gtime.Time, no, payPlatformNo string, payType uint) error {
    return dao.ComboOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
        order, err := service.ComboOrderService.DetailByNo(ctx, no)
        if err != nil {
            return err
        }
        if order.PayState == model.PayStateSuccess {
            return nil
        }
        if err := service.ComboOrderService.PaySuccess(ctx, payAt, no, payPlatformNo, payType); err != nil {
            return err
        }
        if err := service.UserService.RenewalComboSuccess(ctx, order); err != nil {
            return err
        }
        _, _ = service.MessageService.Create(ctx,
            order.UserId,
            model.MessageTypeUserBizRenewalSuccess,
            "支付成功",
            "开通套餐支付成功",
            model.MessageDetail{
                Type:         "续费套餐",
                PayType:      payType,
                PayAt:        payAt,
                Amount:       order.Amount,
                ComboOrderNo: no,
            })
        return nil
    })
}

// comboOrderPenaltySuccess 违约金支付成功处理
func comboOrderPenaltySuccess(ctx context.Context, payAt *gtime.Time, no, payPlatformNo string, payType uint) error {
    return dao.ComboOrder.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
        order, err := service.ComboOrderService.DetailByNo(ctx, no)
        if err != nil {
            return err
        }
        if order.PayState == model.PayStateSuccess {
            return nil
        }
        if err := service.ComboOrderService.PaySuccess(ctx, payAt, no, payPlatformNo, payType); err != nil {
            return err
        }
        if err := service.UserService.PenaltyComboSuccess(ctx, order); err != nil {
            return err
        }
        return nil
    })
}
