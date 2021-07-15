package esign

import (
	"github.com/gogf/gf/frame/g"
)

func init() {
	if conf.host == "" {
		conf.SetHost(g.Cfg().GetString("eSign.host"))
		conf.SetProjectId(g.Cfg().GetString("eSign.appId"))
		conf.SetProjectSecret(g.Cfg().GetString("eSign.appSecret"))
	}
}
