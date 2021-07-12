package model

import (
	"battery/app/model/internal"
)

type SysUsers internal.SysUsers

type SysUserLoginReq struct {
	Username string `json:"username" v:"required|length:5,16"`
	Password string `json:"password" v:"required|length:6,16"`
}

type SysUserLoginRep struct {
	AccessToken string `json:"accessToken"`
}

type SysUserProfileRep struct {
	Roles        []string `json:"roles"`
	Name         string   `json:"name"`
	Avatar       string   `json:"avatar"`
	Introduction string   `json:"introduction"`
}
