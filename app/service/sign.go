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

func (*signService) UserLatestDetail(ctx context.Context, userId uint64) (model.Sign, error) {
	var sign model.Sign
	err := dao.Sign.Ctx(ctx).Where(dao.Sign.Columns.UserId, userId).OrderDesc(dao.Sign.Columns.Id).Limit(1).Scan(&sign)
	return sign, err
}

func (*signService) Done(ctx context.Context, flowId string) error {
	_, err := dao.Sign.Ctx(ctx).Where(dao.Sign.Columns.FlowId, flowId).Update(g.Map{
		dao.Sign.Columns.State: model.SignStateDone,
	})
	return err
}
