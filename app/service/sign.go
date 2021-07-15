package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
)

var SignService = signService{}

type signService struct {
}

func (*signService) Create(ctx context.Context, data model.Sign) (uint64, error) {
	id, err := dao.Sign.Ctx(ctx).InsertAndGetId(data)
	return uint64(id), err
}
