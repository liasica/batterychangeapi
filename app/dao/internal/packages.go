// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// PackagesDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type PackagesDao struct {
	gmvc.M                  // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB          // DB is the raw underlying database management object.
	Table   string          // Table is the table name of the DAO.
	Columns packagesColumns // Columns contains all the columns of Table that for convenient usage.
}

// PackagesColumns defines and stores column names for table packages.
type packagesColumns struct {
	Id          string //
	DeletedAt   string //
	CreatedAt   string //
	UpdatedAt   string // k
	Type        string // 套餐类型 1 个人 2 团体
	BatteryType string // 60 / 72
	Name        string // 名称
	Days        string // 套餐时长天数
	Amount      string // 套餐价格(包含保证金额)
	Price       string //
	Earnest     string // 保证金
	ProvinceId  string // 省级行政编码
	CityId      string // 市级行政编码
	Packagescol string //
}

func NewPackagesDao() *PackagesDao {
	return &PackagesDao{
		M:     g.DB("default").Model("packages").Safe(),
		DB:    g.DB("default"),
		Table: "packages",
		Columns: packagesColumns{
			Id:          "id",
			DeletedAt:   "deletedAt",
			CreatedAt:   "createdAt",
			UpdatedAt:   "updatedAt",
			Type:        "type",
			BatteryType: "batteryType",
			Name:        "name",
			Days:        "days",
			Amount:      "amount",
			Price:       "price",
			Earnest:     "earnest",
			ProvinceId:  "provinceId",
			CityId:      "cityId",
			Packagescol: "packagescol",
		},
	}
}
