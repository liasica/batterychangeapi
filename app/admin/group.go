package admin

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "github.com/gogf/gf/os/gtime"
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
// @Produce json
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

        // _ = service.GroupDailyStatService.GenerateWeek(ctx, group.Id, model.BatteryType60)
        // _ = service.GroupDailyStatService.GenerateWeek(ctx, group.Id, model.BatteryType72)

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
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Produce json
// @Router  /admin/group [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.GroupEntity}} "返回结果"
func (*groupApi) List(r *ghttp.Request) {
    req := new(model.GroupListAdminReq)
    _ = request.ParseRequest(r, req)
    total, items := service.GroupService.ListAdmin(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// Contract
// @Summary 获取合同, 若成功获取则直接返回二进制文件, 此时前端直接处理成文件下载; 若失败, 则返回json失败数据
// @Tags    管理
// @Accept  json
// @Param   id path int true "团签ID"
// @Produce octet-stream
// @Produce json
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
// @Produce json
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

// ListMember
// @Summary 团签成员列表
// @Tags    管理
// @Accept  json
// @Param   id path int true "团签ID"
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Param 	groupId query integer false "团签ID"
// @Param 	realName query string false "成员姓名"
// @Param 	mobile query string false "成员电话"
// @Param 	batteryState query integer false "换电状态, 个签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租 5已逾期; 团签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租" ENUMS(0,1,2,3,4,5)
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Produce json
// @Router  /admin/group/{id}/member [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.UserListItem}} "返回结果"
func (*groupApi) ListMember(r *ghttp.Request) {
    req := new(model.UserListReq)
    _ = request.ParseRequest(r, req)
    req.GroupId = r.GetUint("id")

    total, items := service.UserService.ListPersonalItems(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// ListMemberBiz
// @Summary 团签成员换电记录
// @Tags    管理
// @Accept  json
// @Param   id path int true "团签ID"
// @Param   userId path int true "成员ID"
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Produce json
// @Router  /admin/group/{id}/member/{userId}/biz [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.BizSimpleItem}} "返回结果"
func (*groupApi) ListMemberBiz(r *ghttp.Request) {
    groupId := r.GetUint("id")
    userId := r.GetUint("userId")
    var page = new(model.Page)
    _ = request.ParseRequest(r, page)

    // 查找用户
    c := dao.User.Columns
    var member = new(model.User)
    _ = dao.User.Ctx(r.Context()).Where(c.Id, userId).Where(c.GroupId, groupId).Scan(member)
    if member == nil {
        response.Json(r, response.RespCodeArgs, "未找到该用户")
    }
    total, items := service.UserBizService.ListSimaple(r.Context(), member, page)
    response.ItemsWithTotal(r, total, items)
}

// DeleteMember
// @Summary 删除团签成员
// @Tags    管理
// @Accept  json
// @Param   id path int true "团签ID"
// @Param   userId path int true "成员ID"
// @Produce json
// @Router  /admin/group/{id}/member/{userId} [DELETE]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*groupApi) DeleteMember(r *ghttp.Request) {
    memberId := r.GetUint("userId")
    groupId := r.GetUint("id")

    if err := service.GroupUserService.Delete(r.Context(), groupId, memberId); err == nil {
        response.JsonOkExit(r)
    } else {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
}

// GetSettlement
// @Summary 获取团签账单
// @Tags    管理
// @Accept  json
// @Param   id path int true "团签ID"
// @Param   expDate path string true "截止日期"
// @Produce json
// @Router  /admin/group/{id}/bill/{expDate} [GET]
// @Success 200 {object} response.JsonResponse{data=model.SettlementCache} "返回结果"
func (*groupApi) GetSettlement(r *ghttp.Request) {
    id := r.GetInt("id")
    d := r.GetString("expDate")
    expDate := gtime.NewFromStr(d)
    if expDate.IsZero() || !gtime.Now().After(expDate.AddDate(0, 0, 1)) {
        response.Json(r, response.RespCodeArgs, "请携带正确的日期")
    }

    group := new(model.Group)
    if err := dao.Group.Ctx(r.Context()).Where(dao.Group.Columns.Id, id).Scan(group); err != nil {
        response.Json(r, response.RespCodeArgs, "未找到团签")
    }

    bill, err := service.GroupSettlementDetailService.GetGroupBill(r.Context(), group, expDate)
    if err != nil {
        g.Log("settlement").Errorf("结算单生成失败: %v", err)
        response.JsonErrExit(r, response.RespCodeSystemError, err)
    }
    response.JsonOkExit(r, bill)
}

// PostSettlement
// @Summary 结账(hash从[获取团签账单/admin/group/{id}/bill/{expDate}]拿取)
// @Tags    管理
// @Accept  json
// @Param   entity body model.GroupSettlementCheckoutReq true "结算请求"
// @Produce json
// @Router  /admin/group/bill [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*groupApi) PostSettlement(r *ghttp.Request) {
    var req = new(model.GroupSettlementCheckoutReq)
    _ = request.ParseRequest(r, req)
    // 查找结算单
    err := service.GroupSettlementService.CheckoutBill(r.Context(), req)
    if err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r)
}
