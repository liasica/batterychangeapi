package getui

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gtime"
	"golang.org/x/sync/singleflight"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	CacheKeyGeTuiToken = "CACHE:GETUI:TOKEN"

	UriToken = "/auth"
)

type service struct {
	appId           string
	appKey          string
	appSecret       string
	appMasterSecret string
	baseUrl         string
}

var serv *service

func Service() *service {
	if serv == nil {
		serv = &service{
			appId:           g.Cfg().GetString("getui.appId"),
			appKey:          g.Cfg().GetString("getui.appKey"),
			appSecret:       g.Cfg().GetString("getui.appSecret"),
			appMasterSecret: g.Cfg().GetString("getui.appMasterSecret"),
			baseUrl:         fmt.Sprintf("https://restapi.getui.com/v2/%s", g.Cfg().GetString("getui.appId")),
		}
	}
	return serv
}

func (s *service) rpc(method, uri string, data interface{}) ([]byte, error) {
	client := g.Client().SetHeader("Content-Type", "application/json")
	if uri != UriToken {
		token, err := s.Token()
		if err != nil {
			return nil, err
		}
		client = client.SetHeader("token", token)
	}
	resp, err := client.DoRequest(method, fmt.Sprintf("%s%s", s.baseUrl, uri), data)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return ioutil.ReadAll(resp.Body)
}

// Token 获取请求token
func (s *service) Token() (string, error) {
	if val, err := gcache.Get(CacheKeyGeTuiToken); err == nil && val != nil {
		return val.(string), err
	}
	var singleSetCache singleflight.Group
	token, err, _ := singleSetCache.Do(CacheKeyGeTuiToken, func() (interface{}, error) {
		t := strconv.FormatInt(gtime.TimestampMilli(), 10)
		h := sha256.New()
		_, err := io.WriteString(h, s.appKey+t+s.appMasterSecret)
		if err != nil {
			return "", err
		}
		content, err := s.rpc(http.MethodPost, UriToken, TokenRequest{
			Sign:      fmt.Sprintf("%x", h.Sum(nil)),
			Timestamp: t,
			AppKey:    s.appKey,
		})
		if err != nil {
			return "", err
		}
		var res TokenResponse
		if _err := json.Unmarshal(content, &res); _err != nil {
			return "", _err
		}
		if res.Code != 0 {
			return "", errors.New(res.Msg)
		}
		gcache.Set(CacheKeyGeTuiToken, res.Data.Token, 12*time.Hour)
		return res.Data.Token, err
	})

	return token.(string), err
}

// PushAll 群推
func (s *service) PushAll(msg PushAllRequest) (res PushAllResponse, err error) {
	content, err := s.rpc(http.MethodPost, "/push/all", msg)
	if err != nil {
		return res, err
	}
	_err := json.Unmarshal(content, &res)
	return res, _err
}

// PushSingleCid 单推
func (s *service) PushSingleCid(msg PushSingleRequest) (res PushSingleResponse, err error) {
	content, err := s.rpc(http.MethodPost, "/push/single/cid", msg)
	if err != nil {
		return res, err
	}
	_err := json.Unmarshal(content, &res)
	return res, _err
}
