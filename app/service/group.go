package service

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/model/group"
    "battery/library/mq"
    "context"
    "fmt"
    "strings"
)

var GroupService = groupService{}

type groupService struct {
}

func (*groupService) Detail(ctx context.Context, groupId uint) (group model.Group) {
    _ = dao.Group.Ctx(ctx).WherePri(groupId).Scan(&group)
    return group
}

func (*groupService) GetByIds(ctx context.Context, groupIds []uint) (groupList []model.Group) {
    _ = dao.Group.Ctx(ctx).WhereIn(dao.Combo.Columns.Id, groupIds).Scan(&groupList)
    return groupList
}

// CheckName 检测名称是否可以用
func (*groupService) CheckName(ctx context.Context, groupId uint, name string) bool {
    q := dao.Group.Ctx(ctx)
    if groupId > 0 {
        q.WhereNot(dao.Group.Columns.Id, groupId)
    }
    cnt, err := q.Where(dao.Group.Columns.Name, name).Count()
    return cnt == 0 && err == nil
}

// Create 创建
func (*groupService) Create(ctx context.Context, group model.Group) (uint, error) {
    id, err := dao.Group.Ctx(ctx).InsertAndGetId(group)
    return uint(id), err
}

// ListAdmin 管理列表
func (*groupService) ListAdmin(ctx context.Context, req *model.GroupListAdminReq) (total int, items []model.GroupEntity) {
    query := dao.Group.Ctx(ctx).WithAll()
    if req.Keywords != "" {
        query = query.WhereLike(dao.Group.Columns.Name, fmt.Sprintf("%%%s%%", req.Keywords))
    }

    total, _ = query.Count()

    // 查找成员数量
    fileds := mq.FieldsWithTable(group.Table, dao.Group.Columns)
    query = query.Fields(strings.Join(fileds, ",") + ",memberCnt").
        LeftJoin("(SELECT COUNT(1) AS `memberCnt`, `gu`.`groupId` FROM `group_user` `gu` GROUP BY `gu`.`groupId`) `members` ON `members`.groupId = `group`.id")

    _ = query.Page(req.PageIndex, req.PageLimit).Scan(&items)

    for k, item := range items {
        if item.City != nil {
            items[k].CityName = item.City.Name
        }
        for _, detail := range item.SettlementDetails {
            days := detail.GetDays()
            items[k].Days += days
            if detail.State != model.SettlementSettled {
                items[k].BillDays += days
            }
        }
    }
    return
}
