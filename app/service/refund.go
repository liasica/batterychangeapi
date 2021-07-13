package service

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/library/snowflake"
	"context"
	"github.com/gogf/gf/frame/g"
)

var RefundService = refundService{}

type refundService struct {
}

func (*refundService) No() string {
	return snowflake.Service().Generate().String()
}

// Create 创建退款单
func (*refundService) Create(ctx context.Context, userId, relationId uint64, relationType uint, no string, amount float64) (uint64, error) {
	id, err := dao.Refund.Ctx(ctx).InsertAndGetId(model.Refund{
		No:           no,
		Amount:       amount,
		UserId:       userId,
		RelationId:   relationId,
		RelationType: relationType,
		State:        model.RefundStateStart,
	})
	return uint64(id), err
}

// Success 退款成功
func (*refundService) Success(ctx context.Context, no, platformRefundNo string) error {
	_, err := dao.Refund.Ctx(ctx).Where(dao.Refund.Columns.No, no).Update(g.Map{
		dao.Refund.Columns.State:            model.RefundStateDone,
		dao.Refund.Columns.PlatformRefundNo: platformRefundNo,
	})
	return err
}
