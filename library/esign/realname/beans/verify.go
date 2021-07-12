package beans

// WebIndivIdentityUrlInfo 获取个人实名认证地址请求
type WebIndivIdentityUrlInfo struct {
	AuthType            string       `json:"authType,omitempty"`
	AvailableAuthTypes  []string     `json:"availableAuthTypes,omitempty"`
	AuthAdvancedEnabled []string     `json:"authAdvancedEnabled,omitempty"`
	ReceiveUrlMobileNo  string       `json:"receiveUrlMobileNo,omitempty"`
	ContextInfo         ContextInfo  `json:"contextInfo,omitempty"`
	IndivInfo           IndivInfo    `json:"indivInfo,omitempty"`
	ConfigParams        ConfigParams `json:"configParams,omitempty"`
	RepeatIdentity      bool         `json:"repeatIdentity,omitempty"`
}

type WebIndivIdentityUrlInfoRes struct {
	Code int `json:"code"`
	Data struct {
		FlowId    string `json:"flowId"`
		ShortLink string `json:"shortLink"`
		Url       string `json:"url"`
	} `json:"data"`
	Message string `json:"message"`
}

type ContextInfo struct {
	ContextId      string `json:"contextId,omitempty"`
	NotifyUrl      string `json:"notifyUrl,omitempty"`
	Origin         string `json:"origin,omitempty"`
	RedirectUrl    string `json:"redirectUrl,omitempty"`
	ShowResultPage bool   `json:"showResultPage,omitempty"`
}

type IndivInfo struct {
	BankCardNo string `json:"bankCardNo,omitempty"`
	CertNo     string `json:"certNo,omitempty"`
	CertType   string `json:"certType,omitempty"`
	MobileNo   string `json:"mobileNo,omitempty"`
	Name       string `json:"name,omitempty"`
}

type ConfigParams struct {
	IndivUneditableInfo []string `json:"indivUneditableInfo,omitempty"`
	OrgUneditableInfo   []string `json:"orgUneditableInfo,omitempty"`
}

//IdentityDetail 认证信息查询结果
type IdentityDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		FlowId     string `json:"flowId"`
		Status     string `json:"status"`
		ObjectType string `json:"objectType"`
		AuthType   string `json:"authType"`
		StartTime  int64  `json:"startTime"`
		EndTime    int64  `json:"endTime"`
		OrganInfo  struct {
			AccountId           string `json:"accountId"`
			Name                string `json:"name"`
			CertNo              string `json:"certNo"`
			CertType            string `json:"certType"`
			LegalRepName        string `json:"legalRepName"`
			LegalRepNationality string `json:"legalRepNationality"`
			LegalRepCertNo      string `json:"legalRepCertNo"`
			LegalRepCertType    string `json:"legalRepCertType"`
		} `json:"organInfo"`
		IndivInfo struct {
			AccountId    string `json:"accountId"`
			Name         string `json:"name"`
			CertNo       string `json:"certNo"`
			CertType     string `json:"certType"`
			Nationality  string `json:"nationality"`
			MobileNo     string `json:"mobileNo"`
			BankCardNo   string `json:"bankCardNo"`
			FacePhotoUrl string `json:"facePhotoUrl"`
		} `json:"indivInfo"`
	} `json:"data"`
}
