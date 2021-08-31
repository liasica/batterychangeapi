package service

import (
	"context"
	"github.com/gogf/gf/frame/g"

	"battery/app/dao"
	"battery/app/model"
)

var SignService = signService{}

type signService struct {
}

func (*signService) Create(ctx context.Context, data model.Sign) (uint64, error) {
	id, err := dao.Sign.Ctx(ctx).InsertAndGetId(data)
	return uint64(id), err
}

// GetDetailBayFlowId 根据flowId获取签约信息
func (*signService) GetDetailBayFlowId(ctx context.Context, flowId string) (model.Sign, error) {
	var sign model.Sign
	err := dao.Sign.Ctx(ctx).Where(dao.Sign.Columns.FlowId, flowId).Limit(1).Scan(&sign)
	return sign, err
}

// GetDetailBayFileId 根据fileId获取签约信息
func (*signService) GetDetailBayFileId(ctx context.Context, fileId string) (model.Sign, error) {
	var sign model.Sign
	err := dao.Sign.Ctx(ctx).Where(dao.Sign.Columns.FileId, fileId).Limit(1).Scan(&sign)
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

// UserLatestDoneDetail 用户最后完成的签约并支付成功的签约信息
func (*signService) UserLatestDoneDetail(ctx context.Context, userId, packagesOrderId uint64, groupId uint) (sign *model.Sign, err error) {
	err = dao.Sign.Ctx(ctx).
		Where(dao.Sign.Columns.UserId, userId).
		Where(dao.Sign.Columns.PackagesOrderId, packagesOrderId).
		Where(dao.Sign.Columns.GroupId, groupId).
		Where(dao.Sign.Columns.State, model.SignStateDone).
		OrderDesc(dao.Sign.Columns.Id).
		Limit(1).
		Scan(sign)
	return sign, err
}
