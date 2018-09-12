package vaptcha

import "github.com/herb-go/fetch"

var Server = fetch.Server{
	Host: "https://api.vaptcha.com/v2/validate",
}

var ApiValidate = Server.EndPoint("POST", "/v2/validate")

var MsgTokenExpired = "token-error"
var MsgTokenError = "token-error"

type ResultValidate struct {
	Success int    `json:"success"`
	Score   int    `json:"score"`
	Msg     string `json:"msg"`
}
