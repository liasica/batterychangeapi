package service

import (
    "battery/app/dao"
    "battery/app/model"
    "context"
    "log"
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
        query = query.Where(c.ShopId, req.ShopId)
    }
    if !req.StartTime.IsZero() {
        log.Println(req.StartTime.IsZero(), req.StartTime.Format(layout))
        query = query.WhereGTE(c.CreatedAt, req.StartTime.Format(layout))
    }
    if !req.EndTime.IsZero() {
        query = query.WhereLTE(c.CreatedAt, req.EndTime.Format(layout))
    }

    _ = query.
        Page(req.PageIndex, req.PageLimit).
        LeftJoin("(SELECT cityId, d.name AS cityName, s.id AS extraShopId, s.name AS shopName FROM shop s LEFT JOIN districts d on s.cityId = d.id) extra ON extra.extraShopId = exception.shopId").
        Scan(&items)

    total, _ = query.Count()
    return
}
