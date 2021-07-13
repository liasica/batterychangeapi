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
