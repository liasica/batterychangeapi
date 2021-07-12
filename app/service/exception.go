package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
)

var ExceptionService = exceptionService{}

type exceptionService struct {
}

func (*exceptionService) Create(ctx context.Context, req model.ExceptionReportReq) error {
	_, err := dao.Exception.Insert(req)
	return err
}
