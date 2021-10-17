package admin

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "net/url"
    "os"
    "path/filepath"
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

    // 判断合同文件是否存在
    if _, err := os.Stat(req.ContractFile); os.IsNotExist(err) {
        response.Json(r, response.RespCodeArgs, "合同文件不存在")
    }

    if err := dao.Group.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
        group := model.Group{
            Name:          req.Name,
            ProvinceId:    req.ProvinceId,
            CityId:        req.CityId,
            ContactName:   req.ContactName,
            ContactMobile: req.ContactMobile,
            ContractFile:  req.ContractFile,
        }
        groupId, err := service.GroupService.Create(ctx, group)
        if err != nil {
            return err
        }
        group.Id = groupId

        // 设置团签leader
        boss := model.User{
            GroupId:  groupId,
            RealName: req.ContactName,
            Mobile:   req.ContactMobile,
            Type:     model.UserTypeGroupBoss,
        }

        if err = service.UserService.AddOrSetGroupUser(ctx, boss); err != nil {
            return err
        }

        if err = service.GroupUserService.AddUsers(ctx, group, req.UserList); err != nil {
            return err
        }

        _ = service.GroupDailyStatService.GenerateWeek(ctx, group.Id, 60)
        _ = service.GroupDailyStatService.GenerateWeek(ctx, group.Id, 72)

        return err
    }); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r)
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

// Contract
// @Summary 获取合同, 若成功获取则直接返回二进制文件, 此时前端直接处理成文件下载; 若失败, 则返回json失败数据
// @Tags    管理
// @Accept  json
// @Param   id path int true "团签ID"
// @Produce  octet-stream
// @Produce  json
// @Router  /admin/group/{id}/contract [GET]
// @Success 200 {object} object "合同文件"
// @Failure 400,404 {object} response.JsonResponse "错误结果"
func (*groupApi) Contract(r *ghttp.Request) {
    id := r.GetInt("id")

    var group model.Group

    if err := dao.Group.Ctx(r.GetCtx()).Where("id = ?", id).Scan(&group); err != nil {
        response.Json(r, response.RespCodeArgs, "未找到团签")
    }

    if _, err := os.Stat(group.ContractFile); err != nil {
        response.Json(r, response.RespCodeArgs, "合同文件不存在")
    }

    r.Response.ServeFileDownload(group.ContractFile, url.QueryEscape(group.Name+"团签合同"+filepath.Ext(group.ContractFile)))
}

// AddMember
// @Summary 新增团签用户
// @Tags    管理
// @Accept  json
// @Param   id path int true "团签ID"
// @Param   entity body model.GroupCreateUserReq true "用户详情"
// @Produce  json
// @Router  /admin/group/{id}/member [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*groupApi) AddMember(r *ghttp.Request) {
    var req model.GroupCreateUserReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }

    groupId := r.GetInt("id")
    var group model.Group
    if err := dao.Group.Ctx(r.Context()).Where(g.Map{dao.Group.Columns.Id: groupId}).Scan(&group); err != nil {
        response.Json(r, response.RespCodeSystemError, err.Error())
    }

    if err := service.GroupUserService.AddUsers(r.Context(), group, []model.GroupCreateUserReq{req}); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r)
}

func (*groupApi) ListMember(r *ghttp.Request) {
}

func (*groupApi) DeleteMember(r *ghttp.Request) {
}
