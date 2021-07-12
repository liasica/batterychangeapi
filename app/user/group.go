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

//Stat 团队统计
// @summary 骑手-团签团队统计
// @tags    骑手-团签BOSS
// @Accept  json
// @Produce  json
// @router  /rapi/group/stat  [GET]
// @success 200 {object} response.JsonResponse{data=model.UserGroupStatRep}  "返回结果"
func (*groupApi) Stat(r *ghttp.Request) {
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	response.JsonOkExit(r, model.UserGroupStatRep{
		UserCnt: service.GroupUserService.UserCnt(r.Context(), user.GroupId),
		Days:    service.GroupService.StatDays(r.Context(), user.GroupId),
	})
}

//List 团队详情
// @summary 骑手-团签团队详情
// @tags    骑手-团签BOSS
// @Accept  json
// @Produce  json
// @router  /rapi/group/list  [GET]
// @success 200 {object} response.JsonResponse{data=[]model.UserGroupListRep}  "返回结果"
func (*groupApi) List(r *ghttp.Request) {
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	userIds := service.GroupUserService.UserIds(r.Context(), user.GroupId)
	if len(userIds) > 0 {
		users := service.UserService.GetByIds(r.Context(), userIds)
		list := make([]model.UserGroupListRep, len(users))
		for key, u := range users {
			list[key] = model.UserGroupListRep{
				Name: u.RealName,
				//TODO 使用天数
			}
		}
		response.JsonOkExit(r, list)
	}
	response.JsonOkExit(r, make([]model.UserGroupListRep, 0))
}
