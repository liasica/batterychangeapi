// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-20
// Based on apiv2 by liasica, magicrolan@qq.com.

package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/group_settlement"
    "battery/app/model/group_settlement_detail"
    "battery/library/mongo"
    "context"
    "errors"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/os/gtime"
    jsoniter "github.com/json-iterator/go"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var GroupSettlementService = new(groupSettlementService)

type groupSettlementService struct {
}

// CheckoutBill 发起结算
func (s *groupSettlementService) CheckoutBill(ctx context.Context, req *model.GroupSettlementCheckoutReq) (err error) {
    now := gtime.Now()
    // 查找缓存
    col := mongo.DB.Collection(group_settlement.Table)
    objectID, err := primitive.ObjectIDFromHex(req.Hash)
    if err != nil {
        return errors.New("结算单hash错误，请携带正确参数")
    }
    var bill *model.SettlementCache
    if err = col.Find(ctx, bson.M{"_id": objectID}).One(&bill); err != nil || bill == nil {
        return errors.New("未找到结算单，请重新生成")
    }
    // 获取所有结算单详情
    var ids []uint64
    var details []model.GroupSettlementDetail
    mapItems := map[uint64]*model.SettlementListItem{}
    for _, item := range bill.Items {
        ids = append(ids, item.DetailId)
        mapItems[item.DetailId] = item
    }
    if err = dao.GroupSettlementDetail.Ctx(ctx).
        WhereIn("id", ids).
        Scan(&details); err != nil {
        return err
    }
    if len(bill.Items) != len(details) {
        return errors.New("请重新生成结算单")
    }
    sysUser := ctx.Value(model.ContextAdminKey).(*model.ContextAdmin)
    // 结算账单
    return dao.GroupSettlement.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
        // 存储结算账单
        ds, _ := jsoniter.MarshalToString(bill.Items)
        gs := model.GroupSettlement{
            Hash:      bill.Id.Hex(),
            GroupId:   bill.GroupId,
            Amount:    bill.Amount,
            SysUserId: sysUser.Id,
            SysName:   sysUser.Username,
            Date:      gtime.NewFromStr(bill.ToDate),
            Remark:    req.Remark,
            Detail:    ds,
        }
        r, err := tx.Save(group_settlement.Table, gs)
        if err != nil {
            return err
        }
        id, _ := r.LastInsertId()

        for k, detail := range details {
            item := mapItems[detail.Id]
            end := gtime.NewFromStr(item.StopDate)
            // 分割账单
            if item.NeedSplit {
                newRecord := detail
                newRecord.Id = 0
                newRecord.StartDate = end.AddDate(0, 0, 1)
                newRecord.ParentId = detail.Id
                newRecord.CreatedAt = now
                details = append(details, newRecord)
            }
            details[k].SettlementId = uint64(id)
            details[k].SplitAt = now
            details[k].State = model.SettlementSettled
            details[k].StopDate = end
        }
        _, err = tx.Save(group_settlement_detail.Table, details)
        return err
    })
}
