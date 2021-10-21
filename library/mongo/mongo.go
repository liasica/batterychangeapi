// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-20
// Based on apiv2 by liasica, magicrolan@qq.com.

package mongo

import (
    "context"
    "github.com/gogf/gf/frame/g"
    "github.com/qiniu/qmgo"
)

type Mongo struct {
    *qmgo.Database
}

var DB *Mongo

func Connect() {
    link := g.Cfg().GetString("mongo.link")
    database := g.Cfg().GetString("mongo.database")
    if link == "" || database == "" {
        g.Log().Fatalf("MongoDB配置读取失败")
    }
    ctx := context.Background()
    client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: link})
    if err != nil {
        g.Log().Fatalf("MongoDB启动失败")
    }
    db := client.Database(database)
    DB = &Mongo{db}
}
