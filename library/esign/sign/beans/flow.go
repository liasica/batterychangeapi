package beans

// CreateFlowOneStepReq 一步发起签署请求
type CreateFlowOneStepReq struct {
	Attachments []struct {
		AttachmentName string `json:"attachmentName,omitempty"`
		FileId         string `json:"fileId,omitempty"`
	} `json:"attachments,omitempty"`
	Docs    []CreateFlowOneStepReqDoc `json:"docs,omitempty"`
	Copiers []struct {
		CopierAccountId           string `json:"copierAccountId,omitempty"`
		CopierIdentityAccountId   string `json:"copierIdentityAccountId,omitempty"`
		CopierIdentityAccountType int    `json:"copierIdentityAccountType,omitempty"`
	} `json:"copiers,omitempty"`
	FlowInfo CreateFlowOneStepReqDocFlowInfo `json:"flowInfo,omitempty"`
	Signers  []CreateFlowOneStepReqDocSigner `json:"signers,omitempty"`
}

type CreateFlowOneStepReqDoc struct {
	FileId   string `json:"fileId,omitempty"`
	FileName string `json:"fileName,omitempty"`
}

type CreateFlowOneStepReqDocFlowInfo struct {
	AutoArchive      bool                                          `json:"autoArchive,omitempty"`
	AutoInitiate     bool                                          `json:"autoInitiate,omitempty"`
	BusinessScene    string                                        `json:"businessScene,omitempty"`
	ContractRemind   int                                           `json:"contractRemind,omitempty"`
	ContractValidity int64                                         `json:"contractValidity,omitempty"`
	FlowConfigInfo   CreateFlowOneStepReqDocFlowInfoFlowConfigInfo `json:"flowConfigInfo,omitempty"`
	Remark           string                                        `json:"remark,omitempty"`
	SignValidity     int64                                         `json:"signValidity,omitempty"`
}

type CreateFlowOneStepReqDocFlowInfoFlowConfigInfo struct {
	NoticeDeveloperUrl        string   `json:"noticeDeveloperUrl,omitempty"`
	NoticeType                string   `json:"noticeType,omitempty"`
	RedirectUrl               string   `json:"redirectUrl,omitempty"`
	SignPlatform              string   `json:"signPlatform,omitempty"`
	WillTypes                 []string `json:"willTypes,omitempty"`
	PersonAvailableAuthTypes  []string `json:"personAvailableAuthTypes,omitempty"`
	BatchDropSeal             bool     `json:"batchDropSeal,omitempty"`
	OrgAvailableAuthTypes     []string `json:"orgAvailableAuthTypes,omitempty"`
	PersonAuthAdvancedEnabled []string `json:"personAuthAdvancedEnabled,omitempty"`
	Countdown                 int      `json:"countdown,omitempty"`
}

type CreateFlowOneStepReqDocSigner struct {
	PlatformSign  bool                                 `json:"platformSign,omitempty"`
	SignOrder     int                                  `json:"signOrder,omitempty"`
	SignerAccount CreateFlowOneStepReqDocSignerAccount `json:"signerAccount,omitempty"`
	Signfields    []CreateFlowOneStepReqDocSignerField `json:"signfields,omitempty"`
	ThirdOrderNo  string                               `json:"thirdOrderNo,omitempty"`
}

type CreateFlowOneStepReqDocSignerAccount struct {
	SignerAccountId     string `json:"signerAccountId,omitempty"`
	AuthorizedAccountId string `json:"authorizedAccountId,omitempty"`
	NoticeType          string `json:"noticeType,omitempty"`
}

type CreateFlowOneStepReqDocSignerField struct {
	AutoExecute        bool                                           `json:"autoExecute,omitempty"`
	ActorIndentityType int                                            `json:"actorIndentityType,omitempty"`
	FileId             string                                         `json:"fileId,omitempty"`
	PosBean            CreateFlowOneStepReqDocSignerFieldPosBean      `json:"posBean,omitempty"`
	SealType           string                                         `json:"sealType,omitempty"`
	SignDateBean       CreateFlowOneStepReqDocSignerFieldSignDateBean `json:"signDateBean,omitempty"`
	SignType           int                                            `json:"signType,omitempty"`
	Width              int                                            `json:"width,omitempty"`
}
type CreateFlowOneStepReqDocSignerFieldPosBean struct {
	PosPage string `json:"posPage,omitempty"`
	PosX    int    `json:"posX,omitempty"`
	PosY    int    `json:"posY,omitempty"`
}

type CreateFlowOneStepReqDocSignerFieldSignDateBean struct {
	FontSize int    `json:"fontSize,omitempty"`
	Format   string `json:"format,omitempty"`
}

// CreateFlowOneStepRep 一步发起签署请求
type CreateFlowOneStepRep struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		FlowId string `json:"flowId,omitempty"`
	} `json:"data,omitempty"`
}

// FlowExecuteUrlReq 获取签署url请求
type FlowExecuteUrlReq struct {
	FlowId     string `json:"flowId,omitempty"`
	AccountId  string `json:"accountId,omitempty"`
	OrganizeId string `json:"organizeId,omitempty"`
	UrlType    string `json:"urlType,omitempty"`
	AppScheme  string `json:"appScheme,omitempty"`
}

// FlowExecuteUrlRep 获取签署url响应
type FlowExecuteUrlRep struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Url      string `json:"url,omitempty"`
		ShortUrl string `json:"shortUrl,omitempty"`
	} `json:"data,omitempty"`
}
