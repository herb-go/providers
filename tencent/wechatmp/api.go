package wechatmp

import (
	"encoding/json"

	"github.com/herb-go/fetcher"
)

const ScopeSnsapiBase = "snsapi_base"
const ScopeSnsapiUserinfo = "snsapi_userinfo"

var Server = fetcher.MustPreset(&fetcher.ServerInfo{
	URL: "https://api.weixin.qq.com",
})

var APIGetUserInfo = Server.EndPoint("GET", "/sns/userinfo")
var APIToken = Server.EndPoint("GET", "/cgi-bin/token")
var APIOauth2AccessToken = Server.EndPoint("GET", "/sns/oauth2/access_token")

var APIMenuCreate = Server.EndPoint("POST", "/cgi-bin/menu/create")

var APIMenuGet = Server.EndPoint("GET", "/cgi-bin/menu/get")

var APIQRCodeCreate = Server.EndPoint("POST", "/cgi-bin/qrcode/create")

var APIGetAllPrivateTemplate = Server.EndPoint("GET", "/cgi-bin/template/get_all_private_template")

var APIMessageTemplateSend = Server.EndPoint("POST", "/cgi-bin/message/template/send?")

const APIErrAccessTokenNotLast = 40001
const APIErrAccessTokenWrong = 40014
const APIErrAccessTokenOutOfDate = 42001
const APIErrSuccess = 0
const APIErrUserUnaccessible = 50002
const APIErrOauthCodeWrong = 40029

const APIResultGenderMale = 1
const APIResultGenderFemale = 2

type ResultAPIError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (e *ResultAPIError) IsOK() bool {
	return e.Errcode == 0
}

func (e *ResultAPIError) IsAccessTokenError() bool {
	return e.Errcode == APIErrAccessTokenOutOfDate || e.Errcode == APIErrAccessTokenWrong || e.Errcode == APIErrAccessTokenNotLast
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

type TemplateMessageMiniprogram struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
}
type TemplateMessage struct {
	ToUser      string                      `json:"touser"`
	TemlpateID  string                      `json:"template_id"`
	Miniprogram *TemplateMessageMiniprogram `json:"miniprogram"`
	URL         *string                     `json:"url"`
	Data        json.RawMessage             `json:"data"`
}

func NewTemplateMessage() *TemplateMessage {
	return &TemplateMessage{}
}

type TemplateMessageSendResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	MsgID   int64  `json:"msgid"`
}

type Userinfo struct {
	OpenID       string
	Nickname     string
	Sex          int
	Province     string
	City         string
	Country      string
	HeadimgURL   string
	Privilege    json.RawMessage
	UnionID      string
	AccessToken  string
	RefreshToken string
}

type resultUserInfo struct {
	UserID     string `json:"UserId"`
	UserTicket string `json:"user_ticket"`
}
