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
	if cron.GroupCron.GroupDailyGenerateInit() != nil {
		fmt.Println("团签定时任务初始化失败")
	}
}
