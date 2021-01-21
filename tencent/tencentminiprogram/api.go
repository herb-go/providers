package tencentminiprogram

import (
	"github.com/herb-go/fetcher"
)

var Server = fetcher.Preset{
	fetcher.URL("https://api.weixin.qq.com"),
}

var apiLogin = Server.EndPoint("GET", "/sns/jscode2session")
var APIToken = Server.EndPoint("GET", "/cgi-bin/token")

type ResultAPIError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type ResultUserInfo struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

type resultAccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
