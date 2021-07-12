package service

import (
	"battery/app/dao"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"strconv"
	"time"
)

var GroupDailyStatService = groupDailyStatService{}

type groupDailyStatService struct {
}

func (s *groupDailyStatService) Date(t *gtime.Time) string {
	if t.IsZero() {
		return "0"
	}
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
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

// RiderBizExit 用户退租
func (s *groupDailyStatService) RiderBizExit(ctx context.Context, groupId uint, batteryType uint) error {
	now := gtime.Now()
	_, err := dao.GroupDailyStat.Ctx(ctx).Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		Where(dao.GroupDailyStat.Columns.BatteryType, batteryType).
		WhereGTE(dao.GroupDailyStat.Columns.Date, s.Date(now)).
		Decrement(dao.GroupDailyStat.Columns.Total, 1)
	return err
}

func (s *groupDailyStatService) GenerateWeek(ctx context.Context, groupId uint, batteryType uint) error {
	var max struct {
		GroupId uint
		Date    uint
		Total   uint
	}
	return dao.GroupDailyStat.DB.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		maxTime := gtime.NewFromTimeStamp(gtime.Now().Timestamp() - 86400)
		if err := dao.GroupDailyStat.Ctx(ctx).
			Fields(dao.GroupDailyStat.Columns.Date, dao.GroupDailyStat.Columns.Total, dao.GroupDailyStat.Columns.GroupId).
			Where(dao.GroupDailyStat.Columns.GroupId, groupId).
			Where(dao.GroupDailyStat.Columns.BatteryType, batteryType).
			OrderDesc(dao.GroupDailyStat.Columns.Id).
			LockUpdate().
			Scan(&max); err == nil {
			today, _ := strconv.ParseUint(s.Date(gtime.Now()), 10, 64)
			if max.Date-uint(today) < 7 {
				maxTime = gtime.NewFromStr(fmt.Sprintf("%d-%02d-%02d 23:59:59", max.Date/10000, max.Date%10000/100, max.Date%100))
			} else {
				return nil
			}
		}
		newMaxTime := gtime.Now().Add(24 * 7 * time.Hour)
		newStat := g.List{}
		for {
			if maxTime.Timestamp() >= newMaxTime.Timestamp() {
				break
			}
			maxTime = maxTime.Add(24 * time.Hour)
			newStat = append(newStat, g.Map{
				dao.GroupDailyStat.Columns.GroupId:     groupId,
				dao.GroupDailyStat.Columns.BatteryType: batteryType,
				dao.GroupDailyStat.Columns.Date:        s.Date(maxTime),
				dao.GroupDailyStat.Columns.Total:       max.Total,
			})
		}
		if len(newStat) > 0 {
			res, err := dao.GroupDailyStat.Ctx(ctx).Data(newStat).Insert()
			if err != nil {
				return err
			}
			if row, _err := res.RowsAffected(); int(row) != len(newStat) || _err != nil {
				return errors.New("执行失败")
			}
		}
		return nil
	})
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
