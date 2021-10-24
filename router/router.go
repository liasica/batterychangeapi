package router

import (
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "net/http"

    "battery/app/admin"
    "battery/app/api"
    "battery/app/debug"
    "battery/app/esign"
    "battery/app/payment"
    "battery/app/service"
    "battery/app/shop"
    "battery/app/tools"
    "battery/app/user"
)

func init() {
    s := g.Server()
    s.BindMiddlewareDefault(
        service.Middleware.ErrorHandle,
        // service.Middleware.ETAG,
        service.Middleware.CORS,
        // func(r *ghttp.Request) {
        //     r.Middleware.Next()
        //
        //     g.Log().Printf("url: %s \n method: %s \n header: %v \n requestData: %s \n responseCode: %d \n responseData: %s \n",
        //         r.URL.String(),
        //         r.Method,
        //         r.Header,
        //         string(r.GetBody()),
        //         r.Response.Status,
        //         r.Response.BufferString())
        // },
    )

    s.BindStatusHandler(http.StatusNotFound, func(r *ghttp.Request) {

        // TODO 易签回调404
        if r.URL.String() == "/esign/callback/real_name" {
            esign.CallbackApi.RealName(r)
        }

        if r.URL.String() == "/esign/callback/sign" {
            esign.CallbackApi.Sign(r)
        }

    })

    // 认证签约回调
    s.Group("/esign", func(group *ghttp.RouterGroup) {
        group.POST("/callback/real_name", esign.CallbackApi.RealName)
        group.POST("/callback/sign", esign.CallbackApi.Sign)
        group.GET("/:fileId", esign.CallbackApi.SignState)
    })

    s.Group("/debug", func(group *ghttp.RouterGroup) {
        group.GET("/user/reset", debug.User.Reset)
        group.GET("/user/grouptest", debug.User.GroupTest)
        group.GET("/group/settlement", debug.Group.Settlement)
    })

    // 支付回调
    s.Group("/payment_callback", func(group *ghttp.RouterGroup) {
        group.POST("/combo_new/alipay", payment.AlipayApi.ComboOrderNewSuccessCallback)
        group.POST("/combo_new/wechat", payment.WechatApi.ComboOrderNewSuccessCallback)

        group.POST("/combo_renewal/alipay", payment.AlipayApi.ComboOrderRenewalSuccessCallback)
        group.POST("/combo_renewal/wechat", payment.WechatApi.ComboOrderRenewalSuccessCallback)

        group.POST("/combo_penalty/alipay", payment.AlipayApi.ComboOrderPenaltySuccessCallback)
        group.POST("/combo_penalty/wechat", payment.WechatApi.ComboOrderPenaltySuccessCallback)
    })

    // 公用
    s.Group("/api", func(group *ghttp.RouterGroup) {
        group.POST("/upload/image", api.Upload.Image)
        group.POST("/upload/base64_image", api.Upload.Base64Image)
        group.POST("/upload/file", api.Upload.File)
        group.POST("/sms", api.SmsApi.Send)
        group.GET("/battery", api.BatteryApi.Battery)
    })

    // 工具
    s.Group("/tools", func(group *ghttp.RouterGroup) {
        group.GET("/weather", tools.Weather.Now)
    })

    // 骑手
    s.Group("/rapi", func(group *ghttp.RouterGroup) {
        group.POST("/login", user.RiderApi.Login)
        group.Middleware(
            user.Middleware.Ctx,
            user.Middleware.Auth,
        )
        group.POST("/auth", user.RiderApi.Auth)
        group.GET("/auth", user.RiderApi.AuthGet)
        group.PUT("/device", user.RiderApi.PushToken)
        group.GET("/combo", user.RiderApi.Combo)
        group.GET("/combo_order/qr", user.RiderApi.ComboOrderQr)
        group.GET("/home", user.RiderApi.Profile)

        group.GET("/districts/current_city", user.DistrictsApi.CurrentCity)
        group.GET("/open_city", user.DistrictsApi.OpenCityList)

        group.POST("/biz_sign", user.BizApi.Sign)
        group.POST("/biz_new", user.BizApi.New)
        group.GET("/biz_new/:orderId/payState", user.BizApi.NewComborOrderState)
        group.POST("/biz_renewal", user.BizApi.Renewal)
        group.POST("/biz_new_group", user.BizApi.GroupNew)
        group.GET("/biz_penalty", user.BizApi.PenaltyProfile)
        group.POST("/biz_penalty", user.BizApi.Penalty)

        group.GET("/combos", user.ComboApi.List)
        group.GET("/combos/:id", user.ComboApi.Detail)
        group.GET("/shop", user.ShopApi.List)

        group.POST("/biz_battery_renewal", user.BizApi.BatteryRenewal)

        group.GET("/biz_record/stat", user.BizApi.RecordStat)
        group.GET("/biz_record/list", user.BizApi.RecordList)

        group.GET("/group/stat", user.GroupApi.Stat)
        group.GET("/group/list", user.GroupApi.List)

        group.GET("/sign_file", user.RiderApi.SignFile)

        group.GET("/message", user.MessageApi.List)
        group.PUT("/message/read", user.MessageApi.Read)
    })

    // 门店
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
        group.GET("/shop/profile ", shop.ManagerApi.Profile)

        group.GET("/order_scan/:code", shop.OrderApi.ScanDetail)
        group.POST("/order_claim", shop.OrderApi.Claim)

        group.GET("/order_total", shop.OrderApi.Total)
        group.GET("/order", shop.OrderApi.List)
        group.GET("/order/:id", shop.OrderApi.ListDetail)

        group.GET("/user_biz/:code", shop.UserBizApi.Profile)
        group.POST("/user_biz", shop.UserBizApi.Post)
        group.GET("/user_biz_profile/:code", shop.UserBizApi.Profile)

        group.GET("/biz/record", shop.UserBizApi.Record)

        group.GET("/battery", shop.BatteryApi.Overview)
        group.GET("/battery/record", shop.BatteryApi.Record)

        group.POST("/exception", shop.ExceptionApi.Report)
    })

    // 后台管理员
    s.Group("/admin", func(group *ghttp.RouterGroup) {
        group.POST("/login", admin.UserApi.Login)
        group.Middleware(
            admin.Middleware.Ctx,
            admin.Middleware.Auth,
        )
        group.PUT("/logout", admin.UserApi.Logout)
        group.GET("/profile", admin.UserApi.Profile)
        group.GET("/districts/:id/child", admin.DistrictsApi.Child)
        group.GET("/districts", admin.DistrictsApi.List)

        group.Group("/shop", func(g *ghttp.RouterGroup) {
            g.GET("/", admin.ShopApi.List)
            g.POST("/", admin.ShopApi.Create)
            g.GET("/:id", admin.ShopApi.Detail)
            g.PUT("/:id", admin.ShopApi.Edit)
            g.GET("/simpleList", admin.ShopApi.SimpleList)
        })

        group.GET("/combo", admin.ComboApi.List)
        group.POST("/combo", admin.ComboApi.Create)
        group.PUT("/combo/:id", admin.ComboApi.Edit)

        group.Group("/rider", func(g *ghttp.RouterGroup) {
            g.GET("/verify", admin.RiderApi.Verify)
            g.GET("/personal", admin.RiderApi.Personal)
            g.GET("/note/:userId", admin.RiderApi.NoteList)
            g.POST("/note", admin.RiderApi.Note)
        })

        group.Group("/group", func(g *ghttp.RouterGroup) {
            g.GET("/", admin.GroupApi.List)
            g.POST("/", admin.GroupApi.Create)
            g.POST("/:id/member", admin.GroupApi.AddMember)
            g.GET("/:id/member", admin.GroupApi.ListMember)
            g.GET("/:id/member/:userId/biz", admin.GroupApi.ListMemberBiz)
            g.DELETE("/:id/member/:userId", admin.GroupApi.DeleteMember)
            g.GET("/:id/contract", admin.GroupApi.Contract)
            g.GET("/:id/bill/:expDate", admin.GroupApi.GetSettlement)
            g.POST("/bill", admin.GroupApi.PostSettlement)
        })

        group.Group("/battery", func(g *ghttp.RouterGroup) {
            g.GET("/exception", admin.BatteryApi.Exception)
            g.PUT("/exception/{id}", admin.BatteryApi.ExceptionFix)
            g.GET("/record", admin.BatteryApi.TransferRecord)
            g.POST("/allocate", admin.BatteryApi.Allocate)
        })

        group.GET("/order", admin.OrderApi.List)

        group.GET("/biz", admin.BizApi.List)

        group.Group("/dashboard", func(g *ghttp.RouterGroup) {
            g.GET("/cities", admin.DashboardApi.OpenCities)
            g.GET("/overview", admin.DashboardApi.Overview)
            g.GET("/newly", admin.DashboardApi.Newly)
        })
    })

    s.SetIndexFolder(true)
    s.SetServerRoot("./uploads")
    s.AddStaticPath("/uploads", "./uploads")
}
