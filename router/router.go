package router

import (
	"battery/app/admin"
	"battery/app/api"
	"battery/app/esign"
	"battery/app/payment"
	"battery/app/service"
	"battery/app/shop"
	"battery/app/user"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	s.BindMiddlewareDefault(service.Middleware.ErrorHandle)

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
		group.PUT("/device", user.UserApi.PushToken)
		group.GET("/package", user.UserApi.Packages)
		group.GET("/package_order/qr", user.UserApi.PackagesOrderQr)
		group.GET("/home", user.UserApi.Profile)

		group.GET("/districts/current_city", user.DistrictsApi.CurrentCity)
		group.GET("/open_city", user.DistrictsApi.OpenCityList)

		group.POST("/biz_sign", user.BizApi.Sign)
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

		group.GET("/message", user.MessageApi.List)
		group.PUT("/message/read", user.MessageApi.Read)
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
