// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at https://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-07-12
// Based on apiv2 by liasica, magicrolan@qq.com.

package deploy

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
func StartListen() {
    if h, err := os.Hostname(); err != nil || h != "shiguangju" {
        return
    }

    hook, _ := github.New(github.Options.Secret("ILDUIctIvErOCaREuSHaNgYpudesToTi"))

    http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
        // 仅开启pull事件
        payload, err := hook.Parse(r, github.PullRequestEvent)
        if err != nil {
            if err == github.ErrEventNotFound {
                // ok event wasn;t one of the ones asked to be parsed
                return
            }
        }

        fmt.Printf("%v", payload)

        fmt.Println(exec.Command("bash", "-c", "git pull; go mod download all").Run())
    })
    _ = http.ListenAndServe(":3761", nil)
}
