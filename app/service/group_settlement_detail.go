// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-19
// Based on apiv2 by liasica, magicrolan@qq.com.

package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/group_settlement"
    "battery/library/mongo"
    "context"
    "errors"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/os/gtime"
    "github.com/qiniu/qmgo/options"
    "github.com/shopspring/decimal"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var GroupSettlementDetailService = new(groupSettlementDetailService)

type groupSettlementDetailService struct {
}

// Earning 入账
func (s *groupSettlementDetailService) Earning(ctx context.Context, user model.User) error {
    now := gtime.Now()
    // 获取套餐详情
    var combo = new(model.Combo)
    _ = dao.Combo.Ctx(ctx).Where(dao.Combo.Columns.Id, user.ComboId).Scan(combo)
    if combo == nil {
        return errors.New("未找到有效套餐")
    }

    data := model.GroupSettlementDetail{
        UserId:       user.Id,
        State:        model.SettlementBilling,
        GroupId:      user.GroupId,
        ComboId:      user.ComboId,
        ComboOrderId: user.ComboOrderId,
        BatteryType:  user.BatteryType,
        UnitPrice:    combo.UnitPrice(),
        StartDate:    now,
    }

    return dao.GroupSettlementDetail.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {

        id, err := dao.GroupSettlementDetail.Ctx(ctx).
            Data(data).
            InsertAndGetId()
        if err != nil {
            return err
        }

        // 查询今日是否有多次新签操作，若有，则标记为忽略
        // 将今日其他的入账操作标记为忽略 (startDate = now)
        c := dao.GroupSettlementDetail.Columns
        if _, err := dao.GroupSettlementDetail.Ctx(ctx).
            Data(g.Map{c.Ignorance: 1}).
            Where(
                g.Map{
                    c.UserId:    data.UserId,
                    c.GroupId:   data.GroupId,
                    c.StartDate: now.Format("Y-m-d"),
                }).
            WhereNot(c.Id, id).
            Update(); err != nil {
            return err
        }
        return err
    })
}

// Cancel 退租
func (s *groupSettlementDetailService) Cancel(ctx context.Context, user model.User) error {
    now := gtime.Now()
    c := dao.GroupSettlementDetail.Columns
    detail := new(model.GroupSettlementDetail)
    err := dao.GroupSettlementDetail.Ctx(ctx).
        Where(c.UserId+" = ?", user.Id).
        OrderDesc(c.StartDate).
        Scan(detail)
    if err != nil {
        return err
    }

    return dao.GroupSettlementDetail.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
        detail.CancelDate = now
        _, err = dao.GroupSettlementDetail.Ctx(ctx).Data(detail).Save()
        if err != nil {
            return err
        }
        return err
    })
}

// Checkout 结算
// func (s *groupSettlementDetailService) Checkout(ctx context.Context, user model.DriverBiz) error {
//     now := gtime.Now()
//     detail := new(model.GroupSettlementDetail)
//     err := dao.GroupSettlementDetail.Ctx(ctx).Where("id = ?", user.SettlementDetailId).Scan(detail)
//     if err != nil {
//         return err
//     }
//
//     return dao.GroupSettlementDetail.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
//         detail.CancelDate = now
//         detail.State = model.SettlementSettled
//         _, err = dao.GroupSettlementDetail.Ctx(ctx).Data(detail).Save()
//         if err != nil {
//             return err
//         }
//
//         // 若用户是使用或逾期状态，则创建新的账单
//         switch user.BatteryState {
//         case model.BatteryStateUse, model.BatteryStateOverdue:
//             return s.Earning(ctx, user)
//         }
//         return nil
//     })
// }

// GetBillFromMongo 读取账单缓存
func (s *groupSettlementDetailService) GetBillFromMongo(ctx context.Context, groupId uint, date string) (bill *model.SettlementCache) {
    col := mongo.DB.Collection(group_settlement.Table)
    _ = col.Find(ctx, bson.M{"groupId": groupId, "toDate": date}).Sort("createAt").Limit(1).One(&bill)
    return
}

