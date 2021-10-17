// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package user_biz

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table user_biz.
type Entity struct {
    Id           uint64      `orm:"id,primary"   json:"id"`             //                                             
    CityId       uint        `orm:"cityId"       json:"city_id"`        //                                             
    ShopId       uint        `orm:"shopId"       json:"shop_id"`        // 门店ID                                      
    UserId       uint64      `orm:"userId"       json:"user_id"`        //                                             
    GoroupId     uint        `orm:"goroupId"     json:"goroup_id"`      // 团体ID                                      
    GoroupUserId uint        `orm:"goroupUserId" json:"goroup_user_id"` //                                             
    Type         uint        `orm:"type"         json:"type"`           // 业务类型:  1 新签， 2 换电，3 寄存，4 退租  
    PackagesId   uint        `orm:"packagesId"   json:"packages_id"`    // 套餐ID                                      
    BatteryType  uint        `orm:"batteryType"  json:"battery_type"`   // 电池型号 60 / 72                            
    CreatedAt    *gtime.Time `orm:"createdAt"    json:"created_at"`     // 扫码时间                                    
    UpdatedAt    *gtime.Time `orm:"updatedAt"    json:"updated_at"`     //                                             
}

// OmitEmpty sets OPTION_OMITEMPTY option for the model, which automatically filers
// the data and where attributes for empty values.
// Deprecated.
func (r *Entity) OmitEmpty() *arModel {
	return Model.Data(r).OmitEmpty()
}

// Inserts does "INSERT...INTO..." statement for inserting current object into table.
// Deprecated.
func (r *Entity) Insert() (result sql.Result, err error) {
	return Model.Data(r).Insert()
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for inserting current object into table.
// Deprecated.
func (r *Entity) InsertIgnore() (result sql.Result, err error) {
	return Model.Data(r).InsertIgnore()
}

// Replace does "REPLACE...INTO..." statement for inserting current object into table.
// If there's already another same record in the table (it checks using primary key or unique index),
// it deletes it and insert this one.
// Deprecated.
func (r *Entity) Replace() (result sql.Result, err error) {
	return Model.Data(r).Replace()
}

// Save does "INSERT...INTO..." statement for inserting/updating current object into table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
// Deprecated.
func (r *Entity) Save() (result sql.Result, err error) {
	return Model.Data(r).Save()
}

// Update does "UPDATE...WHERE..." statement for updating current object from table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
// Deprecated.
func (r *Entity) Update() (result sql.Result, err error) {
	where, args, err := gdb.GetWhereConditionOfStruct(r)
	if err != nil {
		return nil, err
	}
	return Model.Data(r).Where(where, args).Update()
}

// Delete does "DELETE FROM...WHERE..." statement for deleting current object from table.
// Deprecated.
func (r *Entity) Delete() (result sql.Result, err error) {
	where, args, err := gdb.GetWhereConditionOfStruct(r)
	if err != nil {
		return nil, err
	}
	return Model.Where(where, args).Delete()
}