package beans

// CreateByTemplateReq 通过模板创建文件请求
type CreateByTemplateReq struct {
	Name             string                              `json:"name,omitempty"`
	SimpleFormFields CreateByTemplateReqSimpleFormFields `json:"simpleFormFields,omitempty"`
	TemplateId       string                              `json:"templateId,omitempty"`
}

type CreateByTemplateReqSimpleFormFields struct {
	Name     string `json:"name,omitempty"`
	IdCardNo string `json:"idCardNo,omitempty"`
}

// CreateByTemplateRep 通过模板创建文件响应
type CreateByTemplateRep struct {
	Code int `json:"code,omitempty"`
	Data struct {
		DownloadUrl string `json:"downloadUrl,omitempty"`
		FileId      string `json:"fileId,omitempty"`
		FileName    string `json:"fileName,omitempty"`
	} `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}
