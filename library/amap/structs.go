// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-08-05
// Based on apiv2 by liasica, magicrolan@qq.com.

package amap

const (
    RESOK = "1"
)

type Regeo struct {
    Status    string `json:"status"`
    Info      string `json:"info"`
    Infocode  string `json:"infocode"`
    Regeocode struct {
        FormattedAddress string `json:"formatted_address"`
        AddressComponent struct {
            Province string `json:"province"`
            Adcode   string `json:"adcode"`
            District string `json:"district"`
            Country  string `json:"country"`
            Township string `json:"township"`
            Citycode string `json:"citycode"`
        } `json:"addressComponent"`
    } `json:"regeocode"`
}
