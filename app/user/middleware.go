package user

import (
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
	"context"
	"github.com/gogf/gf/net/ghttp"
)

// Context 上下文管理服务
var Context = contextUser{}

type contextUser struct{}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *contextUser) Init(r *ghttp.Request, rider *model.ContextRider) {
	r.SetCtxVar(model.ContextRiderKey, rider)
}

// GetUser GetAdmin 获得上下文变量，如果没有设置，那么返回nil
func (s *contextUser) GetUser(ctx context.Context) *model.ContextRider {
	value := ctx.Value(model.ContextRiderKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.ContextRider); ok {
		return localCtx
	}
	return nil
}

// SetUser  将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *contextUser) SetUser(ctx context.Context, r *model.ContextRider) {
	rider := s.GetUser(ctx)
	rider.Id = r.Id
	rider.Mobile = r.Mobile
	rider.GroupId = r.GroupId
	rider.EsignAccountId = r.EsignAccountId
	rider.AuthState = r.AuthState
	rider.SignState = r.SignState
	rider.PackagesId = r.PackagesId
	rider.BatteryType = r.BatteryType
	rider.Qr = r.Qr
	rider.PackagesOrderId = r.PackagesOrderId
	rider.ExpirationAt = r.ExpirationAt
	rider.BizBatteryRenewalCnt = r.BizBatteryRenewalCnt
	rider.BizBatteryRenewalSeconds = r.BizBatteryRenewalSeconds
	rider.BizBatterySecondsStartAt = r.BizBatterySecondsStartAt
}

var Middleware = middleware{}

type middleware struct {
}

func (m *middleware) Ctx(r *ghttp.Request) {
	rider := &model.ContextRider{}
	Context.Init(r, rider)
	if token := r.Header.Get("X-ACCESS-TOKEN"); token != "" {
		if u, err := service.UserService.GetUserByAccessToken(token); err == nil && u.Id > 0 {
			rider.Id = u.Id
			rider.Mobile = u.Mobile
			rider.GroupId = u.GroupId
			rider.EsignAccountId = u.EsignAccountId
			rider.AuthState = u.AuthState
			rider.SignState = u.SignState
			rider.PackagesId = u.PackagesId
			rider.BatteryType = u.BatteryType
			rider.BatteryState = u.BatteryState
			rider.Qr = u.Qr
			rider.PackagesOrderId = u.PackagesOrderId
			rider.ExpirationAt = u.BatteryReturnAt
			rider.BizBatteryRenewalCnt = u.BizBatteryRenewalCnt
			rider.BizBatteryRenewalSeconds = u.BizBatteryRenewalSeconds
			rider.BizBatterySecondsStartAt = u.BizBatterySecondsStartAt
		}
	}
	r.Middleware.Next()
}

// Auth 鉴权中间件，只有登录成功之后才能通过
func (m *middleware) Auth(r *ghttp.Request) {
	if user := Context.GetUser(r.Context()); user != nil && user.Id > 0 {
		r.Middleware.Next()
	} else {
		response.JsonErrExit(r, response.RespCodeUnauthorized)
	}
}
