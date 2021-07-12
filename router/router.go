package router

import (
	"battery/app/admin"
	"battery/app/api"
	"battery/app/esign"
	"battery/app/payment"
	"battery/app/shop"
	"battery/app/user"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()

	//认证签约回调
	s.Group("/esign/callback", func(group *ghttp.RouterGroup) {
		group.Middleware(
			esign.Middleware.Ip,
		)
		group.POST("/real_name", esign.CallbackApi.RealName)
		group.POST("/sign", esign.CallbackApi.Sign)
	})

	//支付回调
	s.Group("/payment_callback", func(group *ghttp.RouterGroup) {
		group.POST("/package_new/alipay", payment.AlipayApi.PackageOrderNewSuccessCallback)
		group.POST("/package_new/wechat", payment.WechatApi.PackageOrderNewSuccessCallback)

		group.POST("/package_renewal/alipay", payment.AlipayApi.PackageOrderRenewalSuccessCallback)
		group.POST("/package_renewal/wechat", payment.WechatApi.PackageOrderRenewalSuccessCallback)

		group.POST("/package_penalty/alipay", payment.AlipayApi.PackageOrderPenaltySuccessCallback)
		group.POST("/package_penalty/wechat", payment.WechatApi.PackageOrderPenaltySuccessCallback)
	})

	//公用
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.POST("/upload/image", api.Upload.Image)
		group.POST("/sms", api.SmsApi.Send)
		group.GET("/test", func(r *ghttp.Request) {
			//resWeb, err := realname.Service().WebIndivIdentityUrl(beans.WebIndivIdentityUrlInfo{
			//	AuthType:           "PSN_FACEAUTH_BYURL",
			//	AvailableAuthTypes: []string{"PSN_FACEAUTH_BYURL"},
			//	ContextInfo: beans.ContextInfo{
			//		NotifyUrl: g.Cfg().GetString("app.host") + "/esign/callback/real_name",
			//	},
			//	ConfigParams: beans.ConfigParams{
			//		IndivUneditableInfo: []string{"name", "certNo", "mobileNo"},
			//	},
			//}, "74b2747324264eb6800cd771e624136f")
			//if err == nil {
			//	_ = r.Response.WriteJson(resWeb)
			//}
			//fmt.Println(resWeb, err)

			//data := beans.CreateByTemplateReq{
			//	Name:       "王麻子的测试合同.pdf",
			//	TemplateId: "a7ac5f9e808840439efb0817ce23b84c",
			//}
			//data.SimpleFormFields.Name = "王麻子"
			//data.SimpleFormFields.IdCardNo = "510902198005183997"
			//res, err := sign.Service().CreateByTemplate(data)
			//fmt.Println(res, err)

			//{0 {https://esignoss.esign.cn/1111564182/89b34614-af9e-481a-ba56-715f2ac6e1ba/%E7%8E%8B%E9%BA%BB%E5%AD%90%E7%9A%84%E6%B5%8B%E8%AF%95%E5%90%88%E5%90%8C.pdf?Expires=1625222633&OSSAccessKeyId=LTAI4G23YViiKnxTC28ygQzF&Signature=MDJVrF3wrQOQhlEuUgM%2BOf%2FZ4w4%3D
			//99f50adc00cb41529c45a8df63c5c6b4
			//王麻子的测试合同.pdf}
			//成功} <nil>

			//data := beans.CreateFlowOneStepReq{
			//	Docs: []beans.CreateFlowOneStepReqDoc{
			//		{
			//			FileId:   "6d91857f3fa04efc959da2436b13e9c0",
			//			FileName: "王麻子的测试合同.pdf",
			//		},
			//	},
			//	FlowInfo: beans.CreateFlowOneStepReqDocFlowInfo{
			//		AutoInitiate:  true,
			//		AutoArchive:   true,
			//		BusinessScene: "本次签署流程的文件主题名称",
			//		FlowConfigInfo: beans.CreateFlowOneStepReqDocFlowInfoFlowConfigInfo{
			//			NoticeDeveloperUrl: g.Cfg().GetString("api.host") + "/esign/callback/sign",
			//		},
			//	},
			//	Signers: []beans.CreateFlowOneStepReqDocSigner{
			//		{
			//			PlatformSign:  true,
			//			SignerAccount: beans.CreateFlowOneStepReqDocSignerAccount{},
			//			Signfields: []beans.CreateFlowOneStepReqDocSignerField{
			//				{
			//					AutoExecute: true,
			//					SignType:    1,
			//					FileId:      "6d91857f3fa04efc959da2436b13e9c0",
			//					PosBean: beans.CreateFlowOneStepReqDocSignerFieldPosBean{
			//						PosPage: "3",
			//						PosX:    400,
			//						PosY:    400,
			//					},
			//				},
			//			},
			//		},
			//		{
			//			PlatformSign: false,
			//			SignerAccount: beans.CreateFlowOneStepReqDocSignerAccount{
			//				SignerAccountId: "74b2747324264eb6800cd771e624136f",
			//			},
			//			Signfields: []beans.CreateFlowOneStepReqDocSignerField{
			//				{
			//					FileId: "6d91857f3fa04efc959da2436b13e9c0",
			//					PosBean: beans.CreateFlowOneStepReqDocSignerFieldPosBean{
			//						PosPage: "3",
			//						PosX:    300,
			//						PosY:    300,
			//					},
			//				},
			//			},
			//		},
			//	},
			//}
			//
			//res, err := sign.Service().CreateFlowOneStep(data)
			//fmt.Println(res, err)

			//{0 成功 {2fd6da656d914c6297891208cfcbb1d1}} <nil>

			//	data := beans.FlowExecuteUrlReq{
			//		FlowId:    "e4c27444a0064fceb0755467a35eb08b",
			//		AccountId: "74b2747324264eb6800cd771e624136f",
			//	}
			//
			//	res, err := sign.Service().FlowExecuteUrl(data)
			//	fmt.Println(res, err)
		})
	})

	//骑手
	s.Group("/rapi", func(group *ghttp.RouterGroup) {
		group.POST("/register", user.UserApi.Register)
		group.POST("/login", user.UserApi.Login)
		group.Middleware(
			user.Middleware.Ctx,
			user.Middleware.Auth,
		)
		group.POST("/auth", user.UserApi.Auth)
		group.GET("/sign", user.UserApi.Sign)
		group.PUT("/device", user.UserApi.PushToken)
		group.GET("/package", user.UserApi.Packages)
		group.GET("/package_order/qr", user.UserApi.PackagesOrderQr)
		group.GET("/home", user.UserApi.Profile)

		group.GET("/districts/current_city", user.DistrictsApi.CurrentCity)
		group.GET("/open_city", user.DistrictsApi.OpenCityList)

		group.POST("/biz_new", user.BizApi.New)
		group.POST("/biz_renewal", user.BizApi.Renewal)
		group.POST("/biz_new_group", user.BizApi.GroupNew)

		group.GET("/packages", user.PackagesApi.List)
		group.GET("/shop", user.ShopApi.List)

		group.POST("/biz_battery_renewal", user.BizApi.BatteryRenewal)

		group.GET("/biz_record/stat", user.BizApi.RecordStat)
		group.GET("/biz_record/list", user.BizApi.RecordList)

		group.GET("/group/stat", user.GroupApi.Stat)
		group.GET("/group/list", user.GroupApi.List)
	})

	//店长
	s.Group("/sapi", func(group *ghttp.RouterGroup) {
		group.POST("/login", shop.ManagerApi.Login)
		group.Middleware(
			shop.Middleware.Ctx,
			shop.Middleware.Auth,
		)
		group.GET("/qr", shop.ManagerApi.Qr)
		group.PUT("/shop/state", shop.ManagerApi.ShopState)
		group.PUT("/shop/device ", shop.ManagerApi.PushToken)
		group.PUT("/shop/mobile ", shop.ManagerApi.ResetMobile)

		group.GET("/order_scan/:code", shop.OrderApi.ScanDetail)
		group.POST("/order_claim", shop.OrderApi.Claim)

		group.GET("/order_total", shop.OrderApi.Total)
		group.GET("/order", shop.OrderApi.List)
		group.GET("/order/:id", shop.OrderApi.ListDetail)

		group.GET("/user_biz/:code", shop.UserBizApi.Profile)
		group.POST("/user_biz", shop.UserBizApi.Post)

		group.GET("biz_record", shop.UserBizApi.RecordUser)
		group.GET("biz_record_total", shop.UserBizApi.RecordUserTotal)

		group.GET("/asset/battery_stat", shop.AssetApi.BatteryStat)
		group.GET("/asset/battery_list", shop.AssetApi.BatteryList)

		group.GET("/exception", shop.ExceptionApi.Report)
	})

	//后台管理员
	s.Group("/adm", func(group *ghttp.RouterGroup) {
		group.Middleware(
			admin.Middleware.CORS,
		)
		group.POST("/login", admin.UserApi.Login)
		group.Middleware(
			admin.Middleware.Ctx,
			admin.Middleware.Auth,
		)
		group.POST("/logout", admin.UserApi.Logout)
		group.GET("/profile", admin.UserApi.Profile)
		group.GET("/districts/:id/child", admin.DistrictsApi.Child)

		group.GET("/shop", admin.ShopApi.List)
		group.POST("/shop", admin.ShopApi.Create)
		group.GET("/shop/:id", admin.ShopApi.Detail)
		group.PUT("/shop/:id", admin.ShopApi.Edit)

		group.GET("/package", admin.PackagesApi.List)
		group.POST("/package", admin.PackagesApi.Create)

		group.GET("/group", admin.GroupApi.List)
		group.POST("/group", admin.GroupApi.Create)
		group.GET("/group/:id", admin.GroupApi.Detail)
		group.PUT("/group/:id", admin.GroupApi.Edit)
	})

}
