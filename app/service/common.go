// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-08-05
// Based on apiv2 by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/gogf/gf/frame/g"
)

type common struct{}

type Weather struct {
    Code       string     `json:"code"`
    UpdateTime string     `json:"updateTime"`
    FxLink     string     `json:"fxLink"`
    Now        WeatherNow `json:"now"`
}

type WeatherNow struct {
    ObsTime   string `json:"obsTime"`   // 数据观测时间
    Temp      string `json:"temp"`      // 温度，默认单位：摄氏度
    FeelsLike string `json:"feelsLike"` // 体感温度，默认单位：摄氏度
    Icon      string `json:"icon"`      // 天气状况和图标的代码
    Text      string `json:"text"`      // 天气状况的文字描述，包括阴晴雨雪等天气状态的描述
    Wind360   string `json:"wind360"`   // 风向360角度
    WindDir   string `json:"windDir"`   // 风向
    WindScale string `json:"windScale"` // 风力等级
    WindSpeed string `json:"windSpeed"` // 风速，公里/小时
    Humidity  string `json:"humidity"`  // 相对湿度，百分比数值
    Precip    string `json:"precip"`    // 当前小时累计降水量，默认单位：毫米
    Pressure  string `json:"pressure"`  // 大气压强，默认单位：百帕
    Vis       string `json:"vis"`       // 能见度，默认单位：公里
    Cloud     string `json:"cloud"`     // 云量，百分比数值
    Dew       string `json:"dew"`       // 露点温度
}

var CommonService = new(common)

// WeatherNow 实时天气
func (*common) WeatherNow(lng, lat string) (now WeatherNow, err error) {
    // geo := new(amap.Regeo)
    // geo, err = amap.GetRegeo(lng, lat)
    // code := geo.Regeocode.AddressComponent.Adcode
    cfg := g.Cfg()
    // 获取天气（和风）
    w := new(Weather)
    err = g.Client().GetVar(fmt.Sprintf("%s/v7/weather/now?key=%s&location=%s,%s", cfg.GetString("qweather.url"), cfg.GetString("qweather.key"), lng, lat)).Scan(w)
    if err != nil {
        return
    }
    return w.Now, nil
}
