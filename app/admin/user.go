package admin

import (
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
)

var UserApi = userApi{}

type userApi struct {
}

func (*userApi) Login(r *ghttp.Request) {
	var req model.SysUserLoginReq
	if err := r.Parse(&req); err != nil {
		fmt.Println(err)
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	rep, err := service.SysUsersService.Login(r.Context(), req)
	if err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	response.JsonOkExit(r, rep)
}

func (*userApi) Logout(r *ghttp.Request) {
	err := service.SysUsersService.Logout(r.Context())
	if err != nil {
		response.JsonErrExit(r)
	}
	response.JsonOkExit(r)
}

func (*userApi) Profile(r *ghttp.Request) {
	//TODO
	rep := model.SysUserProfileRep{
		Name:         "管理员",
		Roles:        []string{"admin"},
		Avatar:       "https://wpimg.wallstcn.com/50530061-851b-4ca5-9dc5-2fead928a939.jpg?imageView2/2/h/150",
		Introduction: "超级管理员测试账号",
	}
	response.JsonOkExit(r, rep)
}
