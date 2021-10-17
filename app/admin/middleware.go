package admin

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "context"
    "github.com/gogf/gf/net/ghttp"
)

// Context 上下文管理服务
var Context = contextAdmin{}

type contextAdmin struct{}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *contextAdmin) Init(r *ghttp.Request, admin *model.ContextAdmin) {
    r.SetCtxVar(model.ContextAdminKey, admin)
}

// GetUser GetAdmin 获得上下文变量，如果没有设置，那么返回nil
func (s *contextAdmin) GetUser(ctx context.Context) *model.ContextAdmin {
    value := ctx.Value(model.ContextAdminKey)
    if value == nil {
        return nil
    }
    if localCtx, ok := value.(*model.ContextAdmin); ok {
        return localCtx
    }
    return nil
}

// SetUser  将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *contextAdmin) SetUser(ctx context.Context, adm *model.ContextAdmin) {
    admin := s.GetUser(ctx)
    admin.Id = adm.Id
    admin.Username = adm.Username
}

var Middleware = middleware{}

type middleware struct {
}

func (m *middleware) CORS(r *ghttp.Request) {
    corsOptions := r.Response.DefaultCORSOptions()
    corsOptions.AllowDomain = []string{"console.shiguangjv.com"}
    r.Response.CORS(corsOptions)
    r.Middleware.Next()
}

func (m *middleware) Ctx(r *ghttp.Request) {
    admin := &model.ContextAdmin{}
    Context.Init(r, admin)
    if token := r.Header.Get("X-ACCESS-TOKEN"); token != "" {
        if u, err := service.SysUsersService.GetUserByAccessToken(token); err == nil && u.Id > 0 {
            admin.Id = u.Id
            admin.Username = u.Username
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
