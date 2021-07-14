package boot

import (
	"battery/app/cron"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/swagger"
)

func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})
	if cron.GroupCron.Init() != nil {
		fmt.Println("团签定时任务初始化失败")
	}

	if cron.RefundCron.Init() != nil {
		fmt.Println("退款定时任务初始化失败")
	}
}
