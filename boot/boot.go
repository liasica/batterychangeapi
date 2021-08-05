package boot

import (
	_ "battery/app/cron"
	"battery/app/service"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/swagger"
)

func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})
	service.MessageService.SendWorkFlowInit()
}
