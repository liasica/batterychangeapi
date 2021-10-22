package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/combo_order"
    "battery/app/model/user"
    "battery/library/mq"
    "battery/library/snowflake"
    "context"
    "errors"
    "fmt"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/os/gtime"
)

var ComboOrderService = comboOrderService{}

type comboOrderService struct {
}

// GenerateOrderNo 获取新的订单号
func (s comboOrderService) GenerateOrderNo() string {
    id := snowflake.Service().Generate()
    now := gtime.Now()
    return fmt.Sprintf("%d%d%d%012d", now.Year(), now.Month(), now.Day(), id%1000000000000)
}

// Detail 套餐订单详情
func (s *comboOrderService) Detail(ctx context.Context, id uint64) (rep model.ComboOrder, err error) {
    err = dao.ComboOrder.Ctx(ctx).WherePri(id).Scan(&rep)
    return
}

// DetailByNo 套餐订单详情
func (s *comboOrderService) DetailByNo(ctx context.Context, no string) (rep model.ComboOrder, err error) {
    err = dao.ComboOrder.Ctx(ctx).Where(dao.ComboOrder.Columns.No, no).Scan(&rep)
    return
}

// New 新购订单
func (s *comboOrderService) New(ctx context.Context, userId uint64, combo model.Combo) (order model.ComboOrder, err error) {
    no := s.GenerateOrderNo()
    c := dao.ComboOrder.Columns
    id, insertErr := dao.ComboOrder.Ctx(ctx).InsertAndGetId(g.Map{
        c.ComboId: combo.Id,
        c.Amount:  combo.Amount,
        c.Deposit: combo.Deposit,
        c.PayType: 0,
        c.No:      no,
        c.UserId:  userId,
        c.Type:    model.ComboTypeNew,
        c.CityId:  combo.CityId,
    })
    if insertErr == nil {
        err = dao.ComboOrder.Ctx(ctx).WherePri(id).Scan(&order)
    }
    return order, err
}

// Renewal 续购订单
func (s *comboOrderService) Renewal(ctx context.Context, payType uint, firstOrder model.ComboOrder) (order model.ComboOrder, err error) {
    no := s.GenerateOrderNo()
    c := dao.ComboOrder.Columns
    id, insertErr := dao.ComboOrder.Ctx(ctx).InsertAndGetId(g.Map{
        c.ComboId:  firstOrder.ComboId,
        c.Amount:   firstOrder.Amount - firstOrder.Deposit,
        c.Deposit:  0,
        c.PayType:  payType,
        c.No:       no,
        c.UserId:   firstOrder.UserId,
        c.ShopId:   firstOrder.ShopId,
        c.Type:     model.ComboTypeRenewal,
        c.CityId:   firstOrder.CityId,
        c.ParentId: firstOrder.Id,
    })
    if insertErr == nil {
        err = dao.ComboOrder.WherePri(id).Scan(&order)
    }
    return order, err
}

// Penalty 订单违约金
func (s *comboOrderService) Penalty(ctx context.Context, payType uint, amount float64, firstOrder model.ComboOrder) (order model.ComboOrder, err error) {
    no := s.GenerateOrderNo()
    c := dao.ComboOrder.Columns
    id, insertErr := dao.ComboOrder.Ctx(ctx).InsertAndGetId(g.Map{
        c.ComboId:  firstOrder.ComboId,
        c.Amount:   amount,
        c.Deposit:  0,
        c.PayType:  payType,
        c.No:       no,
        c.UserId:   firstOrder.UserId,
        c.ShopId:   firstOrder.ShopId,
        c.Type:     model.ComboTypePenalty,
        c.CityId:   firstOrder.CityId,
        c.ParentId: firstOrder.Id,
    })
    if insertErr == nil {
        err = dao.ComboOrder.WherePri(id).Scan(&order)
    }
    return order, err
}

// PaySuccess 订单支付成功处理
func (s *comboOrderService) PaySuccess(ctx context.Context, payAt *gtime.Time, no, PayPlatformNo string, payType uint) error {
    c := dao.ComboOrder.Columns
    _, err := dao.ComboOrder.Ctx(ctx).Where(dao.ComboOrder.Columns.No, no).Update(g.Map{
        c.PayState:      model.PayStateSuccess,
        c.PayPlatformNo: PayPlatformNo,
        c.PayAt:         payAt,
        c.PayType:       payType,
    })
    return err
}

