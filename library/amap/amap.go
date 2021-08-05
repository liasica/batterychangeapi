// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-08-05
// Based on apiv2 by liasica, magicrolan@qq.com.

package amap

import (
    "errors"
    "fmt"
    "github.com/gogf/gf/frame/g"
)

// GetRegeo 逆地理编码
func GetRegeo(lng, lat string) (geo *Regeo, err error) {
    geo = new(Regeo)
    location := fmt.Sprintf("%s,%s", lng, lat)
    err = g.Client().GetVar(fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%s&radius=1000&extensions=base", g.Cfg().GetString("amap.key"), location)).Scan(geo)
    if err != nil {
        return
    }
    if geo.Status != RESOK {
        return nil, errors.New(geo.Info)
    }
    return
}
