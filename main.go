package main

import (
    _ "battery/boot"
    _ "battery/router"

    "github.com/gogf/gf/frame/g"
)

// @title 时光驹API
// @version 1.0

// @host apiv2.shiguangjv.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-ACCESS-TOKEN

func main() {
    g.Server().Run()
}
