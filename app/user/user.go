package user

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"

	"battery/app/model"
	"battery/app/service"
	"battery/library/esign/sign"
	"battery/library/qr"
	"battery/library/response"
)

var UserApi = userApi{}

type userApi struct {
}

// Register
// @summary 骑手-用户注册
// @tags    骑手
// @Accept  json
// @Produce  json
// @param   entity  body model.UserRegisterReq true "注册数据"
// @router  /rapi/register [POST]
// @success 200 {object} response.JsonResponse  "返回结果"
func (*userApi) Register(r *ghttp.Request) {
	var req model.UserRegisterReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	if _, err := service.UserService.Register(r.Context(), req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	response.JsonOkExit(r)
}

// Login
// @summary 骑手-用户登录
// @tags    骑手
// @Accept  json
// @Produce  json
// @param   entity  body model.UserLoginReq true "登录数据"
// @router  /rapi/login [POST]
// @success 200 {object} response.JsonResponse{data=model.UserLoginRep}  "返回结果"
func (*userApi) Login(r *ghttp.Request) {
	var req model.UserLoginReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	if data, err := service.UserService.Login(r.Context(), req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	} else {
		response.JsonOkExit(r, data)
	}
}

// Auth
// @summary 骑手-实名认证提交
// @tags    骑手
// @Accept  json
// @Produce  json
// @param   entity  body model.UserRealNameAuthReq true "认证数据"
// @router  /rapi/auth [POST]
// @success 200 {object} response.JsonResponse{data=model.UserRealNameAuthRep}  "返回结果"
func (*userApi) Auth(r *ghttp.Request) {
	var req model.UserRealNameAuthReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	if user, err := service.UserService.GetUserByIdCardNo(r.Context(), req.IdCardNo); err == nil && user.AuthState == model.AuthStateVerifySuccess {
		response.Json(r, response.RespCodeArgs, fmt.Sprintf("证件号码 %s 已认证超过，请检查证件号码", req.IdCardNo))
	}
	if res, err := service.UserService.RealNameAuthSubmit(r.Context(), req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	} else {
		response.JsonOkExit(r, res)
	}
}

// AuthGet
// @summary 骑手-获取实名认证状态 [每次切换页面都进行查询]
// @tags    骑手
// @Accept  json
// @Produce  json
// @router  /rapi/auth [GET]
// @success 200 {object} response.JsonResponse{data=int}  "返回结果"
func (*userApi) AuthGet(r *ghttp.Request) {
	u := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	response.JsonOkExit(r, u.AuthState)
}

// PushToken
// @summary 骑手-上报推送token
// @tags    骑手-消息
// @Accept  json
// @Produce  json
// @param   entity  body model.PushTokenReq true "登录数据"
// @router  /rapi/device  [PUT]
// @success 200 {object} response.JsonResponse  "返回结果"
func (*userApi) PushToken(r *ghttp.Request) {
	var req model.PushTokenReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	err := service.UserService.PushToken(r.Context(), req)
	if err != nil {
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	response.JsonOkExit(r)
}

// Packages
// @summary 骑手-获取骑手当前套餐详情
// @tags    骑手
// @Accept  json
// @Produce  json
// @router  /rapi/package  [GET]
// @success 200 {object} response.JsonResponse{data=model.UserCurrentPackageOrder}  "返回结果"
func (*userApi) Packages(r *ghttp.Request) {
	rep, err := service.UserService.MyPackage(r.Context())
	if err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	response.JsonOkExit(r, rep)
}

// PackagesOrderQr
// @summary 骑手-获取骑手当前套餐二维码
// @tags    骑手
// @Accept  json
// @Produce  json
// @router  /rapi/package_order/qr  [GET]
// @success 200 {object} response.JsonResponse "返回结果, data字段为二维码图片数据，需要本地生成二维码"
func (*userApi) PackagesOrderQr(r *ghttp.Request) {
	u := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	if u.GroupId > 0 {
		response.JsonOkExit(r, fmt.Sprintf("%d-%s-%d", u.GroupId, u.Qr, u.BatteryType))
	} else {
		if u.PackagesOrderId == 0 {
			response.Json(r, response.RespCodeArgs, "还未购买套餐")
		}
		order, err := service.PackagesOrderService.Detail(r.Context(), u.PackagesOrderId)
		if err != nil {
			response.Json(r, response.RespCodeArgs, "为找到订单")
		}
		response.JsonOkExit(r, order.No)
	}
}

// Profile
// @summary 骑手-首页
// @tags    骑手
// @Accept  json
// @Produce  json
// @router  /rapi/home  [GET]
// @success 200 {object} response.JsonResponse{data=model.UserProfileRep}  "返回结果"
func (*userApi) Profile(r *ghttp.Request) {
	profile := service.UserService.Profile(r.Context())
	profile.Qr = qr.Code.AddPrefix(profile.Qr)
	response.JsonOkExit(r, profile)
}

type UserSignFileRepItem struct {
	FileName string `json:"fileName"` //文件名称
	FileUrl  string `json:"fileUrl"`  //文件地址
}

type UserSignFileRep []*UserSignFileRepItem

// SignFile
// @summary 骑手-签约文件地址
// @tags    骑手
// @Accept  json
// @Produce  json
// @router  /rapi/sign_file  [GET]
// @success 200 {object} response.JsonResponse{data=[]user.UserSignFileRep}  "返回结果"
func (*userApi) SignFile(r *ghttp.Request) {
	u := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	s, err := service.SignService.UserLatestDoneDetail(r.Context(), u.Id, u.PackagesOrderId, u.GroupId)
	if err != nil || s == nil {
		response.JsonErrExit(r, response.RespCodeNotFound)
	}
	res, err := sign.Service().SignFlowDocuments(s.FlowId)
	if err != nil || res.Code != 0 {
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	files := make([]*UserSignFileRepItem, len(res.Data.Docs))
	for i, f := range res.Data.Docs {
		files[i] = &UserSignFileRepItem{
			FileName: f.FileName,
			FileUrl:  f.FileUrl,
		}
	}
	response.JsonOkExit(r, files)
}