// GetGroupBill 获取团队账单
// expDate 截止日期
func (s *groupSettlementDetailService) GetGroupBill(ctx context.Context, group *model.Group, expDate *gtime.Time) (bill *model.SettlementCache, err error) {
    layout := "Y-m-d"
    // 尝试从缓存中读取账单
    bill = s.GetBillFromMongo(ctx, group.Id, expDate.Format(layout))
    if bill != nil {
        bill.Hash = bill.Id.Hex()
        return
    }

    // 重新生成账单
    var rows []*model.SettlementDetailEntity
    var items []*model.SettlementListItem
    var total float64
    var from *gtime.Time
    now := gtime.Now()
    c := dao.GroupSettlementDetail.Columns
    _ = dao.GroupSettlementDetail.Ctx(ctx).With(model.SettlementDetailEntity{}.Combo).
        Where(g.Map{
            c.GroupId:   group.Id,
            c.Ignorance: 0,
            c.State:     model.SettlementBilling,
        }).
        WhereLTE(c.StartDate, expDate).
        OrderAsc(c.StartDate).
        Scan(&rows)

    for _, row := range rows {
        days, needSplit, end := row.GetExpDays(now, expDate)
        num := decimal.NewFromFloat(row.UnitPrice).Mul(decimal.NewFromInt(int64(days))).Round(2)
        amount, _ := num.Float64()
        total, _ = decimal.NewFromFloat(total).Add(num).Float64()
        item := &model.SettlementListItem{
            DetailId:    row.Id,
            ComboName:   row.Combo.Name,
            BatteryType: row.BatteryType,
            BillDays:    days,
            UnitPrice:   row.UnitPrice,
            Amount:      amount,
            StopDate:    end.Format(layout),
            StartDate:   row.StartDate.Format(layout),
            NeedSplit:   needSplit,
        }
        items = append(items, item)
    }

    if len(rows) > 0 {
        from = rows[0].StartDate
    }

    bill = &model.SettlementCache{
        GroupId: group.Id,
        Amount:  total,
        Items:   items,
        ToDate:  expDate.Format(layout),
    }

    // 上次结算日期
    bill.FromDate = from.Format(layout)
    // var last *model.GroupSettlement
    // _ = dao.GroupSettlement.Ctx(ctx).
    //     Where(dao.GroupSettlement.Columns.GroupId, group.Id).
    //     OrderDesc(dao.GroupSettlement.Columns.Date).
    //     Scan(last)
    // if last == nil {
    //     bill.FromDate = group.CreatedAt.Format(layout)
    // } else {
    //     bill.FromDate = last.Date.AddDate(0, 0, 1).Format(layout)
    // }

    // 存入MongoDB数据库
    col := mongo.DB.Collection(group_settlement.Table)
    _ = col.CreateOneIndex(ctx, options.IndexModel{Key: []string{"groupId", "toDate"}})
    r, err := col.InsertOne(ctx, bill)
    if err == nil {
        bill.Hash = r.InsertedID.(primitive.ObjectID).Hex()
    }
    return
}

func (s *groupSettlementDetailService) ListDetails(ctx context.Context, groupId uint) (rows []*model.SettlementDetailEntity) {
    c := dao.GroupSettlementDetail.Columns
    _ = dao.GroupSettlementDetail.Ctx(ctx).
        Where(c.GroupId, groupId).
        Where(c.Ignorance, 0).
        OrderDesc("startDate").
        Scan(&rows)
    return rows
}

// GetDays 获取团队天数
// billDays 未结算天数
// days 总天数
func (s *groupSettlementDetailService) GetDays(ctx context.Context, groupId uint) (billDays uint, days uint) {
    rows := s.ListDetails(ctx, groupId)
    if len(rows) > 0 {
        for _, row := range rows {
            num := row.GetDays()
            days += num
            if row.State != model.SettlementSettled {
                billDays += num
            }
        }
    }
    return
}

// GetDaysGroupByUser 获取团队天数(用户分组)
func (s *groupSettlementDetailService) GetDaysGroupByUser(ctx context.Context, groupId uint) (data map[uint64]model.GroupUsageDays) {
    data = make(map[uint64]model.GroupUsageDays)
    rows := s.ListDetails(ctx, groupId)
    if len(rows) > 0 {
        for _, row := range rows {
            uid := row.UserId
            detail, ok := data[uid]
            if !ok {
                detail = model.GroupUsageDays{
                    Days:     0,
                    BillDays: 0,
                }
            }
            num := row.GetDays()
            detail.Days += num
            if row.State != model.SettlementSettled {
                detail.BillDays += row.GetDays()
            }
            data[uid] = detail
        }
    }
    return
}
