// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-17
// Based on apiv2 by liasica, magicrolan@qq.com.

package mq

import (
    "fmt"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/util/gutil"
    "reflect"
    "strings"
)

// ParseStructToQuery 将interface转换为查询语句(=)
func ParseStructToQuery(req interface{}, without ...string) (params g.Map) {
    ex := make(map[string]struct{})
    for _, s := range without {
        ex[s] = struct{}{}
    }
    t := reflect.TypeOf(req)
    v := reflect.ValueOf(req)
    params = make(g.Map)
    for i := 0; i < t.NumField(); i++ {
        key := t.Field(i).Name
        val := v.Field(i)
        _, ok := ex[key]
        if !val.IsZero() && !ok {
            switch val.Kind() {
            case reflect.String, reflect.Uint, reflect.Int, reflect.Float64:
                params[strings.ToLower(key[:1])+key[1:]] = val.Interface()
            }
        }
    }
    return
}

func FieldsWithTable(table string, columns interface{}) (fields []string) {
    keys := gutil.Values(columns)
    for _, field := range keys {
        fields = append(fields, fmt.Sprintf("`%s`.`%v`", table, field))
    }
    return
}
