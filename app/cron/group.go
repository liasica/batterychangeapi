package cron

import (
	"battery/app/model"
	"battery/app/service"
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/robfig/cron/v3"
)

func init() {
	if GroupCron.start() != nil {
		panic("团签定时任务初始化失败")
	}

	if RefundCron.start() != nil {
		panic("退款定时任务初始化失败")
	}
}

var GroupCron = group{}

type group struct {
}

func (*group) start() error {
	if !g.Cfg().GetBool("cron.group.stat.enable", false) {
		return nil
	}
	c := cron.New()
	_, err := c.AddFunc(g.Cfg().GetString("cron.group.stat.spec"), func() {
		fmt.Println("GroupDailyGenerateInit start !!!")
		req := model.GroupListAdminReq{}
		req.PageIndex = 1
		req.PageLimit = 20
		for {
			total, items := service.GroupService.ListAdmin(context.TODO(), req)
			if total == 0 || len(items) == 0 {
				break
			}
			for _, item := range items {
				if service.GroupDailyStatService.GenerateWeek(context.TODO(), item.Id, 60) != nil {
					g.Log().Error(fmt.Sprintf("团体 %d 型号 %d 生成失败", item.Id, 60))
				}
				if service.GroupDailyStatService.GenerateWeek(context.TODO(), item.Id, 72) != nil {
					g.Log().Error(fmt.Sprintf("团体 %d 型号 %d 生成失败", item.Id, 72))
				}
			}
			req.PageIndex++
		}
	})
	c.Start()

	return err
}
