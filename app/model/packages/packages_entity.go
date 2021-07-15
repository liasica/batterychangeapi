// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package packages

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table packages.
type Entity struct {
    Id          uint        `orm:"id,primary"  json:"id"`           //                         
    DeletedAt   *gtime.Time `orm:"deletedAt"   json:"deleted_at"`   //                         
    CreatedAt   *gtime.Time `orm:"createdAt"   json:"created_at"`   //                         
    UpdatedAt   *gtime.Time `orm:"updatedAt"   json:"updated_at"`   // k                       
    Type        uint        `orm:"type"        json:"type"`         // 套餐类型 1 个人 2 团体  
    BatteryType uint        `orm:"batteryType" json:"battery_type"` // 60 / 72                 
    Name        string      `orm:"name"        json:"name"`         // 名称                    
    Days        uint        `orm:"days"        json:"days"`         // 套餐时长天数            
    Amount      float64     `orm:"amount"      json:"amount"`       // 套餐价格(包含保证金额)  
    Price       float64     `orm:"price"       json:"price"`        //                         
    Earnest     float64     `orm:"earnest"     json:"earnest"`      // 保证金                  
    ProvinceId  uint        `orm:"provinceId"  json:"province_id"`  // 省级行政编码            
    CityId      uint        `orm:"cityId"      json:"city_id"`      // 市级行政编码            
    Packagescol string      `orm:"packagescol" json:"packagescol"`  //                         
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