// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// GroupDailyStatDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type GroupDailyStatDao struct {
	gmvc.M                        // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB                // DB is the raw underlying database management object.
	Table   string                // Table is the table name of the DAO.
	Columns groupDailyStatColumns // Columns contains all the columns of Table that for convenient usage.
}

// GroupDailyStatColumns defines and stores column names for table group_daily_stat.
type groupDailyStatColumns struct {
	Id          string //
	GroupId     string //
	BatteryType string // 电池型号 60 / 72
	IsArrears   string // 是否未付款 1 是 0 不是
	Date        string // 日期 如 20210705
	Total       string // 使用人数
	UserIds     string //
	CreatedAt   string // 创建时间
	UpdatedAt   string //
}

func NewGroupDailyStatDao() *GroupDailyStatDao {
	return &GroupDailyStatDao{
		M:     g.DB("default").Model("group_daily_stat").Safe(),
		DB:    g.DB("default"),
		Table: "group_daily_stat",
		Columns: groupDailyStatColumns{
			Id:          "id",
			GroupId:     "groupId",
			BatteryType: "batteryType",
			IsArrears:   "isArrears",
			Date:        "date",
			Total:       "total",
			UserIds:     "userIds",
			CreatedAt:   "createdAt",
			UpdatedAt:   "updatedAt",
		},
	}
}
