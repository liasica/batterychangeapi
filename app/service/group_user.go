package service

import (
    "battery/app/dao"
    "battery/app/model"
    "context"
    "errors"
    "fmt"
    "github.com/gogf/gf/frame/g"
)

var GroupUserService = groupUserService{}

type groupUserService struct {
}

// AddUsers 添加团签用户
// 联系人（团队管理员是否可用）
func (*groupUserService) AddUsers(ctx context.Context, group model.Group, usersReq []model.GroupCreateUserReq) error {
    if len(usersReq) < 1 {
        return nil
    }

    var userMobiles []string
    var batchInsert []model.User

    userMap := make(map[string]struct{})
    userMap[group.ContactMobile] = struct{}{}

    for _, user := range usersReq {
        if _, exists := userMap[user.Mobile]; exists {
            return errors.New(fmt.Sprintf("手机号 %s 重复", user.Mobile))
        }
        userMap[user.Mobile] = struct{}{}
        userMobiles = append(userMobiles, user.Mobile)

        batchInsert = append(batchInsert, model.User{
            GroupId:  group.Id,
            RealName: user.Name,
            Mobile:   user.Mobile,
            Type:     model.UserTypeGroupRider,
            IdCardNo: user.IdCardNo,
        })
    }

    // 用户状态验证
    users := UserService.GetByMobiles(ctx, userMobiles)
    var usersIds []uint64
    for _, user := range users {
        if user.GroupId > 0 {
            return errors.New(fmt.Sprintf("手机号 %s 已经是其它团签成员，无法添加", user.Mobile))
        }
        if user.BatteryState != model.BatteryStateDefault && user.BatteryState != model.BatteryStateExit {
            return errors.New(fmt.Sprintf("手机号 %s 正在使用的中电池，无法添加", user.Mobile))
        }
        // userMobilesMap[user.Mobile] = false
        usersIds = append(usersIds, user.Id)
    }

    // 批量添加团签用户
    if len(batchInsert) > 0 {
        if err := UserService.AddGroupUsers(ctx, batchInsert); err != nil {
            return err
        }
    }
    return nil
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
