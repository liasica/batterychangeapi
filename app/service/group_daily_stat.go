package service

import (
	"battery/app/dao"
	"context"
	"fmt"
	"github.com/gogf/gf/os/gtime"
)

var GroupDailyStatService = groupDailyStatService{}

type groupDailyStatService struct {
}

func (s *groupDailyStatService) Date(t *gtime.Time) string {
	if t.IsZero() {
		return "0"
	}
	return fmt.Sprintf("%d%d%d", t.Year(), t.Month(), t.Day())
}

// RiderBizNew 用户新领取
func (s *groupDailyStatService) RiderBizNew(ctx context.Context, groupId uint, batteryType uint) error {
	now := gtime.Now()
	_, err := dao.GroupDailyStat.Ctx(ctx).Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		Where(dao.GroupDailyStat.Columns.BatteryType, batteryType).
		WhereGTE(dao.GroupDailyStat.Columns.Date, s.Date(now)).
		Increment(dao.GroupDailyStat.Columns.Total, 1)
	return err
}

type groupDailyStatDateRangeItem struct {
	ArrearsDaysCnt60 uint
	ArrearsDaysCnt72 uint
	PaidDaysCnt60    uint
	PaidDaysCnt72    uint
}

// StatDateRange 按时间统计使用情况
func (s *groupDailyStatService) StatDateRange(ctx context.Context, groupIds []uint, startDate, endDate *gtime.Time) (stat map[uint]*groupDailyStatDateRangeItem, err error) {
	m := dao.GroupDailyStat.Ctx(ctx).
		Fields(dao.GroupDailyStat.Columns.GroupId,
			dao.GroupDailyStat.Columns.BatteryType,
			dao.GroupDailyStat.Columns.IsArrears,
			"sum("+dao.GroupDailyStat.Columns.Total+") as total").
		WhereIn(dao.GroupDailyStat.Columns.GroupId, groupIds)
	if !startDate.IsZero() {
		m = m.WhereGTE(dao.GroupDailyStat.Columns.Date, s.Date(startDate))
	}
	if !endDate.IsZero() {
		m = m.WhereLTE(dao.GroupDailyStat.Columns.Date, s.Date(startDate))
	}
	var res []struct {
		GroupId     uint
		BatteryType uint
		IsArrears   uint
		Total       uint
	}
	err = m.Group(fmt.Sprintf("%s,%s,%s", dao.GroupDailyStat.Columns.GroupId, dao.GroupDailyStat.Columns.BatteryType, dao.GroupDailyStat.Columns.IsArrears)).Scan(&res)
	if err != nil {
		return
	}
	stat = make(map[uint]*groupDailyStatDateRangeItem, len(groupIds))
	for _, groupId := range groupIds {
		stat[groupId] = &groupDailyStatDateRangeItem{}
	}
	for _, row := range res {
		if row.BatteryType == 60 {
			if row.IsArrears == 0 {
				stat[row.GroupId].PaidDaysCnt60 = row.Total
			} else {
				stat[row.GroupId].ArrearsDaysCnt60 = row.Total
			}
		} else {
			if row.IsArrears == 0 {
				stat[row.GroupId].PaidDaysCnt72 = row.Total
			} else {
				stat[row.GroupId].ArrearsDaysCnt72 = row.Total
			}
		}
	}
	return
}
