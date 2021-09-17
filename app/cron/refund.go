package cron

import (
	"battery/app/model"
	"battery/app/service"
	"battery/library/payment/alipay"
	"battery/library/payment/wechat"
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/robfig/cron/v3"
	"strconv"
)

var RefundCron = refund{}

type refund struct {
}

func (*refund) Start() error {
	if !g.Cfg().GetBool("cron.order.refund.enable", false) {
		return nil
	}
	g.Log().Info("退款任务启动")
	c := cron.New()
	_, err := c.AddFunc(g.Cfg().GetString("cron.order.refund.spec"), func() {
		var minId uint64 = 0
		page := model.Page{
			PageIndex: 1,
			PageLimit: 10,
		}
		for {
			list := service.RefundService.WaitList(context.TODO(), page, minId)
			if len(list) == 0 {
				return
			}
			for _, refundOrder := range list {
				minId = uint64(refundOrder.Id)
				page.PageIndex++
				if refundOrder.RelationType == model.RefundRelationTypePackagesOrder {
					packagesOrder, err := service.PackagesOrderService.Detail(context.TODO(), refundOrder.RelationId)
					if err != nil {
						g.Log().Error(err.Error())
					}
					if packagesOrder.PayType == model.PayTypeAliPay {
						platformRefundNo, err := alipay.Service().Refund(context.TODO(),
							packagesOrder.PayPlatformNo,
							packagesOrder.No,
							refundOrder.No,
							strconv.FormatFloat(refundOrder.Amount, 'f', 2, 64),
							refundOrder.Reason)
						if err != nil {
							g.Log().Error("退款失败", err.Error())
							continue
						}
						_ = service.RefundService.Success(context.TODO(), refundOrder.No, platformRefundNo)
					}
					if packagesOrder.PayType == model.PayTypeWechat {
						platformRefundNo, err := wechat.Service().Refund(context.TODO(),
							packagesOrder.PayPlatformNo,
							packagesOrder.No,
							refundOrder.No,
							refundOrder.Reason,
							refundOrder.Amount,
							packagesOrder.Amount)
						if err != nil {
							g.Log().Error("退款失败", err.Error())
							continue
						}
						_ = service.RefundService.Success(context.TODO(), refundOrder.No, platformRefundNo)
					}
				}
			}
		}
	})
	c.Start()
	return err
}
