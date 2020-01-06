package wechatwork

import (
	"errors"
	"strconv"

	"github.com/herb-go/fetcher"
)

var Server = fetcher.MustPreset(&fetcher.ServerInfo{
	URL: "https://qyapi.weixin.qq.com",
})

var apiGetUserInfo = Server.EndPoint("GET", "/cgi-bin/user/getuserinfo")
var apiGetToken = Server.EndPoint("GET", "/cgi-bin/gettoken")
var apiGetUserDetail = Server.EndPoint("POST", "/cgi-bin/user/getuserdetail")
var apiUserGet = Server.EndPoint("GET", "/cgi-bin/user/get")
var apiMessagePost = Server.EndPoint("POST", "/cgi-bin/message/send")
var apiDepartmentList = Server.EndPoint("GET", "/cgi-bin/department/list")
var apiMediaUpload = Server.EndPoint("POST", "/cgi-bin/media/upload")

const APIErrAccessTokenNotLast = 40001
const APIErrAccessTokenWrong = 40014
const APIErrAccessTokenOutOfDate = 42001
const APIErrSuccess = 0
const APIErrUserUnaccessible = 50002
const APIErrOauthCodeWrong = 40029
const APIErrNoPrivilege = 60011
const APIResultGenderMale = "1"
const APIResultGenderFemale = "2"

type resultAPIError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (e *resultAPIError) IsOK() bool {
	return e.Errcode == 0
}

func (e *resultAPIError) IsAccessTokenError() bool {
	return e.Errcode == APIErrAccessTokenOutOfDate || e.Errcode == APIErrAccessTokenWrong || e.Errcode == APIErrAccessTokenNotLast
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

type resultMediaUpload struct {
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
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

func NewResultError(code int, msg string) error {
	return errors.New("wechat work resuld error :" + strconv.Itoa(code) + " - " + msg)
}
