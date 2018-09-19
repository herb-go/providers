package waterproofwall

import "github.com/herb-go/fetch"

var Server = fetch.Server{
	Host: "https://ssl.captcha.qq.com",
}

var ApiValidate = Server.EndPoint("GET", "/ticket/verify")

type ResultValidate struct {
	Response string `json:"response"`
	EviLevel string `json:"evil_level"`
	ErrMsg   string `json:"err_msg"`
}
