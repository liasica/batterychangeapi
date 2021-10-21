package user

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var GroupApi = groupApi{}

type groupApi struct {
}

// Stat 团队统计
// @Summary 骑手-团签团队统计
// @Tags    骑手-团签BOSS
// @Accept  json
// @Produce  json
// @Router  /rapi/group/stat  [GET]
// @Success 200 {object} response.JsonResponse{data=model.UserGroupStatRep}  "返回结果"
func (*groupApi) Stat(r *ghttp.Request) {
    user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
    billDays, _ := service.GroupSettlementDetailService.GetDays(r.Context(), user.GroupId)
    response.JsonOkExit(r, model.UserGroupStatRep{
        MemberCnt: service.GroupUserService.UserCnt(r.Context(), user.GroupId),
        BillDays:  billDays,
    })
}

// List 团队详情
// @Summary 骑手-团签团队详情
// @Tags    骑手-团签BOSS
// @Accept  json
// @Produce json
// @Router  /rapi/group/list  [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.UserGroupListRep}  "返回结果"
func (*groupApi) List(r *ghttp.Request) {
    user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
    userIds := service.GroupUserService.UserIds(r.Context(), user.GroupId)
    if len(userIds) > 0 {
        users := service.UserService.GetByIds(r.Context(), userIds)
        var list []model.UserGroupListRep
        m := service.GroupSettlementDetailService.GetDaysGroupByUser(r.Context(), user.GroupId)
        for key, u := range users {
            uid := u.Id
            data, ok := m[uid]
            var billDays, days uint
            if ok {
                billDays = data.BillDays
                days = data.Days
            }
            list[key] = model.UserGroupListRep{
                RealName: u.RealName,
                Days:     days,
                BillDays: billDays,
            }
        }
        response.JsonOkExit(r, list)
    }
    response.JsonOkExit(r, make([]model.UserGroupListRep, 0))
}
