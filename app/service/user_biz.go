package service

import (
    "battery/app/model/user"
    "battery/app/model/user_biz"
    "battery/library/mq"
    "context"
    "fmt"
    "github.com/gogf/gf/os/gtime"
    "time"

    "battery/app/dao"
    "battery/app/model"
)

var UserBizService = userBizService{}

type userBizService struct {
}

// Create 添加记录
func (*userBizService) Create(ctx context.Context, req model.UserBiz) (uint64, error) {
    at := gtime.Now()
    req.CreatedAt = at
    req.UpdatedAt = at
    id, err := dao.UserBiz.Ctx(ctx).InsertAndGetId(req)
    return uint64(id), err
}

// ListUser 骑手获取换电记录
func (*userBizService) ListUser(ctx context.Context, req model.Page) (rep []model.UserBiz) {
    user := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
    _ = dao.UserBiz.Ctx(ctx).
        Where(dao.UserBiz.Columns.UserId, user.Id).
        WhereIn(dao.UserBiz.Columns.Type, []int{model.UserBizNew, model.UserBizBatteryRenewal, model.UserBizBatteryUnSave}).
        OrderDesc(dao.UserBiz.Columns.Id).
        Page(req.PageIndex, req.PageLimit).
        Scan(&rep)
    return
}

// ListShop 门店获取换电记录
func (*userBizService) ListShop(ctx context.Context, req model.UserBizShopRecordReq) (rep []model.UserBiz) {
    manager := ctx.Value(model.ContextShopManagerKey).(*model.ContextShopManager)
    m := dao.UserBiz.Ctx(ctx).
        Where(dao.UserBiz.Columns.ShopId, manager.ShopId).
        OrderDesc(dao.UserBiz.Columns.Id).
        Page(req.PageIndex, req.PageLimit)
    if req.BizType == 0 {
        m = m.WhereIn(dao.UserBiz.Columns.Type, []uint{model.UserBizBatteryRenewal, model.UserBizBatterySave, model.UserBizClose})
    } else {
        m = m.Where(dao.UserBiz.Columns.Type, req.BizType)
    }

    if req.Keywords != "" {
        var users []struct {
            Id uint64
        }
        _ = dao.User.Where(dao.User.Columns.Mobile, req.Keywords).
            WhereOrLike(dao.User.Columns.RealName, fmt.Sprintf("%%%s%%", req.Keywords)).
            Fields(dao.User.Columns.Id).Scan(&users)

        if len(users) == 0 {
            return []model.UserBiz{}
        }
        userIds := make([]uint64, len(users))
        for i, u := range users {
            userIds[i] = u.Id
        }
        m = m.WhereIn(dao.UserBiz.Columns.UserId, userIds)
    }

    if req.Month > 0 {
        year := int(req.Month / 100)
        month := time.Month(req.Month % 100)
        firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
        m = m.WhereGTE(dao.UserBiz.Columns.CreatedAt, firstOfMonth)
        lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
        et, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%d-%02d-%02d 23:59:59", lastOfMonth.Year(), lastOfMonth.Month(), lastOfMonth.Day()))
        m = m.WhereLTE(dao.UserBiz.Columns.CreatedAt, et)
    }
    if req.UserType == 1 {
        m = m.Where(dao.UserBiz.Columns.GoroupId, 0)
    } else if req.UserType == 2 {
        m = m.WhereGT(dao.UserBiz.Columns.GoroupId, 0)
    }
    _ = m.Scan(&rep)
    return
}

// ListShopMonthTotal 门店获取换电记录按月统计
func (*userBizService) ListShopMonthTotal(ctx context.Context, req model.UserBizShopRecordMonthTotalReq) (rep model.UserBizShopRecordMonthTotalRep) {
    manager := ctx.Value(model.ContextShopManagerKey).(*model.ContextShopManager)
    m := dao.UserBiz.Ctx(ctx).
        Where(dao.UserBiz.Columns.ShopId, manager.ShopId).
        WhereIn(dao.UserBiz.Columns.Type, []uint{model.UserBizBatteryRenewal, model.UserBizBatterySave, model.UserBizClose}).
        OrderDesc(dao.UserBiz.Columns.Id)
    if req.UserType == 1 {
        m = m.Where(dao.UserBiz.Columns.GoroupId, 0)
    } else if req.UserType == 2 {
        m = m.WhereGT(dao.UserBiz.Columns.GoroupId, 0)
    }
    if req.Month > 0 {
        year := int(req.Month / 100)
        month := time.Month(req.Month % 100)
        firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
        m = m.WhereGTE(dao.UserBiz.Columns.CreatedAt, firstOfMonth)
        lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
        et, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%d-%02d-%02d 23:59:59", lastOfMonth.Year(), lastOfMonth.Month(), lastOfMonth.Day()))
        m = m.WhereLTE(dao.UserBiz.Columns.CreatedAt, et)
    }
    if req.BizType == 0 {
        m = m.WhereIn(dao.UserBiz.Columns.Type, []uint{model.UserBizBatteryRenewal, model.UserBizBatterySave, model.UserBizClose})
    } else {
        m = m.Where(dao.UserBiz.Columns.Type, req.BizType)
    }
    rep.Cnt, _ = m.Count()
    return
}

// ListAdmin 业务列表查询
func (*userBizService) ListAdmin(ctx context.Context, req *model.BizListReq) (total int, items []model.BizEntity) {
    query := dao.UserBiz.Ctx(ctx)

    c := dao.UserBiz.Columns
    layout := "Y-m-d"

    params := mq.ParseStructToQuery(*req, "RealName", "Mobile")
    query = query.Where(params)

    if !req.StartDate.IsZero() {
        query = query.WhereGTE(c.CreatedAt, req.StartDate.Format(layout))
    }
    if !req.EndDate.IsZero() {
        query = query.WhereLTE(c.CreatedAt, req.EndDate.Format(layout))
    }

    if req.Mobile != "" {
        query = query.LeftJoin(user.Table, fmt.Sprintf("%s.%s=%s.%s", user.Table, dao.User.Columns.Id, user_biz.Table, c.UserId)).
            Where(fmt.Sprintf("%s.%s = ?", user.Table, dao.User.Columns.Mobile), req.Mobile)
    }

    if req.RealName != "" {
        query = query.LeftJoin(user.Table, fmt.Sprintf("%s.%s=%s.%s", user.Table, dao.User.Columns.Id, user_biz.Table, c.UserId)).
            WhereLike(fmt.Sprintf("%s.%s", user.Table, dao.User.Columns.RealName), "%"+req.RealName+"%")
    }

    fields := mq.FieldsWithTable(user_biz.Table, c)

    _ = query.WithAll().
        Page(req.PageIndex, req.PageLimit).
        OrderDesc(c.CreatedAt).
        Fields(fields).
        Scan(&items)

    for k, row := range items {
        if row.Shop != nil {
            items[k].ShopName = row.Shop.Name
        }

        if row.ComboDetail != nil {
            items[k].ComboName = row.ComboDetail.Name
        }

        if row.Group != nil {
            items[k].GroupName = row.Group.Name
        }

        if row.City != nil {
            items[k].CityName = row.City.Name
        }

        if row.User != nil {
            items[k].RealName = row.User.RealName
            items[k].Mobile = row.User.Mobile
        }
    }

    total, _ = query.Count()
    return
}
