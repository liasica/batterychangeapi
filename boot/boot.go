package boot

import (
    "battery/app/cron"
    _ "battery/app/cron"
    "battery/app/service"
    "battery/library/mongo"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/swagger"
    "sync"
)

func init() {
    s := g.Server()
    s.Plugin(&swagger.Swagger{})

    mongo.Connect()

    o := sync.Once{}
    o.Do(func() {
        _ = cron.RefundCron.Start()
        service.MessageService.SendWorkFlowInit()
    })
}
