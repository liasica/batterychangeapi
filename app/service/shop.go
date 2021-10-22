package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/shop"
    "battery/library/snowflake"
    "context"
    "fmt"
    "github.com/gogf/gf/frame/g"
)

var ShopService = shopService{}

type shopService struct {
}

// ListUser 骑手获取门店列表
func (s *shopService) ListUser(ctx context.Context, req model.ShopListUserReq) (rep []model.ShopListUserRep) {
    c := dao.Shop.Columns
    m := dao.Shop.Ctx(ctx).Fields(
        c.Id,
        c.Name,
        c.Lat,
        c.Lng,
        c.Mobile,
        c.Address,
        c.State,
        fmt.Sprintf("(%s + %s) as batteryTotal", c.V60, c.V72),
        fmt.Sprintf("ST_DISTANCE_SPHERE(POINT(lng, lat), POINT(%f, %f)) as distance", req.Lng, req.Lat)).
        Where(dao.Shop.Columns.CityId, req.CityId).
        OrderAsc("distance").
        Page(req.PageIndex, req.PageLimit)
    if req.Name != "" {
        m = m.WhereLike(dao.Shop.Columns.Name, fmt.Sprintf("%%%s%%", req.Name))
    }
    _ = m.Scan(&rep)
    return
}

// GetShop 获取门店
func (s *shopService) GetShop(ctx context.Context, id uint) (rep *model.Shop, err error) {
    rep = new(model.Shop)
    err = dao.Shop.Ctx(ctx).WherePri(id).Scan(rep)
    return
}

// DetailByQr 获取门店详情
func (s *shopService) DetailByQr(ctx context.Context, qr string) (rep model.Shop, err error) {
    err = dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Qr, qr).Limit(1).Scan(&rep)
    return
}

// State 修改门店状态
func (s *shopService) State(ctx context.Context, id uint, state uint) error {
    _, err := dao.Shop.Ctx(ctx).Where(dao.Shop.Columns.Id, id).Update(g.Map{dao.Shop.Columns.State: state})
    return err
}

// MapIdName 获取店名IDMap
func (*shopService) MapIdName(ctx context.Context, ids []uint) map[uint]string {
    var list []model.Shop
    rep := map[uint]string{}
    _ = dao.Shop.Ctx(ctx).WhereIn(dao.Shop.Columns.Id, ids).Fields(dao.Shop.Columns.Id, dao.Shop.Columns.Name).Scan(&list)
    for _, shop := range list {
        rep[shop.Id] = shop.Name
    }
    return rep
}

// CheckMobile 检测电话是否可用
func (*shopService) CheckMobile(ctx context.Context, shopId uint, mobile string) bool {
    var cnt int
    c := dao.Shop.Columns
    query := dao.Shop.Ctx(ctx).Where(c.Mobile, mobile)
    if shopId == 0 {
        cnt, _ = query.Count()
    } else {
        cnt, _ = query.WhereNot(c.Id, shopId).Count()
    }
    return cnt == 0
}

// CheckName 检测名称是否可用
func (*shopService) CheckName(ctx context.Context, shopId uint, name string) bool {
    c := dao.Shop.Columns
    query := dao.Shop.Ctx(ctx).Where(c.Name, name)
    if shopId > 0 {
        query = query.WhereNot(c.Id, shopId)
    }
    cnt, _ := query.Count()
    return cnt == 0
}

// Create 创建门店
func (*shopService) Create(ctx context.Context, shop *model.Shop) error {
    shop.Qr = fmt.Sprintf("%d", snowflake.Service().Generate())
    id, err := dao.Shop.Ctx(ctx).InsertAndGetId(shop)
    shop.Id = uint(id)
    return err
}

// ListAdmin 管理员获取门店列表
func (s *shopService) ListAdmin(ctx context.Context, req model.ShopListAdminReq) (total int, items []model.ShopListItem) {
    m := dao.Shop.Ctx(ctx).Page(req.PageIndex, req.PageLimit)
    if req.Name != "" {
        m = m.WhereLike(fmt.Sprintf("%s.%s", shop.Table, dao.Shop.Columns.Name), fmt.Sprintf("%%%s%%", req.Name))
    }
    total, _ = m.Count()
    if total > 0 {
        _ = m.LeftJoin(`(SELECT shopId, COUNT(1) AS exceptionCnt FROM exception WHERE exception.state = 0 GROUP BY exception.shopId) exceptions ON exceptions.shopId = shop.id`).
            Scan(&items)
    }
    return
}
