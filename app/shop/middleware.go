package shop

import (
	"battery/app/dao"
	"battery/app/model"
	"battery/library/response"
	"context"
	"github.com/gogf/gf/net/ghttp"
)

// Context 上下文管理服务
var Context = contextShop{}

type contextShop struct{}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *contextShop) Init(r *ghttp.Request, manager *model.ContextShopManager) {
	r.SetCtxVar(model.ContextShopManagerKey, manager)
}

// GetManager 获得上下文变量，如果没有设置，那么返回nil
func (s *contextShop) GetManager(ctx context.Context) *model.ContextShopManager {
	value := ctx.Value(model.ContextShopManagerKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.ContextShopManager); ok {
		return localCtx
	}
	return nil
}

// SetManager 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *contextShop) SetManager(ctx context.Context, r *model.ContextShopManager) {
	manager := s.GetManager(ctx)
	r.Id = manager.Id
	r.ShopId = manager.ShopId
	r.Mobile = manager.Mobile
}

var Middleware = middleware{}

type middleware struct {
}

func (m *middleware) Ctx(r *ghttp.Request) {
	manager := &model.ContextShopManager{}
	Context.Init(r, manager)
	if token := r.Header.Get("X-ACCESS-TOKEN"); token != "" {
		var u model.ShopManager
		if err := dao.ShopManager.Where(dao.User.Columns.AccessToken, token).Scan(&u); err == nil && u.Id > 0 {
			manager.Id = u.Id
			manager.ShopId = u.ShopId
			manager.Mobile = u.Mobile
		}
	}
	r.Middleware.Next()
}

// Auth 鉴权中间件，只有登录成功之后才能通过
func (m *middleware) Auth(r *ghttp.Request) {
	if user := Context.GetManager(r.Context()); user != nil && user.ShopId > 0 {
		r.Middleware.Next()
	} else {
		response.JsonErrExit(r, response.RespCodeUnauthorized)
	}
}
