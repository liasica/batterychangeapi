package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/user"
    "battery/app/model/user_biz"
    "battery/library/mq"
    "context"
    "fmt"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/os/gtime"
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
    u := ctx.Value(model.ContextRiderKey).(*model.ContextRider)
    _ = dao.UserBiz.Ctx(ctx).
        Where(dao.UserBiz.Columns.UserId, u.Id).
        WhereIn(dao.UserBiz.Columns.BizType, []int{model.UserBizNew, model.UserBizBatteryRenewal, model.UserBizBatteryRecover}).
        OrderDesc(dao.UserBiz.Columns.Id).
        Page(req.PageIndex, req.PageLimit).
        Scan(&rep)
    return
}

// ShopFilter 门店业务记录
func (b *userBizService) ShopFilter(ctx context.Context, req *model.BizShopFilterReq) (total int, items []model.BizShopFilterResp) {
    layout := "Y-m"

    filterReq := &model.BizListReq{
        Page:     req.Page,
        ShopId:   req.ShopId,
        RealName: req.RealName,
        UserType: req.UserType,
        BizType:  req.BizType,
    }
    if req.Month != "" {
        start := gtime.NewFromStr(req.Month + "-01")
        end := start.EndOfMonth()
        filterReq.StartDate = start
        filterReq.EndDate = end
    }
    total, rows := b.Filter(ctx, filterReq)

    // 组装数据
    tmp := make(map[string]model.BizShopFilterResp)
    for k, row := range rows {
        month := row.CreatedAt.Format(layout)
        item, ok := tmp[month]
        if !ok {
            item = model.BizShopFilterResp{
                Month: month,
                Total: 0,
            }
        }
        item.Total++
        item.Items = append(item.Items, &model.BizShopFilterItem{
            RealName:  row.RealName,
            ComboName: row.ComboName,
            GroupName: row.GroupName,
            Mobile:    row.Mobile,
            BizType:   row.BizType,
            CreatedAt: row.CreatedAt,
        })
        if k == len(rows)-1 {
            // 最后一项所在月份的数量
            if total > 0 && total != len(rows) {
                last := rows[len(rows)-1]
                filterReq.StartDate = last.CreatedAt.StartOfMonth()
                filterReq.EndDate = last.CreatedAt.EndOfMonth()
                r, _ := b.FilterQuery(ctx, filterReq).All()
                item.Total = r.Size()
            }
        }
        tmp[month] = item
    }

    for _, item := range tmp {
        items = append(items, item)
    }

    return
}

// FilterQuery 获取查询命令
func (b *userBizService) FilterQuery(ctx context.Context, req *model.BizListReq) (query *gdb.Model) {
    layout := "Y-m-d"
    c := dao.UserBiz.Columns
    t := user_biz.Table
    params := mq.ParseStructToQuery(*req, "RealName", "Mobile")

    query = dao.UserBiz.Ctx(ctx).OrderDesc(c.CreatedAt).Where(params)

    if req.ShopId > 0 {
        query = query.Where(fmt.Sprintf("%s.%s", t, c.ShopId), req.ShopId)
    }

    switch req.UserType {
    case model.UserTypePersonal:
        query = query.Where(fmt.Sprintf("%s.%s", t, c.GroupId), 0)
    case model.UserTypeGroupMember:
        query = query.WhereGT(fmt.Sprintf("%s.%s", t, c.GroupId), 0)
    }

    if req.BizType > 0 {
        query = query.Where(fmt.Sprintf("%s.%s", t, c.BizType), req.BizType)
    }

    if !req.StartDate.IsZero() {
        query = query.WhereGTE(fmt.Sprintf("%s.%s", t, c.CreatedAt), req.StartDate.Format(layout))
    }
    if !req.EndDate.IsZero() {
        query = query.WhereLTE(fmt.Sprintf("%s.%s", t, c.CreatedAt), req.EndDate.Format(layout))
    }

    if req.Mobile != "" {
        query = query.LeftJoin(fmt.Sprintf("%s ut", user.Table), fmt.Sprintf("ut.%s=%s.%s", dao.User.Columns.Id, user_biz.Table, c.UserId)).
            Where(fmt.Sprintf("ut.%s = ?", dao.User.Columns.Mobile), req.Mobile)
    }

    if req.RealName != "" {
        query = query.LeftJoin(fmt.Sprintf("%s ut", user.Table), fmt.Sprintf("ut.%s=%s.%s", dao.User.Columns.Id, user_biz.Table, c.UserId)).
            WhereLike(fmt.Sprintf("ut.%s", dao.User.Columns.RealName), "%"+req.RealName+"%")
    }

    query = query.Fields(mq.FieldsWithTable(user_biz.Table, c))
    return
}

// Filter 业务列表查询
func (b *userBizService) Filter(ctx context.Context, req *model.BizListReq) (total int, items []model.BizEntity) {
    query := b.FilterQuery(ctx, req)

    _ = query.WithAll().
        Page(req.PageIndex, req.PageLimit).
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

    r, _ := query.All()
    total = r.Size()
    return
}

func (*userBizService) ListSimaple(ctx context.Context, member *model.User, page *model.Page) (total int, items []model.BizSimpleItem) {
    c := dao.UserBiz.Columns
    query := dao.UserBiz.Ctx(ctx).
        WithAll().
        Where(c.GroupId, member.GroupId).
        Where(c.UserId, member.Id).
        Page(page.PageIndex, page.PageLimit).
        OrderDesc(c.CreatedAt)

    _ = query.Scan(&items)

    for k, item := range items {
        if item.Shop != nil {
            items[k].ShopName = item.Shop.Name
        }
    }

    total, _ = query.Count()
    return
}
