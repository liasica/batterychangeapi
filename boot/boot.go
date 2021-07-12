package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/swagger"
)

func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})
}
