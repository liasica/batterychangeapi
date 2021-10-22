package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/shop"
    "battery/app/model/shop_battery_record"
    "battery/library/sutil"
    "context"
    "errors"
    "fmt"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/os/gtime"
)

var ShopBatteryRecordService = shopBatteryRecordService{}

type shopBatteryRecordService struct{}

// Transfer 出入库
func (s *shopBatteryRecordService) Transfer(ctx context.Context, record model.ShopBatteryRecord, shopModel *model.Shop) error {
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

// Allocate 转移库存
func (s *shopBatteryRecordService) Allocate(ctx context.Context, req *model.BatteryAllocateReq) error {
    // 开始转移库存
    return dao.ShopBatteryRecord.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
        var err error
        st := shop.Table
        now := gtime.Now()

        sysUser := ctx.Value(model.ContextAdminKey).(*model.ContextAdmin)
        // 从门店转移
        if req.From > 0 {
            fromShop, _ := ShopService.GetShop(ctx, req.From)
            // 判断库存是否足够转移
            if fromShop == nil || sutil.StructGetFieldByString(fromShop, "V"+req.BatteryType).(int) < req.Num {
                return errors.New("库存不足")
            }
            // 减数量
            if _, err = tx.Update(st, fmt.Sprintf("`V%s` = `V%s` - %d", req.BatteryType, req.BatteryType, req.Num), "id = ?", req.From); err != nil {
                return err
            }
            // 出库记录
            if _, err = tx.Save(shop_battery_record.Table, model.ShopBatteryRecord{
                ShopId:      req.From,
                BatteryType: req.BatteryType,
                Type:        model.ShopBatteryRecordTypeOut,
                Num:         req.Num,
                Date:        now,
                SysUserId:   sysUser.Id,
                SysUserName: sysUser.Username,
            }); err != nil {
                return err
            }
        }

        // 转移到门店
        if req.To > 0 {
            // 加数量
            if _, err = tx.Update(st, fmt.Sprintf("`V%s` = `V%s` + %d", req.BatteryType, req.BatteryType, req.Num), "id = ?", req.To); err != nil {
                return err
            }

            // 入库记录
            if _, err = tx.Save(shop_battery_record.Table, model.ShopBatteryRecord{
                ShopId:      req.To,
                BatteryType: req.BatteryType,
                Type:        model.ShopBatteryRecordTypeIn,
                Num:         req.Num,
                Date:        now,
                SysUserId:   sysUser.Id,
                SysUserName: sysUser.Username,
            }); err != nil {
                return err
            }
        }

        return err
    })
}

// GetBatteryNumber 获取电池数量
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
