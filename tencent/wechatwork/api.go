package wechatwork

import (
	"github.com/herb-go/fetch"
)

var Server = fetch.Server{
	Host: "https://qyapi.weixin.qq.com",
}

var apiGetUserInfo = Server.EndPoint("GET", "/cgi-bin/user/getuserinfo")
var apiGetToken = Server.EndPoint("GET", "/cgi-bin/gettoken")
var apiGetUserDetail = Server.EndPoint("POST", "/cgi-bin/user/getuserdetail")
var apiUserGet = Server.EndPoint("GET", "/cgi-bin/user/get")
var apiMessagePost = Server.EndPoint("POST", "/cgi-bin/message/send")
var apiDepartmentList = Server.EndPoint("GET", "/cgi-bin/department/list")

const ApiErrAccessTokenWrong = 40014
const ApiErrAccessTokenOutOfDate = 42001
const ApiErrSuccess = 0
const ApiErrUserUnaccessible = 50002
const ApiErrOauthCodeWrong = 40029
const ApiErrNoPrivilege = 60011
const ApiResultGenderMale = "1"
const ApiResultGenderFemale = "2"

type bodyMessagePost struct {
	ToUser  string               `json:"touser"`
	ToParty string               `json:"toparty"`
	ToTag   string               `json:"totag"`
	MsgType string               `json:"msgtype"`
	AgentID int                  `json:"agentid"`
	Safe    int                  `json:"safe"`
	Text    *bodyMessagePostText `json:"text"`
}
type bodyMessagePostText struct {
	Content string `json:"content"`
}

type resultMessagePost struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

type resultAPIError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type resultAccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type resultUserInfo struct {
	UserID     string `json:"UserId"`
	UserTicket string `json:"user_ticket"`
}
type paramsUserDetail struct {
	UserTicket string `json:"user_ticket"`
}
type resultUserDetail struct {
	UserID   string `json:"userid"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type resultUserGet struct {
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	Mobile     string `json:"mobile"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Department []int  `json:"department"`
}
type DepartmentInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parentid"`
	Order    int    `json:"order"`
}
type resultDepartmentList struct {
	Department *[]DepartmentInfo `json:"department"`
}
