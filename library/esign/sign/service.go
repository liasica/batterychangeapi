package sign

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"battery/library/esign"
	"battery/library/esign/sign/beans"
)

type service struct {
}

var serv *service

func Service() *service {
	return serv
}

// CreateByTemplate 使用模板创建文件
func (*service) CreateByTemplate(req beans.CreateByTemplateReq) (res beans.CreateByTemplateRep, err error) {
	apiUrl := "/v1/files/createByTemplate"
	initResult, httpStatus := esign.SendCommHttp(apiUrl, req, "POST")
	if httpStatus == http.StatusOK {
		err = json.Unmarshal(initResult, &res)
	} else {
		err = errors.New("请求失败")
	}
	return
}

// CreateFlowOneStep 一步发起签署
func (*service) CreateFlowOneStep(req beans.CreateFlowOneStepReq) (res beans.CreateFlowOneStepRep, err error) {
	apiUrl := "/api/v2/signflows/createFlowOneStep"
	initResult, httpStatus := esign.SendCommHttp(apiUrl, req, "POST")
	if httpStatus == http.StatusOK {
		err = json.Unmarshal(initResult, &res)
	} else {
		err = errors.New("请求失败")
	}
	return
}

// FlowExecuteUrl 获取签约Url
func (*service) FlowExecuteUrl(req beans.FlowExecuteUrlReq) (res beans.FlowExecuteUrlRep, err error) {
	apiUrl := fmt.Sprintf("/v1/signflows/%s/executeUrl?accountId=%s", req.FlowId, req.AccountId)
	initResult, httpStatus := esign.SendCommHttp(apiUrl, nil, "GET")
	if httpStatus == http.StatusOK {
		err = json.Unmarshal(initResult, &res)
	} else {
		err = errors.New("请求失败")
	}
	return
}

// SignFlowDocuments 获取签约文件地址
func (*service) SignFlowDocuments(flowId string) (res beans.SignFlowDocumentsRep, err error) {
	apiUrl := fmt.Sprintf("/v1/signflows/{flowId}/documents", flowId)
	initResult, httpStatus := esign.SendCommHttp(apiUrl, nil, "GET")
	if httpStatus == http.StatusOK {
		err = json.Unmarshal(initResult, &res)
	} else {
		err = errors.New("请求失败")
	}
	return
}