// ShopClaim 门店认领订单
func (s *comboOrderService) ShopClaim(ctx context.Context, no string, shopId uint) error {
    res, err := dao.ComboOrder.Ctx(ctx).Where(dao.ComboOrder.Columns.No, no).Where(dao.ComboOrder.Columns.ShopId, 0).Update(g.Map{
        dao.ComboOrder.Columns.ShopId:     shopId,
        dao.ComboOrder.Columns.FirstUseAt: gtime.Now(),
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
func (s *comboOrderService) ShopMonthTotal(ctx context.Context, month, shopId uint, orderType uint) (res model.ShopOrderTotalRep) {
    m := dao.ComboOrder.Ctx(ctx).
        Fields("sum(amount) as amount, count(*) as cnt").
        Where(dao.ComboOrder.Columns.Month, month).
        Where(dao.ComboOrder.Columns.PayState, model.PayStateSuccess).
        Where(dao.ComboOrder.Columns.ShopId, shopId)
    if orderType != 0 {
        m = m.Where(dao.ComboOrder.Columns.Type, orderType)
    }
    _ = m.Scan(&res)
    return
}

// ShopMonthList 门店订单月度列表
func (s *comboOrderService) ShopMonthList(ctx context.Context, shopId uint, filter model.ShopOrderListReq) (res []model.ComboOrder) {
    m := dao.ComboOrder.Ctx(ctx).
        Where(dao.ComboOrder.Columns.ShopId, shopId).
        Where(dao.ComboOrder.Columns.PayState, model.PayStateSuccess).
        OrderDesc(dao.ComboOrder.Columns.Id).
        Page(filter.PageIndex, filter.PageLimit)

    if filter.Keywords != "" {
        var users []struct {
            Id uint64
        }
        _ = dao.User.Where(dao.User.Columns.Mobile, filter.Keywords).
            WhereOrLike(dao.User.Columns.RealName, fmt.Sprintf("%%%s%%", filter.Keywords)).
            Fields(dao.User.Columns.Id).Scan(&users)
        if len(users) == 0 {
            return []model.ComboOrder{}
        }
        userIds := make([]uint64, len(users))
        for i, u := range users {
            userIds[i] = u.Id
        }
        m = m.WhereIn(dao.ComboOrder.Columns.UserId, userIds)
    }

    if filter.Month != 0 {
        m = m.Where(dao.ComboOrder.Columns.Month, filter.Month)
    }
    if filter.Type != 0 {
        m = m.Where(dao.ComboOrder.Columns.Type, filter.Type)
    }

    _ = m.Scan(&res)
    return
}

// ListAdmin 订单列表查询
func (*comboOrderService) ListAdmin(ctx context.Context, req *model.OrderListReq) (total int, items []model.OrderListItem) {
    query := dao.ComboOrder.Ctx(ctx)
    c := dao.ComboOrder.Columns
    t := combo_order.Table
    layout := "Y-m-d"

    params := mq.ParseStructToQuery(*req, "RealName", "Mobile")
    query = query.Where(params)

    if !req.StartDate.IsZero() {
        query = query.WhereGTE(fmt.Sprintf("%s.%s", t, c.CreatedAt), req.StartDate.Format(layout))
    }
    if !req.EndDate.IsZero() {
        query = query.WhereLTE(fmt.Sprintf("%s.%s", t, c.CreatedAt), req.EndDate.Format(layout))
    }

    if req.Mobile != "" {
        query = query.LeftJoin(user.Table, fmt.Sprintf("%s.%s=%s.%s", user.Table, dao.User.Columns.Id, combo_order.Table, c.UserId)).
            Where(fmt.Sprintf("%s.%s = ?", user.Table, dao.User.Columns.Mobile), req.Mobile)
    }

    if req.RealName != "" {
        query = query.LeftJoin(user.Table, fmt.Sprintf("%s.%s=%s.%s", user.Table, dao.User.Columns.Id, combo_order.Table, c.UserId)).
            WhereLike(fmt.Sprintf("%s.%s", user.Table, dao.User.Columns.RealName), "%"+req.RealName+"%")
    }

    fields := mq.FieldsWithTable(combo_order.Table, c)

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

        if row.ComboDetail != nil {
            pn = row.ComboDetail.Name
        }

        items = append(items, model.OrderListItem{
            Id:         row.Id,
            No:         row.No,
            ShopId:     row.ShopId,
            UserId:     row.UserId,
            RealName:   row.User.RealName,
            Mobile:     row.User.Mobile,
            Type:       row.Type,
            ShopName:   sn,
            ComboName:  pn,
            CityName:   row.City.Name,
            Amount:     row.Amount,
            Deposit:    row.Deposit,
            PayType:    row.PayType,
            PayState:   row.PayState,
            PayAt:      row.PayAt,
            FirstUseAt: row.FirstUseAt,
            CreatedAt:  row.CreatedAt,
        })
    }

    total, _ = query.Count()
    return
}
