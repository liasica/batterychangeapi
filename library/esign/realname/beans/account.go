package beans

// CreatePersonByThirdPartyUserIdInfo 创建个人账号json信息配置
type CreatePersonByThirdPartyUserIdInfo struct {
	ThirdPartyUserId string `json:"thirdPartyUserId,omitempty"`
	Name             string `json:"name,omitempty"`
	IdType           string `json:"idType,omitempty"`
	IdNumber         string `json:"idNumber,omitempty"`
	Mobile           string `json:"mobile,omitempty"`
	Email            string `json:"email,omitempty"`
}

type CreatePersonByThirdPartyUserIdInfoRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccountId string `json:"accountId"`
	} `json:"data"`
}

// CreateOrganizationsByThirdPartyUserIdInfo 创建企业账号json信息配置
type CreateOrganizationsByThirdPartyUserIdInfo struct {
	ThirdPartyUserId string `json:"thirdPartyUserId,omitempty"`
	Creator          string `json:"creator,omitempty"`
	Name             string `json:"name,omitempty"`
	IdType           string `json:"idType,omitempty"`
	IdNumber         string `json:"idNumber,omitempty"`
	OrgLegalIdNumber string `json:"orgLegalIdNumber,omitempty"`
	OrgLegalName     string `json:"orgLegalName,omitempty"`
}
