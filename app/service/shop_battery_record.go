package service

import (
    "battery/app/dao"
    "battery/app/model"
    "context"
    "github.com/gogf/gf/os/gtime"
)

var ShopBatteryRecordService = shopBatteryRecordService{}

type shopBatteryRecordService struct {
}

// DriverBiz 骑手业务调拨记录
func (*shopBatteryRecordService) DriverBiz(ctx context.Context, recordType, bizType, shopId uint, bizId uint64, user model.User) error {
    c := dao.ShopBatteryRecord.Columns
    _, err := dao.ShopBatteryRecord.Ctx(ctx).
        Fields(
            c.ShopId,
            c.BizId,
            c.BizType,
            c.UserName,
            c.BatteryType,
            c.Num,
            c.Type,
            c.UserId,
        ).
        Insert(model.ShopBatteryRecord{
            ShopId:      shopId,
            BizId:       bizId,
            BizType:     bizType,
            UserName:    user.RealName,
            BatteryType: user.BatteryType,
            Num:         1,
            Type:        recordType,
            UserId:      user.Id,
        })
    return err
}

// Platform 平台调拨
func (*shopBatteryRecordService) Platform(ctx context.Context, recordType, shopId, num uint, batteryType string) error {
    c := dao.ShopBatteryRecord.Columns
    sysUser := ctx.Value(model.ContextAdminKey).(*model.ContextAdmin)
    _, err := dao.ShopBatteryRecord.Ctx(ctx).
        Fields(
            c.ShopId,
            c.Num,
            c.Type,
            c.BatteryType,
            c.SysUserId,
            c.SysUserName,
        ).
        Insert(model.ShopBatteryRecord{
            ShopId:      shopId,
            BatteryType: batteryType,
            Num:         num,
            Type:        recordType,
            SysUserId:   sysUser.Id,
            SysUserName: sysUser.Username,
        })
    return err
}

// ShopList 门店获取电池记录
func (*shopBatteryRecordService) ShopList(ctx context.Context, shopId uint, recordType uint, st *gtime.Time, et *gtime.Time) (list []model.ShopBatteryRecord) {
    c := dao.ShopBatteryRecord.Columns
    layout := "Y-m-d"
    m := dao.ShopBatteryRecord.Ctx(ctx).
        Where(c.ShopId, shopId).
        Where(c.Type, recordType).
        OrderDesc(c.Id)
    if !st.IsZero() {
        m = m.WhereGTE(c.CreatedAt, st.Format(layout))
    }
    if !et.IsZero() {
        m = m.WhereLTE(c.CreatedAt, et.Format(layout))
    }
    _ = m.Scan(&list)
    return
}

// ShopDaysTotal 门店获取电池记录按天统计
func (*shopBatteryRecordService) ShopDaysTotal(ctx context.Context, days []int, recordType uint) (list []struct {
    Day int
    Cnt uint
}) {
    c := dao.ShopBatteryRecord.Columns
    _ = dao.ShopBatteryRecord.Ctx(ctx).
        Fields(c.Day, "count(*) cnt").
        WhereIn(c.Day, days).
        Where(c.Type, recordType).
        Group(c.Day).
        Scan(&list)

    return
}

// ListAdmin 所有门店电池记录
func (*shopBatteryRecordService) ListAdmin(ctx context.Context, req *model.BatteryRecordListReq) (total int, items []model.BatteryRecordListItem) {
    layout := "Y-m-d"
    c := dao.ShopBatteryRecord.Columns
    query := dao.ShopBatteryRecord.Ctx(ctx).
        OrderDesc(c.CreatedAt)
    if req.Type > 0 {
        query = query.Where(c.Type, req.Type)
    }
    if !req.StartDate.IsZero() {
        query = query.WhereGTE(c.CreatedAt, req.StartDate.Format(layout))
    }
    if !req.EndDate.IsZero() {
        query = query.WhereLTE(c.CreatedAt, req.EndDate.Format(layout))
    }
    if req.ShopId > 0 {
        query = query.Where(c.ShopId, req.ShopId)
    }
    _ = query.Page(req.PageIndex, req.PageLimit).Scan(&items)
    total, _ = query.Count()

    return
}
