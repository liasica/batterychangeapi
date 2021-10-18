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
func (*refundService) Create(ctx context.Context, userId, relationId uint64, relationType uint, no, reason string, amount float64) (uint64, error) {
    id, err := dao.Refund.Ctx(ctx).InsertAndGetId(model.Refund{
        No:           no,
        Amount:       amount,
        UserId:       userId,
        RelationId:   relationId,
        Reason:       reason,
        RelationType: relationType,
        State:        model.RefundStatePending,
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

// WaitList 获取待退款列表
func (*refundService) WaitList(ctx context.Context, page model.Page, minId uint64) []model.Refund {
    var res []model.Refund
    _ = dao.Refund.Ctx(ctx).WhereGT(dao.Refund.Columns.Id, minId).
        Where(dao.Refund.Columns.State, model.RefundStatePending).
        OrderAsc(dao.Refund.Columns.Id).
        Page(page.PageIndex, page.PageLimit).
        Scan(&res)
    return res
}
