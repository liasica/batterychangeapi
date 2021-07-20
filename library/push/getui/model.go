package getui

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type TokenRequest struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	AppKey    string `json:"appkey"`
}

type TokenResponse struct {
	Response
	Data struct {
		ExpireTime string `json:"expire_time"`
		Token      string `json:"token"`
	} `json:"data"`
}

// PushSingleRequest 【toSingle】执行cid单推
type PushSingleRequest struct {
	RequestId string `json:"request_id"`
	Settings  struct {
		Ttl int `json:"ttl"`
	} `json:"settings,omitempty"`
	Audience    PushSingleRequestAudience    `json:"audience"`
	PushMessage PushSingleRequestPushMessage `json:"push_message"`
}

type PushSingleRequestAudience struct {
	Cid []string `json:"cid"`
}

type PushSingleRequestPushMessage struct {
	Notification PushMessageNotification `json:"notification"`
}

type PushSingleResponse struct {
	Response
	Data struct {
		Taskid struct {
			Cid string `json:"$cid"`
		} `json:"$taskid"`
	} `json:"data"`
}

// CustomTagRequest 【标签】一个用户绑定一批标签
type CustomTagRequest struct {
	CustomTag []string `json:"custom_tag"`
}

// PushAllRequest 群推送
type PushAllRequest struct {
	RequestId string `json:"request_id"`
	GroupName string `json:"group_name"`
	Settings  struct {
		Ttl int `json:"ttl"`
	} `json:"settings,omitempty"`
	Audience    string `json:"audience"`
	PushMessage struct {
		PushMessageNotification `json:"notification"`
	} `json:"push_message"`
}

type PushMessageNotification struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	ClickType string `json:"click_type"`
	Url       string `json:"url"`
}

type PushAllResponse struct {
	Response
	Data struct {
		Taskid string `json:"taskid"`
	} `json:"data"`
}
