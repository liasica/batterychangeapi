package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/shop"
    "battery/app/model/shop_battery_record"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/os/gtime"
)

var ShopBatteryRecordService = shopBatteryRecordService{}

type shopBatteryRecordService struct{}

// Transfer 出入库
func (s shopBatteryRecordService) Transfer(ctx context.Context, record model.ShopBatteryRecord, shopModel *model.Shop) error {
    return dao.ShopBatteryRecord.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
        var err error
        // 数量 正数为入库 负数为出库
        num := record.Num
        if record.Type == model.ShopBatteryRecordTypeOut {
            num *= -1
        }
        // 记录入库
        record.Date = gtime.Now()
        if _, err = tx.Save(shop_battery_record.Table, record); err != nil {
            return err
        }
        // 操作店铺库存
        switch record.BatteryType {
        case model.BatteryType60:
            shopModel.V60 += num
        case model.BatteryType72:
            shopModel.V72 += num
        }
        _, err = tx.Save(shop.Table, shopModel)
        return err
    })
}

// DriverBiz 骑手业务调拨记录
func (*shopBatteryRecordService) DriverBiz(ctx context.Context, recordType, bizType, shopId uint, bizId uint64, user model.User) error {
    _, err := dao.ShopBatteryRecord.Ctx(ctx).
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

func (*shopBatteryRecordService) GetBatteryNumber(ctx context.Context, shopId uint) (data model.ShopBatteryRecordStatRep) {
    var items []*model.ShopBatteryRecord
    c := dao.ShopBatteryRecord.Columns
    _ = dao.ShopBatteryRecord.Ctx(ctx).
        Where(c.ShopId, shopId).
        Scan(&items)

    for _, item := range items {
        if item.Type == model.ShopBatteryRecordTypeIn {
        }
        switch item.Type {
        case model.ShopBatteryRecordTypeIn:
            data.InTotal += item.Num
        case model.ShopBatteryRecordTypeOut:
            data.OutTotal += item.Num
        }
    }

    return
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
