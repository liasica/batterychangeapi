package service

import (
	"github.com/gogf/gf/frame/g"
	umeng_push "github.com/huangfuhui/umeng-push"
)

var PushService = pushService{
	appId:     g.Cfg().GetString("umeng.appKey"),
	appSecret: g.Cfg().GetString("umeng.appMasterSecret"),
}

type pushService struct {
	appId     string
	appSecret string
}

func (s *pushService) UPush() *umeng_push.UmengPush {
	return umeng_push.NewUmengPush(s.appId, s.appSecret)
}

func (s *pushService) Broadcast() error {
	_, err := s.UPush().Send(&umeng_push.SendParam{
		Types: umeng_push.TypeBroadcast,
	})
	//TODO
	return err
}

func (s *pushService) Unicast() error {
	_, err := s.UPush().Send(&umeng_push.SendParam{
		Types: umeng_push.TypeUnicast,
	})
	//TODO
	return err
}
