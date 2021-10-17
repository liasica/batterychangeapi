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

func (*exceptionService) PageList(ctx context.Context, req *model.ExceptionListReq) (total int, items []model.ExceptionListItem) {
    c := dao.Exception.Columns
    layout := "2006-01-02"
    query := dao.Exception.Ctx(ctx)
    if req.ShopId > 0 {
        query.Where(c.ShopId, req.ShopId)
    }
    if !req.StartTime.IsZero() {
        query.WhereGTE(c.CreatedAt, req.StartTime.Format(layout))
    }
    if !req.EndTime.IsZero() {
        query.WhereLTE(c.CreatedAt, req.EndTime.Format(layout))
    }

    var rows []model.ExceptionListItem
    _ = query.
        WithAll().
        Page(req.PageIndex, req.PageLimit).
        Scan(&rows)

    total, _ = query.Count()
    return
}
