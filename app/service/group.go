package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"fmt"
)

var GroupService = groupService{}

type groupService struct {
}

// StatDays 获取团体总使用天数
func (*groupService) StatDays(ctx context.Context, groupId uint) uint {
	days, _ := dao.GroupDailyStat.Ctx(ctx).Where(dao.GroupDailyStat.Columns.GroupId, groupId).Sum(dao.GroupDailyStat.Columns.Total)
	return uint(days)
}

func (*groupService) Detail(ctx context.Context, groupId uint) (group model.Group) {
	_ = dao.Group.Ctx(ctx).WherePri(groupId).Scan(&group)
	return group
}

func (*groupService) GetByIds(ctx context.Context, groupIds []uint) (groupList []model.Group) {
	_ = dao.Group.Ctx(ctx).WhereIn(dao.Packages.Columns.Id, groupIds).Scan(&groupList)
	return groupList
}

// CheckName 检测名称是否可以用
func (*groupService) CheckName(ctx context.Context, groupId uint, name string) bool {
	var cnt int
	var err error
	if groupId > 0 {
		cnt, err = dao.Group.Ctx(ctx).WherePri(dao.Group.Columns.Name, name).Count()
	} else {
		cnt, err = dao.Group.Ctx(ctx).WherePri(groupId).Where(dao.Group.Columns.Name, name).Count()
	}
	return cnt == 0 && err == nil
}

// Create 创建
func (*groupService) Create(ctx context.Context, group model.Group) (uint, error) {
	id, err := dao.Group.Ctx(ctx).InsertAndGetId(group)
	return uint(id), err
}

// ListAdmin 管理列表
func (*groupService) ListAdmin(ctx context.Context, req model.GroupListAdminReq) (total int, items []model.Group) {
	m := dao.Group.Ctx(ctx)
	if req.Keywords != "" {
		m = m.WhereLike(dao.Group.Columns.Name, fmt.Sprintf("%%%s%%", req.Keywords))
	}
	total, _ = m.Count()
	if total > 0 {
		_ = m.Page(req.PageIndex, req.PageLimit).Scan(&items)
	}
	return
}
