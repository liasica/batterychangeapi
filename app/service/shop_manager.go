package service

import (
    "battery/app/dao"
    "battery/app/model"
    "context"
    "errors"
    "github.com/gogf/gf/frame/g"
)

var ShopManagerService = shopManagerService{}

type shopManagerService struct {
}

// Login 门店登录
func (s *shopManagerService) Login(ctx context.Context, req model.ShopManagerLoginReq) (rep model.ShopManagerLoginRep, err error) {
    var manager model.ShopManager
    if err = dao.ShopManager.Ctx(ctx).Where(dao.ShopManager.Columns.Mobile, req.Mobile).Scan(&manager); err != nil || manager.Id == 0 || manager.ShopId == 0 {
        err = errors.New("不是门店管理员不能登录")
        return
    }
    if !SmsServer.Verify(ctx, model.SmsVerifyReq{
        Mobile: req.Mobile,
        Code:   req.Sms,
    }) {
        err = errors.New("手机号或验证码错误，登录失败")
        return
    }
    accessToken := dao.ShopManager.GenerateAccessToken(manager.Id)
    _, err = dao.ShopManager.Where(dao.ShopManager.Columns.Id, manager.Id).Update(g.Map{dao.ShopManager.Columns.AccessToken: accessToken})
    rep.AccessToken = accessToken
    return
}

// PushToken 门店修改推送token
func (*shopManagerService) PushToken(ctx context.Context, req model.PushTokenReq) error {
    manager := ctx.Value(model.ContextShopManagerKey).(*model.ContextShopManager)
    _, err := dao.User.WherePri(manager.Id).Update(g.Map{
        dao.ShopManager.Columns.DeviceType:  req.DeviceType,
        dao.ShopManager.Columns.DeviceToken: req.DeviceToken,
    })
    return err
}

// ResetMobile 门店修改手机号码
func (*shopManagerService) ResetMobile(ctx context.Context, req model.ShopManagerResetMobileReq) error {
    manager := ctx.Value(model.ContextShopManagerKey).(*model.ContextShopManager)
    _, err := dao.User.WherePri(manager.Id).Update(g.Map{
        dao.ShopManager.Columns.Mobile: req.Mobile,
    })
    return err
}

// Create 创建
func (*shopManagerService) Create(ctx context.Context, req model.ShopManager) (uint, error) {
    req.AccessToken = dao.ShopManager.GenerateAccessToken(0)
    req.DeviceToken = req.AccessToken
    id, err := dao.ShopManager.Ctx(ctx).InsertAndGetId(req)
    return uint(id), err
}

// Delete 删除
func (*shopManagerService) Delete(ctx context.Context, mobile string) error {
    _, err := dao.ShopManager.Ctx(ctx).Where(dao.ShopManager.Columns.Mobile, mobile).Delete()
    return err
}
