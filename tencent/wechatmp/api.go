package wechatmp

import (
	"encoding/json"

	"github.com/herb-go/fetch"
)

const ScopeSnsapiBase = "snsapi_base"
const ScopeSnsapiUserinfo = "snsapi_userinfo"

var Server = fetch.Server{
	Host: "https://api.weixin.qq.com",
}

var apiGetUserInfo = Server.EndPoint("GET", "/sns/userinfo")
var apiToken = Server.EndPoint("GET", "/cgi-bin/token")
var apiOauth2AccessToken = Server.EndPoint("GET", "/sns/oauth2/access_token")

var APIMenuCreate = Server.EndPoint("POST", "/cgi-bin/menu/create")

var APIMenuGet = Server.EndPoint("GET", "/cgi-bin/menu/get")

var APIQRCodeCreate = Server.EndPoint("POST", "/cgi-bin/qrcode/create")

var APIGetAllPrivateTemplate = Server.EndPoint("GET", "/cgi-bin/template/get_all_private_template?")

const ApiErrAccessTokenNotLast = 40001
const ApiErrAccessTokenWrong = 40014
const ApiErrAccessTokenOutOfDate = 42001
const ApiErrSuccess = 0
const ApiErrUserUnaccessible = 50002
const ApiErrOauthCodeWrong = 40029

const ApiResultGenderMale = 1
const ApiResultGenderFemale = 2

type ResultAPIError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type resultAccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type resultOauthToken struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}
type resultUserDetail struct {
	Errcode    int             `json:"errcode"`
	Errmsg     string          `json:"errmsg"`
	OpenID     string          `json:"openid"`
	Nickname   string          `json:"nickname"`
	Sex        int             `json:"sex"`
	Province   string          `json:"province"`
	City       string          `json:"city"`
	Country    string          `json:"country"`
	HeadimgURL string          `json:"headimgurl"`
	Privilege  json.RawMessage `json:"privilege"`
	UnionID    string          `json:"unionid"`
}

type ResultQRCodeCreate struct {
	Errcode       int    `json:"errcode"`
	Errmsg        string `json:"errmsg"`
	Ticket        string `json:"ticket"`
	ExpireSeconds *int64 `json:"expire_seconds"`
	URL           string `json:"url"`
}

type PrivateTemplate struct {
	TemplateID      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

type AllPrivateTemplateResult struct {
	Errcode      int               `json:"errcode"`
	Errmsg       string            `json:"errmsg"`
	TemplateList []PrivateTemplate `json:"template_list"`
}
