package wechatwork

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/herb-go/fetch"
)

type Agent struct {
	CorpID      string
	AgentID     string
	Secret      string
	Clients     fetch.Clients
	accessToken string
	lock        sync.Mutex
}

func (a *Agent) AccessToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.accessToken
}
func (a *Agent) sendMessage(b *bodyMessagePost) (*resultMessagePost, error) {
	result := &resultMessagePost{}
	err := a.CallApiWithAccessToken(apiMessagePost, nil, b, result)
	return result, err
}
func (a *Agent) SendTextMessageToUsers(users []string, content string) (*resultMessagePost, error) {
	var err error
	message := &bodyMessagePost{}
	message.ToUser = strings.Join(users, "|")
	message.AgentID, err = strconv.Atoi(a.AgentID)
	message.MsgType = "text"
	message.Text = &bodyMessagePostText{
		Content: content,
	}
	if err != nil {
		return nil, err
	}
	return a.sendMessage(message)
}

func (a *Agent) GrantAccessToken() error {
	params := url.Values{}
	params.Set("corpid", a.CorpID)
	params.Set("corpsecret", a.Secret)
	req, err := apiGetToken.NewRequest(params, nil)
	if err != nil {
		return err
	}
	rep, err := a.Clients.Fetch(req)
	if err != nil {
		return err
	}
	if rep.StatusCode != http.StatusOK {
		return rep
	}
	result := &resultAccessToken{}
	err = rep.UnmarshalAsJSON(result)
	if err != nil {
		return err
	}
	if result.Errcode != 0 || result.Errmsg == "" || result.AccessToken == "" {
		return rep.NewAPICodeErr(result.Errcode)
	}
	a.accessToken = result.AccessToken
	return nil
}

func (a *Agent) CallApiWithAccessToken(api *fetch.EndPoint, params url.Values, body interface{}, v interface{}) error {
	var apierr resultAPIError
	var err error
	if a.AccessToken() == "" {
		err := a.GrantAccessToken()
		if err != nil {
			return err
		}
	}
	p := url.Values{}
	if params != nil {
		for k, vs := range params {
			for _, v := range vs {
				p.Add(k, v)
			}
		}
	}
	p.Set("access_token", a.AccessToken())
	req, err := api.NewJSONRequest(p, body)
	if err != nil {
		return err
	}
	resp, err := a.Clients.Fetch(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return resp
	}
	apierr = resultAPIError{}
	err = resp.UnmarshalAsJSON(&apierr)
	if err != nil {
		return err
	}
	if fetch.CompareAPIErrCode(err, ApiErrAccessTokenOutOfDate) || fetch.CompareAPIErrCode(err, ApiErrAccessTokenWrong) {
		err := a.GrantAccessToken()
		if err != nil {
			return err
		}
		p.Set("access_token", a.AccessToken())
		req, err := api.NewJSONRequest(p, body)
		if err != nil {
			return err
		}
		resp, err := a.Clients.Fetch(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return resp
		}
		apierr = resultAPIError{}
		err = resp.UnmarshalAsJSON(&apierr)
		if err != nil {
			return err
		}
	}
	if apierr.Errcode != 0 {
		return resp.NewAPICodeErr(apierr.Errcode)
	}
	return resp.UnmarshalAsJSON(&v)
}

type Userinfo struct {
	UserID     string
	Name       string
	Mobile     string
	Email      string
	Gender     string
	Avatar     string
	Department []int
}

func (a *Agent) GetUserInfo(code string) (*Userinfo, error) {
	var info = &Userinfo{}
	if code == "" {
		return nil, nil
	}
	var result = &resultUserInfo{}
	params := url.Values{}
	params.Set("code", code)
	err := a.CallApiWithAccessToken(apiGetUserInfo, params, nil, result)
	if err != nil {
		return nil, err
	}
	if result.UserID == "" {
		return nil, nil
	}
	var getuser = &resultUserGet{}
	userGetParam := url.Values{}
	userGetParam.Add("userid", result.UserID)
	err = a.CallApiWithAccessToken(apiUserGet, userGetParam, nil, getuser)
	if err != nil {
		if fetch.CompareAPIErrCode(err, ApiErrUserUnaccessible) || fetch.CompareAPIErrCode(err, ApiErrNoPrivilege) {
			return nil, nil
		}
		return nil, err
	}
	info.UserID = result.UserID
	info.Avatar = getuser.Avatar
	info.Email = getuser.Email
	info.Gender = getuser.Gender
	info.Mobile = getuser.Mobile
	info.Name = getuser.Name
	info.Department = getuser.Department
	return info, nil
}

func (a *Agent) GetDepartmentList(id string) (*[]DepartmentInfo, error) {
	params := url.Values{}
	if id != "" {
		params.Set("id", id)
	}
	var result = &resultDepartmentList{}
	err := a.CallApiWithAccessToken(apiDepartmentList, params, nil, result)
	if err != nil {
		return nil, err
	}
	return result.Department, nil
}
