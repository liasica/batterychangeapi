package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"github.com/gogf/gf/frame/g"
)

var GroupUserService = groupUserService{}

type groupUserService struct {
}

// UserCnt 获取团体总人数
func (*groupUserService) UserCnt(ctx context.Context, groupId uint) uint {
	cnt, _ := dao.GroupUser.Ctx(ctx).Where(dao.GroupUser.Columns.GroupId, groupId).Group(dao.GroupUser.Columns.UserId).Count()
	return uint(cnt)
}

// UserIds 获取团体人员ID
func (*groupUserService) UserIds(ctx context.Context, groupId uint) []uint64 {
	var res []struct {
		UserId uint64
	}
	_ = dao.GroupUser.Ctx(ctx).Unscoped().
		Where(dao.GroupUser.Columns.GroupId, groupId).
		Fields(dao.GroupUser.Columns.UserId).
		Scan(&res)
	if len(res) > 0 {
		var ids = make([]uint64, len(res))
		for key, user := range res {
			ids[key] = user.UserId
		}
		return ids
	}
	return make([]uint64, 0)
}

// GetBuyUserId 获取骑手ID获取团体人员信息
func (*groupUserService) GetBuyUserId(ctx context.Context, userId uint64) (groupUser model.GroupUser) {
	_ = dao.GroupUser.Ctx(ctx).
		Where(dao.GroupUser.Columns.UserId, userId).
		Scan(&groupUser)
	return
}

// BatchCreate 创建
func (*groupUserService) BatchCreate(ctx context.Context, userIds []uint64, groupId uint) error {
	data := g.List{}
	for _, userId := range userIds {
		data = append(data, g.Map{
			dao.GroupUser.Columns.UserId:  userId,
			dao.GroupUser.Columns.GroupId: groupId,
		})
	}
	_, err := dao.GroupUser.Ctx(ctx).Data(data).Insert()
	return err
}
