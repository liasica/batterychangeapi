// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at https://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-07-12
// Based on apiv2 by liasica, magicrolan@qq.com.

package main

import (
    "fmt"
    "github.com/go-playground/webhooks/v6/github"
    "net/http"
    "os"
    "os/exec"
)

const (
    path = "/webhooks"
)

// StartListen GitHub webhook，仅服务器生效（通过hostname判断）
func main() {
    if h, err := os.Hostname(); err != nil || h != "shiguangju" {
        return
    }

    hook, _ := github.New(github.Options.Secret("ILDUIctIvErOCaREuSHaNgYpudesToTi"))

    http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
        // 仅开启push事件
        _, err := hook.Parse(r, github.PushEvent)
        if err != nil {
            if err == github.ErrEventNotFound {
                // ok event wasn;t one of the ones asked to be parsed
                return
            }
        }

        output, err := exec.Command("bash", "-c", "cd /var/www/apiv2; git pull; go mod download all; pm2 restart apiv2").Output()
        if err != nil {
            fmt.Printf("执行失败: %v", err)
            return
        }
        fmt.Println(string(output))
    })
    _ = http.ListenAndServe(":3761", nil)
}
