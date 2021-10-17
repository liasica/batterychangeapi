package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/packages_order"
    "battery/app/model/user"
    "battery/library/request"
    "battery/library/snowflake"
    "context"
    "errors"
    "fmt"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/os/gtime"
    "github.com/gogf/gf/util/gutil"
)

var PackagesOrderService = packagesOrderService{}

type packagesOrderService struct {
}

// GenerateOrderNo 获取新的订单号
func (s packagesOrderService) GenerateOrderNo() string {
    id := snowflake.Service().Generate()
    now := gtime.Now()
    return fmt.Sprintf("%d%d%d%012d", now.Year(), now.Month(), now.Day(), id%1000000000000)
}

// Detail 套餐订单详情
func (s *packagesOrderService) Detail(ctx context.Context, id uint64) (rep model.PackagesOrder, err error) {
    err = dao.PackagesOrder.Ctx(ctx).WherePri(id).Scan(&rep)
    return
}

// DetailByNo 套餐订单详情
func (s *packagesOrderService) DetailByNo(ctx context.Context, no string) (rep model.PackagesOrder, err error) {
    err = dao.PackagesOrder.Ctx(ctx).Where(dao.PackagesOrder.Columns.No, no).Scan(&rep)
    return
}

// New 新购订单
func (s *packagesOrderService) New(ctx context.Context, userId uint64, packages model.Packages) (order model.PackagesOrder, err error) {
    no := s.GenerateOrderNo()
    c := dao.PackagesOrder.Columns
    id, insertErr := dao.PackagesOrder.Ctx(ctx).InsertAndGetId(g.Map{
        c.PackageId: packages.Id,
        c.Amount:    packages.Amount,
        c.Earnest:   packages.Earnest,
        c.PayType:   0,
        c.No:        no,
        c.UserId:    userId,
        c.Type:      model.PackageTypeNew,
        c.CityId:    packages.CityId,
    })
    if insertErr == nil {
        err = dao.PackagesOrder.Ctx(ctx).WherePri(id).Scan(&order)
    }
    return order, err
}

// Renewal 续购订单
func (s *packagesOrderService) Renewal(ctx context.Context, payType uint, firstOrder model.PackagesOrder) (order model.PackagesOrder, err error) {
    no := s.GenerateOrderNo()
    c := dao.PackagesOrder.Columns
    id, insertErr := dao.PackagesOrder.Ctx(ctx).InsertAndGetId(g.Map{
        c.PackageId: firstOrder.PackageId,
        c.Amount:    firstOrder.Amount - firstOrder.Earnest,
        c.Earnest:   0,
        c.PayType:   payType,
        c.No:        no,
        c.UserId:    firstOrder.UserId,
        c.ShopId:    firstOrder.ShopId,
        c.Type:      model.PackageTypeRenewal,
        c.CityId:    firstOrder.CityId,
        c.ParentId:  firstOrder.Id,
    })
    if insertErr == nil {
        err = dao.PackagesOrder.WherePri(id).Scan(&order)
    }
    return order, err
}

// Penalty 订单违约金
func (s *packagesOrderService) Penalty(ctx context.Context, payType uint, amount float64, firstOrder model.PackagesOrder) (order model.PackagesOrder, err error) {
    no := s.GenerateOrderNo()
    c := dao.PackagesOrder.Columns
    id, insertErr := dao.PackagesOrder.Ctx(ctx).InsertAndGetId(g.Map{
        c.PackageId: firstOrder.PackageId,
        c.Amount:    amount,
        c.Earnest:   0,
        c.PayType:   payType,
        c.No:        no,
        c.UserId:    firstOrder.UserId,
        c.ShopId:    firstOrder.ShopId,
        c.Type:      model.PackageTypePenalty,
        c.CityId:    firstOrder.CityId,
        c.ParentId:  firstOrder.Id,
    })
    if insertErr == nil {
        err = dao.PackagesOrder.WherePri(id).Scan(&order)
    }
    return order, err
}

// PaySuccess 订单支付成功处理
func (s *packagesOrderService) PaySuccess(ctx context.Context, payAt *gtime.Time, no, PayPlatformNo string, payType uint) error {
    c := dao.PackagesOrder.Columns
    _, err := dao.PackagesOrder.Ctx(ctx).Where(dao.PackagesOrder.Columns.No, no).Update(g.Map{
        c.PayState:      model.PayStateSuccess,
        c.PayPlatformNo: PayPlatformNo,
        c.PayAt:         payAt,
        c.PayType:       payType,
    })
    return err
}

// ShopClaim 门店认领订单
func (s *packagesOrderService) ShopClaim(ctx context.Context, no string, shopId uint) error {
    res, err := dao.PackagesOrder.Ctx(ctx).Where(dao.PackagesOrder.Columns.No, no).Where(dao.PackagesOrder.Columns.ShopId, 0).Update(g.Map{
        dao.PackagesOrder.Columns.ShopId:     shopId,
        dao.PackagesOrder.Columns.FirstUseAt: gtime.Now(),
    })
    if err == nil {
        rows, err := res.RowsAffected()
        if rows > 0 && err == nil {
            return nil
        }
        err = errors.New("认领失败")
    }
    return err
}

