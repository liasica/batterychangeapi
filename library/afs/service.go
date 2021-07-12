package afs

import (
	afs "github.com/alibabacloud-go/afs-20180112/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/gogf/gf/frame/g"
)

type service struct {
}

var serv *service

func Service() *service {
	if serv == nil {
		serv = &service{}
	}
	return serv
}

func (*service) Verify(sig, sessionId, token, remoteIp string) bool {
	config := new(rpc.Config)
	config.SetAccessKeyId(g.Cfg().GetString("afs.accessKeyId")).
		SetAccessKeySecret(g.Cfg().GetString("afs.accessKeySecret")).
		SetRegionId("cn-hangzhou").
		SetEndpoint("afs.aliyuncs.com")

	client, _ := afs.NewClient(config)

	request := new(afs.AuthenticateSigRequest)
	request.SetSig(sig)
	request.SetSessionId(sessionId)
	request.SetToken(token)
	request.SetRemoteIp(remoteIp)
	request.SetScene(g.Cfg().GetString("afs.scene"))
	request.SetAppKey(g.Cfg().GetString("afs.appKey"))
	response, err := client.AuthenticateSig(request)
	return err == nil && *response.Code == 100
}
