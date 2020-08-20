package tencentminiprogram

import "github.com/herb-go/deprecated/fetch"

var Server = fetch.Server{
	Host: "https://api.weixin.qq.com",
}

var apiLogin = Server.EndPoint("GET", "/sns/jscode2session")

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
