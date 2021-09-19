package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"errors"
	"fmt"
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
func (s *groupDailyStatService) RiderBizNew(ctx context.Context, groupId uint, batteryType uint, userId uint64) error {
	var max model.GroupDailyStat
	if err := dao.GroupDailyStat.Ctx(ctx).
		Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		Where(dao.GroupDailyStat.Columns.BatteryType, batteryType).
		Scan(&max); err != nil {
		return err
	}
	now := gtime.Now()
	_, err := dao.GroupDailyStat.Ctx(ctx).Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		Where(dao.GroupDailyStat.Columns.BatteryType, batteryType).
		WhereGTE(dao.GroupDailyStat.Columns.Date, s.Date(now)).
		Update(g.Map{
			dao.GroupDailyStat.Columns.Total:   max.Total + 1,
			dao.GroupDailyStat.Columns.UserIds: append(max.UserIds, userId),
		})
	return err
}

// RiderBizExit 用户退租
func (s *groupDailyStatService) RiderBizExit(ctx context.Context, groupId uint, batteryType uint, userId uint64) error {
	var max model.GroupDailyStat
	if err := dao.GroupDailyStat.Ctx(ctx).
		Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		Where(dao.GroupDailyStat.Columns.BatteryType, batteryType).
		Scan(&max); err != nil {
		return err
	}
	now := gtime.Now()
	newUserIds := make([]uint64, 1)
	for _, id := range max.UserIds {
		if id != userId {
			newUserIds = append(newUserIds, id)
		}
	}
	_, err := dao.GroupDailyStat.Ctx(ctx).Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		Where(dao.GroupDailyStat.Columns.BatteryType, batteryType).
		WhereGT(dao.GroupDailyStat.Columns.Date, s.Date(now)).
		Update(g.Map{
			dao.GroupDailyStat.Columns.Total:   max.Total - 1,
			dao.GroupDailyStat.Columns.UserIds: newUserIds,
		})
	return err
}

func (s *groupDailyStatService) GenerateWeek(ctx context.Context, groupId uint, batteryType uint) error {
	var max struct {
		GroupId uint
		UserIds []uint64
		Date    uint
		Total   uint
	}
	maxTime := gtime.NewFromTimeStamp(gtime.Now().Timestamp() - 86400)
	if err := dao.GroupDailyStat.Ctx(ctx).
		Fields(dao.GroupDailyStat.Columns.Date, dao.GroupDailyStat.Columns.UserIds, dao.GroupDailyStat.Columns.Total, dao.GroupDailyStat.Columns.GroupId).
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
			dao.GroupDailyStat.Columns.UserIds:     max.UserIds,
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

type resArrearsDays []struct {
	BatteryType uint
	Cnt         uint
}

// ArrearsDays 获取团队未付款天数
func (s *groupDailyStatService) ArrearsDays(ctx context.Context, groupId uint) (res resArrearsDays, err error) {
	err = dao.GroupDailyStat.Ctx(ctx).Fields("sum(total) as cnt", dao.GroupDailyStat.Columns.BatteryType).
		Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		WhereLTE(dao.GroupDailyStat.Columns.Date, time.Now().Format("20060102")).
		Where(dao.GroupDailyStat.Columns.IsArrears, model.GroupDailyStatIsArrearsYes).
		Group(dao.GroupDailyStat.Columns.BatteryType).
		Scan(&res)
	return
}

// ArrearsList 获取团队未付款统计记录
func (s *groupDailyStatService) ArrearsList(ctx context.Context, groupId uint) (list []model.GroupDailyStat, err error) {
	err = dao.GroupDailyStat.Ctx(ctx).
		Where(dao.GroupDailyStat.Columns.GroupId, groupId).
		WhereLTE(dao.GroupDailyStat.Columns.Date, time.Now().Format("20060102")).
		Where(dao.GroupDailyStat.Columns.IsArrears, model.GroupDailyStatIsArrearsYes).
		Scan(&list)
	return
}
