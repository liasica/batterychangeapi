package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"github.com/gogf/gf/frame/g"
)

var SignService = signService{}

type signService struct {
}

func (*signService) Create(ctx context.Context, data model.Sign) (uint64, error) {
	id, err := dao.Sign.Ctx(ctx).InsertAndGetId(data)
	return uint64(id), err
}

func (*signService) GetDetailBayFlowId(ctx context.Context, flowId string) (model.Sign, error) {
	var sign model.Sign
	err := dao.Sign.Ctx(ctx).Where(dao.Sign.Columns.FlowId, flowId).Limit(1).Scan(&sign)
	return sign, err
}

// UserLatestDetail 根据流程ID查询用户成功签约流程（和甲方确认过签约成功后不支付下次需要重新签约）
// 使用指针返回并进行简单的nil判定能有效避免空数据导致的后续流程中断情况
func (*signService) UserLatestDetail(ctx context.Context, userId uint64, flowId string) (sign *model.Sign, err error) {
	sign = new(model.Sign)
	err = dao.Sign.Ctx(ctx).
		Where(dao.Sign.Columns.UserId, userId).
		Where(dao.Sign.Columns.FlowId, flowId).
		Where(dao.Sign.Columns.State, model.SignStateDone).
		OrderDesc(dao.Sign.Columns.Id).
		Limit(1).
		Scan(sign)
	return sign, err
}

func (*signService) Done(ctx context.Context, flowId string) error {
	_, err := dao.Sign.Ctx(ctx).Where(dao.Sign.Columns.FlowId, flowId).Update(g.Map{
		dao.Sign.Columns.State: model.SignStateDone,
	})
	return err
}
