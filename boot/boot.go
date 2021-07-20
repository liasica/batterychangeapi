package boot

import (
	"battery/app/cron"
	"battery/app/service"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/swagger"
)

func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})
	if cron.GroupCron.Init() != nil {
		panic("团签定时任务初始化失败")
	}

	if cron.RefundCron.Init() != nil {
		panic("退款定时任务初始化失败")
	}
	service.MessageService.SendWorkFlowInit()
}
