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
)

var GroupApi = groupApi{}

type groupApi struct {
}

// Create
// @Summary 创建团签
// @Tags    管理
// @Accept  json
// @Param   entity body model.GroupFormReq true "团签详情"
// @Produce  json
// @Router  /admin/group [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*groupApi) Create(r *ghttp.Request) {
    var req model.GroupFormReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if !service.GroupService.CheckName(r.Context(), 0, req.Name) {
        response.Json(r, response.RespCodeArgs, "名称已被使用")
    }
    // 重复手机号验证
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
    // 用户状态验证
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
                return _err
            }
        }
        if len(userInsertData) > 0 {
            if _err := service.UserService.CreateGroupUsers(ctx, userInsertData); _err != nil {
                return _err
            }
        }
        if len(usersIds) > 0 {
            if _err := service.UserService.SetUsersGroupId(ctx, usersIds, groupId); _err != nil {
                return _err
            }
        }
        groupUsers := service.UserService.GetByMobiles(ctx, userMobiles)
        groupUserIds := make([]uint64, len(groupUsers))
        for i, user := range groupUsers {
            groupUserIds[i] = user.Id
        }
        if _err := service.GroupUserService.BatchCreate(ctx, groupUserIds, groupId); _err != nil {
            return _err
        }
        if _err := service.GroupDailyStatService.GenerateWeek(ctx, groupId, 60); _err != nil {
            return _err
        }
        if _err := service.GroupDailyStatService.GenerateWeek(ctx, groupId, 72); _err != nil {
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

// List
// @Summary 团签列表
// @Tags    管理
// @Accept  json
// @Param   entity body model.GroupListReq true "请求参数"
// @Produce  json
// @Router  /admin/group [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.GroupListItem}}  "返回结果"
func (*groupApi) List(r *ghttp.Request) {
    var req model.GroupListReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    total, items := service.GroupService.ListAdmin(r.Context(), model.GroupListAdminReq{
        Page:     req.Page,
        Keywords: req.Keywords,
    })
    var result model.ItemsWithTotal
    if total > 0 {
        groupIds := make([]uint, len(items))
        for key, group := range items {
            groupIds[key] = group.Id
        }
        groupUserCnt := service.UserService.GroupUserCnt(r.Context(), groupIds)
        stat, err := service.GroupDailyStatService.StatDateRange(r.Context(), groupIds, req.StartDate, req.EndDate)
        if err != nil {
            response.JsonErrExit(r)
        }
        for _, group := range items {
            result.Items = append(result.Items, model.GroupListItem{
                Id:               group.Id,
                Name:             group.Name,
                CityId:           group.CityId,
                ProvinceId:       group.ProvinceId,
                ContactName:      group.ContactName,
                ContactMobile:    group.ContactMobile,
                UserCnt:          groupUserCnt[group.Id],
                ArrearsDaysCnt60: stat[group.Id].ArrearsDaysCnt60,
                DaysCnt60:        stat[group.Id].ArrearsDaysCnt60 + stat[group.Id].PaidDaysCnt60,
                ArrearsDaysCnt72: stat[group.Id].ArrearsDaysCnt72,
                DaysCnt72:        stat[group.Id].ArrearsDaysCnt60 + stat[group.Id].PaidDaysCnt72,
            })
        }
    }
    response.JsonOkExit(r, result)
}
