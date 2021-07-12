package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"errors"
	"github.com/gogf/gf/frame/g"
)

var SysUsersService = sysUsersService{}

type sysUsersService struct {
}

// GetUserByAccessToken 使用accessToken 获取用户信息
func (*sysUsersService) GetUserByAccessToken(accessToken string) (user model.SysUsers, err error) {
	err = dao.SysUsers.Where(dao.SysUsers.Columns.AccessToken, accessToken).Scan(&user)
	return
}

// Login 用户登录
func (s *sysUsersService) Login(ctx context.Context, req model.SysUserLoginReq) (rep model.SysUserLoginRep, err error) {
	var user model.SysUsers
	err = dao.SysUsers.Ctx(ctx).Where(dao.SysUsers.Columns.Username, req.Username).Scan(&user)
	if err != nil || user.Id == 0 {
		err = errors.New("登录失败")
		return
	}
	if dao.SysUsers.EncryptPassword(req.Password, user.Salt) != user.Password {
		err = errors.New("用户名或密码错误登录失败")
		return
	}
	token := dao.SysUsers.GenerateAccessToken(user.Id, user.Salt)
	_, err = dao.SysUsers.Where(dao.SysUsers.Columns.Id, user.Id).Update(g.Map{
		dao.User.Columns.AccessToken: token,
	})
	if err == nil {
		rep.AccessToken = token
	}
	return
}

// Logout 登出
func (s *sysUsersService) Logout(ctx context.Context) (err error) {
	id := ctx.Value(model.ContextAdminKey).(*model.ContextAdmin).Id
	_, err = dao.SysUsers.Where(dao.SysUsers.Columns.Id, id).Update(g.Map{
		dao.SysUsers.Columns.AccessToken: dao.SysUsers.GenerateAccessToken(id, "123"),
	})
	return err
}
