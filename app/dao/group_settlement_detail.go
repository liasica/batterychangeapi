// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"battery/app/dao/internal"
)

// groupSettlementDetailDao is the manager for logic model data accessing
// and custom defined data operations functions management. You can define
// methods on it to extend its functionality as you wish.
type groupSettlementDetailDao struct {
	*internal.GroupSettlementDetailDao
}

var (
	// GroupSettlementDetail is globally public accessible object for table group_settlement_detail operations.
	GroupSettlementDetail groupSettlementDetailDao
)

func init() {
	GroupSettlementDetail = groupSettlementDetailDao{
		internal.NewGroupSettlementDetailDao(),
	}
}

// Fill with you ideas below.