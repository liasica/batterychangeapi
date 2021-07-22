package user

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/app/service"
	"battery/library/esign/sign"
	beansSign "battery/library/esign/sign/beans"
	"battery/library/payment/alipay"
	"battery/library/payment/wechat"
	"battery/library/response"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/golang-module/carbon"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
)

var BizApi = bizApi{}

type bizApi struct {
}

//RecordStat 骑手换电记录统计
// @summary 骑手-换电记录统计
// @Accept  json
// @Produce  json
// @tags    骑手
// @router  /rapi/biz_record/stat [GET]
// @success 200 {object} response.JsonResponse{data=model.UserBizRecordStatRep}  "返回结果"
func (*bizApi) RecordStat(r *ghttp.Request) {
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	days := user.BizBatteryRenewalDays
	if !user.BizBatteryRenewalDaysStartAt.IsZero() {
		days = days + uint(carbon.Parse(user.BizBatteryRenewalDaysStartAt.String()).DiffInDays(carbon.Parse(gtime.Now().String())))
	}
	response.JsonOkExit(r, model.UserBizRecordStatRep{
		Count: user.BizBatteryRenewalCnt,
		Days:  days,
	})
}

//RecordList 骑手换电记录列表
// @summary 骑手-换电记录列表
// @Accept  json
// @Produce  json
// @tags    骑手
// @param 	pageIndex query integer  true "当前页码"
// @param 	pageLimit query integer  true "每页行数"
// @router  /rapi/biz_record/list [GET]
// @success 200 {object} response.JsonResponse{data=[]model.UserBizRecordListRep}  "返回结果"
func (*bizApi) RecordList(r *ghttp.Request) {
	var req model.Page
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	records := service.UserBizService.ListUser(r.Context(), req)
	if len(records) > 0 {
		list := make([]model.UserBizRecordListRep, len(records))
		var shopIds []uint
		var cityIds []uint
		for _, record := range records {
			shopIds = append(shopIds, record.ShopId)
			cityIds = append(cityIds, record.CityId)
		}
		shopMapIdName := service.ShopService.MapIdName(r.Context(), shopIds)
		cityMapIdName := service.DistrictsService.MapIdName(r.Context(), cityIds)
		for key, record := range records {
			list[key] = model.UserBizRecordListRep{
				ShopName: shopMapIdName[record.ShopId],
				ScanAt:   record.CreatedAt,
				CityName: cityMapIdName[record.CityId],
			}
		}
		response.JsonOkExit(r, list)
	}
	response.JsonOkExit(r, make([]model.UserBizRecordListRep, 0))
}

