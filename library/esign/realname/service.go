package realname

import (
	"battery/library/esign"
	"battery/library/esign/realname/beans"
	"encoding/json"
	"errors"
	"net/http"
)

type service struct {
}

var serv *service

func Service() *service {
	return serv
}

// CreatePersonByThirdPartyUserId 创建个人账号
func (*service) CreatePersonByThirdPartyUserId(info beans.CreatePersonByThirdPartyUserIdInfo) (res beans.CreatePersonByThirdPartyUserIdInfoRes, err error) {
	apiUrl := "/v1/accounts/createByThirdPartyUserId"
	initResult, httpStatus := esign.SendCommHttp(apiUrl, info, "POST")
	if httpStatus == http.StatusOK {
		err = json.Unmarshal(initResult, &res)
	} else {
		err = errors.New("请求失败")
	}
	return
}

// WebIndivIdentityUrl 获取个人实名认证web地址
func (*service) WebIndivIdentityUrl(data beans.WebIndivIdentityUrlInfo, accountId string) (res beans.WebIndivIdentityUrlInfoRes, err error) {
	apiUrl := "/v2/identity/auth/web/" + accountId + "/indivIdentityUrl"
	initResult, httpStatus := esign.SendCommHttp(apiUrl, data, "POST")
	if httpStatus == http.StatusOK {
		err = json.Unmarshal(initResult, &res)
	} else {
		err = errors.New("请求失败")
	}
	return
}
