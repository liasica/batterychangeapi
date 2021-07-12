package admin

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
	"context"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
)

var GroupApi = groupApi{}

type groupApi struct {
}

type groupFormReq struct {
	Id            uint                 `v:"integer|min:0" json:"id"`
	Name          string               `v:"required" json:"name"`
	ContactName   string               `v:"required" json:"contactName"`
	ContactMobile string               `v:"required" json:"contactMobile"`
	ProvinceId    uint                 `v:"required|integer|min:1" json:"provinceId"`
	CityId        uint                 `v:"required|integer|min:1" json:"cityId"`
	UserList      []groupCreateReqUser `v:"required" json:"userList"`
}

type groupCreateReqUser struct {
	Id     int    `v:"required|integer|min:0" json:"id"`
	Name   string `v:"required|length=6,30" json:"name"`
	Mobile string `v:"required|phone-loose" json:"mobile"`
}

func (*groupApi) Create(r *ghttp.Request) {
	var req groupFormReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	if !service.GroupService.CheckName(r.Context(), req.Id, req.Name) {
		response.Json(r, response.RespCodeArgs, "公司名称已被使用")
	}
	//重复手机号验证
	userMobilesMap := make(map[string]bool, len(req.UserList)+1)
	userMobilesMap[req.ContactMobile] = true
	userMobiles := make([]string, len(req.UserList)+1)
	userMobiles = append(userMobiles, req.ContactMobile)
	for _, user := range req.UserList {
		if _, ok := userMobilesMap[user.Mobile]; ok {
			response.Json(r, response.RespCodeArgs, fmt.Sprintf("手机号 %s 重复", user.Mobile))
		}
		userMobilesMap[user.Mobile] = true
		userMobiles = append(userMobiles, user.Mobile)
	}
	//用户状态验证
	users := service.UserService.GetByMobiles(r.Context(), userMobiles)
	usersIds := make([]uint64, 0)
	for _, user := range users {
		if user.GroupId > 0 {
			response.Json(r, response.RespCodeArgs, fmt.Sprintf("手机号 %s 已经是其它团签成员，无法添加", user.Mobile))
		}
		if user.BatteryState != model.BatteryStateDefault && user.BatteryState != model.BatteryStateExit {
			response.Json(r, response.RespCodeArgs, fmt.Sprintf("手机号 %s 正在使用的中电池，无法添加", user.Mobile))
		}
		userMobilesMap[user.Mobile] = false
		usersIds = append(usersIds, user.Id)
	}
	if dao.Group.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
		groupId, err := service.GroupService.Create(ctx, model.Group{
			Name:          req.Name,
			ProvinceId:    req.ProvinceId,
			CityId:        req.CityId,
			ContactName:   req.ContactName,
			ContactMobile: req.ContactMobile,
		})

		if err != nil {
			return err
		}
		userInsertData := make([]model.User, 0)
		for _, user := range req.UserList {
			if userMobilesMap[user.Mobile] {
				userInsertData = append(userInsertData, model.User{
					GroupId:  groupId,
					RealName: user.Name,
					Mobile:   user.Mobile,
					Type:     model.UserTypeGroupRider,
				})
			}
		}
		if newBossUser := userMobilesMap[req.ContactMobile]; newBossUser {
			userInsertData = append(userInsertData, model.User{
				GroupId:  groupId,
				RealName: req.ContactName,
				Mobile:   req.ContactMobile,
				Type:     model.UserTypeGroupBoss,
			})
		} else {
			if _err := service.UserService.SetUserTypeGroupBoss(ctx, req.ContactName, groupId); _err != nil {
				fmt.Println(err)
				return _err
			}
		}
		if len(userInsertData) > 0 {
			if _err := service.UserService.CreateGroupUsers(ctx, userInsertData); _err != nil {
				fmt.Println(err)
				return _err
			}
		}
		if len(usersIds) > 0 {
			if _err := service.UserService.SetUsersGroupId(ctx, usersIds, groupId); _err != nil {
				fmt.Println(err)
				return _err
			}
		}
		groupUsers := service.UserService.GetByMobiles(ctx, userMobiles)
		groupUserIds := make([]uint64, len(groupUsers))
		for i, user := range groupUsers {
			groupUserIds[i] = user.Id
		}
		if _err := service.GroupUserService.BatchCreate(ctx, groupUserIds, groupId); _err != nil {
			fmt.Println(err)
			return _err
		}
		return nil
	}) != nil {
		response.JsonErrExit(r)
	}
	response.JsonOkExit(r)
}

func (*groupApi) Edit(r *ghttp.Request) {

}

func (*groupApi) Detail(r *ghttp.Request) {

}

type groupListItem struct {
	Id               uint   `json:"id"`
	Name             string `json:"name"`
	CityName         string `json:"cityName"`
	UserCnt          uint   `json:"userCnt"`
	DaysCnt60        uint   `json:"daysCnt60"`
	DaysCnt72        uint   `json:"daysCnt72"`
	ArrearsDaysCnt60 uint   `json:"arrearsDaysCnt60"`
	ArrearsDaysCnt72 uint   `json:"arrearsDaysCnt72"`
	ContactName      string `json:"contactName"`
	ContactMobile    string `json:"contactMobile"`
}

type groupListReq struct {
	model.GroupListAdminReq
	StartDate *gtime.Time `json:"startDate"`
	EndDate   *gtime.Time `json:"endDate"`
}

func (*groupApi) List(r *ghttp.Request) {
	var req groupListReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	total, items := service.GroupService.ListAdmin(r.Context(), model.GroupListAdminReq{
		Page:     req.Page,
		Keywords: req.Keywords,
	})
	rep := struct {
		Total int             `json:"total"`
		Items []groupListItem `json:"items"`
	}{Total: total}
	if total > 0 {
		groupIds := make([]uint, len(items))
		cityIds := make([]uint, len(items))
		rep.Items = make([]groupListItem, len(items))
		for key, group := range items {
			cityIds[key] = group.CityId
			groupIds[key] = group.Id
		}
		cityIdName := service.DistrictsService.MapIdName(r.Context(), cityIds)
		groupUserCnt := service.UserService.GroupUserCnt(r.Context(), groupIds)
		stat, err := service.GroupDailyStatService.StatDateRange(r.Context(), groupIds, req.StartDate, req.EndDate)
		if err != nil {
			response.JsonErrExit(r)
		}
		for key, group := range items {
			rep.Items[key] = groupListItem{
				Id:               group.Id,
				Name:             group.Name,
				CityName:         cityIdName[group.CityId],
				ContactName:      group.ContactName,
				ContactMobile:    group.ContactMobile,
				UserCnt:          groupUserCnt[group.Id],
				ArrearsDaysCnt60: stat[group.Id].ArrearsDaysCnt60,
				DaysCnt60:        stat[group.Id].ArrearsDaysCnt60 + stat[group.Id].PaidDaysCnt60,
				ArrearsDaysCnt72: stat[group.Id].ArrearsDaysCnt72,
				DaysCnt72:        stat[group.Id].ArrearsDaysCnt60 + stat[group.Id].PaidDaysCnt72,
			}
		}
	}
	response.JsonOkExit(r, rep)
}