// BatteryRenewal 骑手扫码换电
// @summary 骑手-骑手扫码换电
// @Accept  json
// @Produce  json
// @tags    骑手-业务办理
// @param   entity  body model.UserBizBatteryRenewalReq true "请求数据"
// @router  /rapi/biz_battery_renewal [POST]
// @success 200 {object} response.JsonResponse{data=model.UserBizBatteryRenewalRep}  "返回结果"
func (*bizApi) BatteryRenewal(r *ghttp.Request) {
	var req model.UserBizBatteryRenewalReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	if user.BatteryState != model.BatteryStateUse {
		response.Json(r, response.RespCodeArgs, "没有正在租借中的电池，不能办理换电")
	}
	shop, err := service.ShopService.DetailByQr(r.Context(), req.Code)
	if err != nil {
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	if shop.State != model.ShopStateOpen {
		response.Json(r, response.RespCodeArgs, "店铺没有营业，不能办理换电")
	}
	at := gtime.Now()
	err = dao.User.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
		if _, err := service.UserBizService.Create(ctx, model.UserBiz{
			ShopId:      shop.Id,
			CityId:      shop.CityId,
			UserId:      user.Id,
			GoroupId:    user.GroupId,
			Type:        model.UserBizBatteryRenewal,
			PackagesId:  user.PackagesId,
			BatteryType: user.BatteryType,
			CreatedAt:   at,
			UpdatedAt:   at,
		}); err != nil {
			return err
		}
		if err := service.UserService.IncrBizBatteryRenewalCnt(ctx, user.Id, 1); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	response.JsonOkExit(r, model.UserBizBatteryRenewalRep{
		ShopName:    shop.Name,
		BatteryType: user.BatteryType,
		At:          at,
	})
}

// New
// @summary 骑手-个签骑手签约之后获取支付信息
// @Accept  json
// @Produce  json
// @tags    骑手-业务办理
// @param   entity  body model.UserBizNewReq true "请求数据"
// @router  /rapi/biz_new [POST]
// @success 200 {object} response.JsonResponse{data=model.UserBizNewRep}"返回结果"
func (*bizApi) New(r *ghttp.Request) {
	var req model.UserBizNewReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	if user.BatteryState != model.BatteryStateDefault && user.BatteryState != model.BatteryStateExit {
		response.Json(r, response.RespCodeArgs, "没有待支付订单")
	}
	if user.GroupId > 0 {
		response.Json(r, response.RespCodeArgs, "团签用户无此操作")
	}
	s, err := service.SignService.UserLatestDetail(r.Context(), user.Id, req.FlowId)
	if err != nil || s == nil {
		response.Json(r, response.RespCodeArgs, "没有完成签约的合同")
		return
	}
	var order model.PackagesOrder
	order, err = service.PackagesOrderService.Detail(r.Context(), s.PackagesOrderId)
	if err == nil {
		packages, _ := service.PackagesService.Detail(r.Context(), order.PackageId)
		switch req.PayType {
		case model.PayTypeWechat:
			var res *app.PrepayWithRequestPaymentResponse
			if res, err = wechat.Service().App(r.Context(), model.Prepay{
				Description: packages.Name,
				No:          order.No,
				Amount:      order.Amount,
				NotifyUrl:   g.Cfg().GetString("api.host") + "/payment_callback/package_new/wechat",
			}); err == nil {
				b, _ := json.Marshal(res)
				response.JsonOkExit(r, model.UserBizNewRep{
					PayOrderInfo: string(b),
				})
				return
			}
			break
		case model.PayTypeAliPay:
			var res string
			if res, err = alipay.Service().App(r.Context(), model.Prepay{
				Description: packages.Name,
				No:          order.No,
				Amount:      order.Amount,
				NotifyUrl:   g.Cfg().GetString("api.host") + "/payment_callback/package_new/alipay",
			}); err == nil {
				response.JsonOkExit(r, model.UserBizNewRep{
					PayOrderInfo: res,
				})
				return
			}
			break
		default:
			err = errors.New("支付方式无效")
			break
		}
	}

	// 经过错误处理之后遇到需要中断的[错误返回]可以直接panic
	panic(err)
}

// Sign 新签
// @summary 骑手-个签用户新签约套餐
// @Accept  json
// @Produce  json
// @tags    骑手-业务办理
// @param   entity  body model.UserBizSignReq true "请求数据"
// @router  /rapi/biz_sign [POST]
// @success 200 {object} response.JsonResponse{data=model.SignRep}  "返回结果"
func (*bizApi) Sign(r *ghttp.Request) {
	var req model.UserBizSignReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	u := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	if u.AuthState != model.AuthStateVerifySuccess {
		response.Json(r, response.RespCodeArgs, "未完成实名认证，请先实名认证")
	}
	if u.BatteryState != model.BatteryStateExit && u.BatteryState != model.BatteryStateDefault {
		response.Json(r, response.RespCodeArgs, "有正在使用中的套餐，请先办理退租")
	}
	if u.GroupId > 0 {
		response.Json(r, response.RespCodeArgs, "团签用户，无需办理购买")
	}
	packages, err := service.PackagesService.Detail(r.Context(), req.PackagesId)
	if err != nil {
		response.Json(r, response.RespCodeArgs, "套餐不存")
	}

	user := service.UserService.Detail(r.Context(), u.Id)
	// 创建代签签文件
	res, err := sign.Service().CreateByTemplate(beansSign.CreateByTemplateReq{
		TemplateId: g.Cfg().GetString("eSign.personal.templateId"),
		SimpleFormFields: beansSign.CreateByTemplateReqSimpleFormFields{
			Name:     user.RealName,
			IdCardNo: user.IdCardNo,
		},
		Name: g.Cfg().GetString("eSign.personal.fileName"),
	})
	if err != nil || res.Code != 0 {
		g.Log().Error(err)
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	// 发起签署
	resFlow, err := sign.Service().CreateFlowOneStep(beansSign.CreateFlowOneStepReq{
		Docs: []beansSign.CreateFlowOneStepReqDoc{
			{
				FileId:   res.Data.FileId,
				FileName: g.Cfg().GetString("eSign.personal.fileName"),
			},
		},
		FlowInfo: beansSign.CreateFlowOneStepReqDocFlowInfo{
			AutoInitiate:  true,
			AutoArchive:   true,
			BusinessScene: g.Cfg().GetString("eSign.personal.businessScene"),
			FlowConfigInfo: beansSign.CreateFlowOneStepReqDocFlowInfoFlowConfigInfo{
				NoticeDeveloperUrl: g.Cfg().GetString("api.host") + "/esign/callback/sign",
				RedirectUrl:        "https://h5.shiguangjv.com/pages/sign-success.html",
			},
		},
		Signers: []beansSign.CreateFlowOneStepReqDocSigner{
			{
				PlatformSign:  true,
				SignerAccount: beansSign.CreateFlowOneStepReqDocSignerAccount{},
				Signfields: []beansSign.CreateFlowOneStepReqDocSignerField{
					{
						AutoExecute: true,
						SignType:    1,
						FileId:      res.Data.FileId,
						PosBean: beansSign.CreateFlowOneStepReqDocSignerFieldPosBean{
							PosPage: "3",
							PosX:    400,
							PosY:    400,
						},
					},
				},
			},
			{
				PlatformSign: false,
				SignerAccount: beansSign.CreateFlowOneStepReqDocSignerAccount{
					SignerAccountId: user.EsignAccountId,
				},
				Signfields: []beansSign.CreateFlowOneStepReqDocSignerField{
					{
						FileId: res.Data.FileId,
						PosBean: beansSign.CreateFlowOneStepReqDocSignerFieldPosBean{
							PosPage: "3",
							PosX:    300,
							PosY:    300,
						},
					},
				},
			},
		},
	})
	if err != nil || resFlow.Code != 0 {
		g.Log().Error(err)
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	// 获取签署地址
	resUrl, err := sign.Service().FlowExecuteUrl(beansSign.FlowExecuteUrlReq{
		FlowId:    resFlow.Data.FlowId,
		AccountId: user.EsignAccountId,
	})
	if err != nil || resUrl.Code != 0 {
		g.Log().Error(err)
		response.JsonErrExit(r, response.RespCodeSystemError)
	}

	if err := dao.PackagesOrder.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
		order, err := service.PackagesOrderService.New(ctx, u.Id, packages)
		if err != nil {
			return err
		}
		if _, _err := service.SignService.Create(ctx, model.Sign{
			UserId:          user.Id,
			GroupId:         0,
			PackagesOrderId: order.Id,
			BatteryType:     packages.BatteryType,
			State:           0,
			FileId:          res.Data.FileId,
			FlowId:          resFlow.Data.FlowId,
		}); _err != nil {
			return _err
		}
		return nil
	}); err != nil {
		g.Log().Error(err)
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	response.JsonOkExit(r, model.SignRep{
		Url:      resUrl.Data.Url,
		ShortUrl: resUrl.Data.ShortUrl,
		FlowId:   resFlow.Data.FlowId,
	})
}

// Renewal 续约
// @summary 骑手-个签用户续购套餐
// @Accept  json
// @Produce  json
// @tags    骑手-业务办理
// @param   entity  body model.UserBizRenewalReq true "请求数据"
// @router  /rapi/biz_renewal [POST]
// @success 200 {object} response.JsonResponse{data=model.UserBizRenewalRep}  "返回结果"
func (*bizApi) Renewal(r *ghttp.Request) {
	var req model.UserBizRenewalReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	if user.GroupId > 0 {
		response.Json(r, response.RespCodeArgs, "团签用户，不能办理续约")
	}
	if user.BatteryState == model.BatteryStateOverdue {
		response.Json(r, response.RespCodeArgs, "套餐已逾期请先交纳违约金")
	}
	if user.BatteryState != model.BatteryStateUse && user.BatteryState != model.BatteryStateSave {
		response.Json(r, response.RespCodeArgs, "没有使用的中套餐，不能办理续约")
	}
	packages, _ := service.PackagesService.Detail(r.Context(), user.PackagesId)
	firstOrder, _ := service.PackagesOrderService.Detail(r.Context(), user.PackagesOrderId)
	order, err := service.PackagesOrderService.Renewal(r.Context(), req.PaymentType, firstOrder)
	if err == nil && req.PaymentType == model.PayTypeWechat {
		if res, err := wechat.Service().App(r.Context(), model.Prepay{
			Description: packages.Name,
			No:          order.No,
			Amount:      order.Amount,
			NotifyUrl:   g.Cfg().GetString("api.host") + "/payment_callback/package_renewal/wechat",
		}); err == nil {
			b, _ := json.Marshal(res)
			response.JsonOkExit(r, model.UserBizNewRep{
				PayOrderInfo: string(b),
			})
		}
	}

	if err == nil && req.PaymentType == model.PayTypeAliPay {
		if res, err := alipay.Service().App(r.Context(), model.Prepay{
			Description: packages.Name,
			No:          order.No,
			Amount:      order.Amount,
			NotifyUrl:   g.Cfg().GetString("api.host") + "/payment_callback/package_renewal/alipay",
		}); err == nil {
			response.JsonOkExit(r, model.UserBizNewRep{
				PayOrderInfo: res,
			})
		}
	}

	response.JsonErrExit(r, response.RespCodeSystemError)
}

// Penalty 违约金
// @summary 骑手-个签用户逾期缴纳违约金
// @Accept  json
// @Produce  json
// @tags    骑手-业务办理
// @param   entity  body model.UserBizPenaltyReq true "请求数据"
// @router  /rapi/biz_penalty [POST]
// @success 200 {object} response.JsonResponse{data=model.UserBizPenaltyRep}  "返回结果"
func (*bizApi) Penalty(r *ghttp.Request) {
	var req model.UserBizPenaltyReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	if user.GroupId > 0 {
		response.Json(r, response.RespCodeArgs, "团签用户无需办理违约")
	}
	if user.BatteryState != model.BatteryStateOverdue {
		response.Json(r, response.RespCodeArgs, "当前套餐正常使用")
	}
	days := carbon.Parse(user.BatteryReturnAt.String()).DiffInDays(carbon.Parse(gtime.Now().String()))
	amount, err := service.PackagesService.PenaltyAmount(r.Context(), user.PackagesId, uint(days))
	if amount <= 0 || err != nil {
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	firstOrder, _ := service.PackagesOrderService.Detail(r.Context(), user.PackagesOrderId)
	packages, _ := service.PackagesService.Detail(r.Context(), user.PackagesId)
	order, err := service.PackagesOrderService.Penalty(r.Context(), req.PaymentType, amount, firstOrder)

	if err == nil && req.PaymentType == model.PayTypeWechat {
		if res, err := wechat.Service().App(r.Context(), model.Prepay{
			Description: packages.Name,
			No:          order.No,
			Amount:      order.Amount,
			NotifyUrl:   g.Cfg().GetString("api.host") + "/payment_callback/package_penalty/wechat",
		}); err == nil {
			b, _ := json.Marshal(res)
			response.JsonOkExit(r, model.UserBizNewRep{
				PayOrderInfo: string(b),
			})
		}
	}

	if err == nil && req.PaymentType == model.PayTypeAliPay {
		if res, err := alipay.Service().App(r.Context(), model.Prepay{
			Description: packages.Name,
			No:          order.No,
			Amount:      order.Amount,
			NotifyUrl:   g.Cfg().GetString("api.host") + "/payment_callback/package_penalty/alipay",
		}); err == nil {
			response.JsonOkExit(r, model.UserBizNewRep{
				PayOrderInfo: res,
			})
		}
	}

	response.JsonErrExit(r, response.RespCodeSystemError)
}

// GroupNew
// @summary 骑手-团签用户新签约
// @Accept  json
// @Produce  json
// @tags    骑手-业务办理
// @param   entity  body model.UserBizGroupNewReq true "请求数据"
// @router  /rapi/biz_new_group [POST]
// @success 200 {object} response.JsonResponse{data=model.SignRep}  "返回结果"
func (*bizApi) GroupNew(r *ghttp.Request) {
	var req model.UserBizGroupNewReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	u := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
	if u.AuthState != model.AuthStateVerifySuccess {
		response.Json(r, response.RespCodeArgs, "未完成实名认证，请先实名认证")
	}
	if u.BatteryState != model.BatteryStateExit && u.BatteryState != model.BatteryStateDefault {
		response.Json(r, response.RespCodeArgs, "有正在使用中的电池，请先办理退租")
	}
	if u.GroupId == 0 {
		response.Json(r, response.RespCodeArgs, "个签用户请购买套餐")
	}
	if u.BatteryState != model.BatteryStateDefault && u.BatteryState != model.BatteryStateExit {
		response.Json(r, response.RespCodeArgs, "已选择过电池型号，请前往店铺办理业务")
	}
	user := service.UserService.Detail(r.Context(), u.Id)
	// 创建代签签文件
	res, err := sign.Service().CreateByTemplate(beansSign.CreateByTemplateReq{
		TemplateId: g.Cfg().GetString("eSign.group.templateId"),
		SimpleFormFields: beansSign.CreateByTemplateReqSimpleFormFields{
			Name:     user.RealName,
			IdCardNo: user.IdCardNo,
		},
		Name: g.Cfg().GetString("eSign.group.fileName"),
	})

	if err != nil || res.Code != 0 {
		fmt.Println(470, err.Error())
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	// 发起签署
	resFlow, err := sign.Service().CreateFlowOneStep(beansSign.CreateFlowOneStepReq{
		Docs: []beansSign.CreateFlowOneStepReqDoc{
			{
				FileId:   res.Data.FileId,
				FileName: g.Cfg().GetString("eSign.group.fileName"),
			},
		},
		FlowInfo: beansSign.CreateFlowOneStepReqDocFlowInfo{
			AutoInitiate:  true,
			AutoArchive:   true,
			BusinessScene: g.Cfg().GetString("eSign.group.businessScene"),
			FlowConfigInfo: beansSign.CreateFlowOneStepReqDocFlowInfoFlowConfigInfo{
				NoticeDeveloperUrl: g.Cfg().GetString("api.host") + "/esign/callback/sign",
				RedirectUrl:        "https://h5.shiguangjv.com/pages/sign-success.html",
			},
		},
		Signers: []beansSign.CreateFlowOneStepReqDocSigner{
			{
				PlatformSign:  true,
				SignerAccount: beansSign.CreateFlowOneStepReqDocSignerAccount{},
				Signfields: []beansSign.CreateFlowOneStepReqDocSignerField{
					{
						AutoExecute: true,
						SignType:    1,
						FileId:      res.Data.FileId,
						PosBean: beansSign.CreateFlowOneStepReqDocSignerFieldPosBean{
							PosPage: "3",
							PosX:    400,
							PosY:    400,
						},
					},
				},
			},
			{
				PlatformSign: false,
				SignerAccount: beansSign.CreateFlowOneStepReqDocSignerAccount{
					SignerAccountId: user.EsignAccountId,
				},
				Signfields: []beansSign.CreateFlowOneStepReqDocSignerField{
					{
						FileId: res.Data.FileId,
						PosBean: beansSign.CreateFlowOneStepReqDocSignerFieldPosBean{
							PosPage: "3",
							PosX:    300,
							PosY:    300,
						},
					},
				},
			},
		},
	})
	if err != nil || res.Code != 0 {
		fmt.Println(525, err)
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	// 获取签署地址
	resUrl, err := sign.Service().FlowExecuteUrl(beansSign.FlowExecuteUrlReq{
		FlowId:    resFlow.Data.FlowId,
		AccountId: user.EsignAccountId,
	})
	if err != nil || res.Code != 0 {
		fmt.Println(534, err.Error())
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	if _, _err := service.SignService.Create(r.Context(), model.Sign{
		UserId:          user.Id,
		GroupId:         user.GroupId,
		PackagesOrderId: 0,
		BatteryType:     req.BatteryType,
		State:           0,
		FileId:          res.Data.FileId,
		FlowId:          resFlow.Data.FlowId,
	}); _err != nil {
		fmt.Println(546, _err.Error())
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	response.JsonOkExit(r, model.SignRep{
		Url:      resUrl.Data.Url,
		ShortUrl: resUrl.Data.ShortUrl,
		FlowId:   resFlow.Data.FlowId,
	})
}