// ShopMonthTotal 门店订单月统计
func (s *packagesOrderService) ShopMonthTotal(ctx context.Context, month, shopId uint, orderType uint) (res model.ShopOrderTotalRep) {
    m := dao.PackagesOrder.Ctx(ctx).
        Fields("sum(amount) as amount, count(*) as cnt").
        Where(dao.PackagesOrder.Columns.Month, month).
        Where(dao.PackagesOrder.Columns.PayState, model.PayStateSuccess).
        Where(dao.PackagesOrder.Columns.ShopId, shopId)
    if orderType != 0 {
        m = m.Where(dao.PackagesOrder.Columns.Type, orderType)
    }
    _ = m.Scan(&res)
    return
}

// ShopMonthList 门店订单月度列表
func (s *packagesOrderService) ShopMonthList(ctx context.Context, shopId uint, filter model.ShopOrderListReq) (res []model.PackagesOrder) {
    m := dao.PackagesOrder.Ctx(ctx).
        Where(dao.PackagesOrder.Columns.ShopId, shopId).
        Where(dao.PackagesOrder.Columns.PayState, model.PayStateSuccess).
        OrderDesc(dao.PackagesOrder.Columns.Id).
        Page(filter.PageIndex, filter.PageLimit)

    if filter.Keywords != "" {
        var users []struct {
            Id uint64
        }
        _ = dao.User.Where(dao.User.Columns.Mobile, filter.Keywords).
            WhereOrLike(dao.User.Columns.RealName, fmt.Sprintf("%%%s%%", filter.Keywords)).
            Fields(dao.User.Columns.Id).Scan(&users)
        if len(users) == 0 {
            return []model.PackagesOrder{}
        }
        userIds := make([]uint64, len(users))
        for i, u := range users {
            userIds[i] = u.Id
        }
        m = m.WhereIn(dao.PackagesOrder.Columns.UserId, userIds)
    }

    if filter.Month != 0 {
        m = m.Where(dao.PackagesOrder.Columns.Month, filter.Month)
    }
    if filter.Type != 0 {
        m = m.Where(dao.PackagesOrder.Columns.Type, filter.Type)
    }

    _ = m.Scan(&res)
    return
}

func (*packagesOrderService) ListAdmin(ctx context.Context, req *model.OrderListReq) (total int, items []model.OrderListItem) {
    query := dao.PackagesOrder.Ctx(ctx)
    c := dao.PackagesOrder.Columns
    layout := "Y-m-d"

    params := request.ParseStructToQuery(*req)
    delete(params, "RealName")
    delete(params, "Mobile")
    query = query.Where(params)

    if !req.StartDate.IsZero() {
        query = query.WhereGTE(c.CreatedAt, req.StartDate.Format(layout))
    }
    if !req.EndDate.IsZero() {
        query = query.WhereLTE(c.CreatedAt, req.EndDate.Format(layout))
    }

    if req.Mobile != "" {
        query = query.LeftJoin(user.Table, fmt.Sprintf("%s.%s=%s.%s", user.Table, dao.User.Columns.Id, packages_order.Table, c.UserId)).
            Where(fmt.Sprintf("%s.%s = ?", user.Table, dao.User.Columns.Mobile), req.Mobile)
    }

    if req.RealName != "" {
        query = query.LeftJoin(user.Table, fmt.Sprintf("%s.%s=%s.%s", user.Table, dao.User.Columns.Id, packages_order.Table, c.UserId)).
            WhereLike(fmt.Sprintf("%s.%s", user.Table, dao.User.Columns.RealName), "%"+req.RealName+"%")
    }

    t := packages_order.Table
    var fields []string
    keys := gutil.Values(c)
    for _, field := range keys {
        fields = append(fields, fmt.Sprintf("%s.%v", t, field))
    }

    var rows []model.OrderEntity
    _ = query.WithAll().
        Page(req.PageIndex, req.PageLimit).
        OrderDesc(c.CreatedAt).
        Fields(fields).
        Scan(&rows)

    for _, row := range rows {
        var sn, pn string

        if row.Shop != nil {
            sn = row.Shop.Name
        }

        if row.PackageDetail != nil {
            pn = row.PackageDetail.Name
        }

        items = append(items, model.OrderListItem{
            Id:          row.Id,
            No:          row.No,
            ShopId:      row.ShopId,
            UserId:      row.UserId,
            RealName:    row.User.RealName,
            Mobile:      row.User.Mobile,
            Type:        row.Type,
            ShopName:    sn,
            PackageName: pn,
            CityName:    row.City.Name,
            Amount:      row.Amount,
            Earnest:     row.Earnest,
            PayType:     row.PayType,
            PayState:    row.PayState,
            PayAt:       row.PayAt,
            FirstUseAt:  row.FirstUseAt,
            CreatedAt:   row.CreatedAt,
        })
    }

    total, _ = query.Count()
    return
}
