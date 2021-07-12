package esign

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

var Middleware = middleware{}

type middleware struct {
}

// Ip 鉴权中间件，只有指定IP才能访问
func (m *middleware) Ip(r *ghttp.Request) {
	if r.GetClientIp() != g.Cfg().GetString("eSign.ip") {
		r.Response.Status = http.StatusForbidden
		r.Exit()
	}
}
