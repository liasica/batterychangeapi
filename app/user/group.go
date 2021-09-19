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
// @summary 骑手-团签团队统计
// @tags    骑手-团签BOSS
// @Accept  json
// @Produce  json
// @router  /rapi/group/stat  [GET]
// @success 200 {object} response.JsonResponse{data=model.UserGroupStatRep}  "返回结果"
func (*groupApi) Stat(r *ghttp.Request) {
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	var days uint = 0
	res, _ := service.GroupDailyStatService.ArrearsDays(r.Context(), user.GroupId)
	for _, day := range res {
		days = days + day.Cnt
	}
	response.JsonOkExit(r, model.UserGroupStatRep{
		UserCnt: service.GroupUserService.UserCnt(r.Context(), user.GroupId),
		Days:    days,
	})
}

// List 团队详情
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
		userArrearsDays := make(map[uint64]uint, len(userIds))
		for _, userId := range userIds {
			userArrearsDays[userId] = 0
		}
		arrearsList, _ := service.GroupDailyStatService.ArrearsList(r.Context(), user.GroupId)
		for _, row := range arrearsList {
			for _, userId := range row.UserIds {
				if _, ok := userArrearsDays[userId]; ok {
					userArrearsDays[userId]++
				}
			}
		}
		for key, u := range users {
			list[key] = model.UserGroupListRep{
				Name: u.RealName,
				Days: userArrearsDays[u.Id],
			}
		}
		response.JsonOkExit(r, list)
	}
	response.JsonOkExit(r, make([]model.UserGroupListRep, 0))
}
