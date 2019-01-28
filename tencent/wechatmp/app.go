package wechatmp

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"

	"github.com/herb-go/fetch"
)

type App struct {
	AppID       string
	AppSecret   string
	Clients     fetch.Clients
	accessToken string
	lock        sync.Mutex
}

func (a *App) AccessToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.accessToken
}

func (a *App) GrantAccessToken() error {
	params := url.Values{}
	params.Set("appid", a.AppID)
	params.Set("secret", a.AppSecret)
	params.Set("grant_type", "client_credential")
	req, err := apiToken.NewRequest(params, nil)
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
	a.lock.Lock()
	defer a.lock.Unlock()
	a.accessToken = result.AccessToken
	return nil
}

func (a *App) callApiWithAccessToken(api *fetch.EndPoint, APIRequestBuilder func(accesstoken string) (*http.Request, error), v interface{}) error {
	var apierr resultAPIError
	var err error
	if a.AccessToken() == "" {
		err := a.GrantAccessToken()
		if err != nil {
			return err
		}
	}

	req, err := APIRequestBuilder(a.AccessToken())
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
		req, err := APIRequestBuilder(a.AccessToken())
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

func (a *App) CallJSONApiWithAccessToken(api *fetch.EndPoint, params url.Values, body interface{}, v interface{}) error {
	jsonAPIRequestBuilder := func(accesstoken string) (*http.Request, error) {
		p := url.Values{}
		if params != nil {
			for k, vs := range params {
				for _, v := range vs {
					p.Add(k, v)
				}
			}
		}
		p.Set("access_token", accesstoken)
		return api.NewJSONRequest(p, body)
	}
	return a.callApiWithAccessToken(api, jsonAPIRequestBuilder, v)
}

func (a *App) GetUserInfo(code string, scope string, lang string) (*Userinfo, error) {
	var info = &Userinfo{}
	if code == "" {
		return nil, nil
	}
	var result = &resultOauthToken{}
	params := url.Values{}
	params.Set("appid", a.AppID)
	params.Set("secret", a.AppSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	req, err := apiOauth2AccessToken.NewJSONRequest(params, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.Clients.Fetch(req)
	if err != nil {
		return nil, err
	}
	err = resp.UnmarshalAsJSON(result)
	if err != nil {
		return nil, err
	}
	if result.AccessToken == "" {
		return nil, resp
	}
	info.OpenID = result.OpenID
	info.AccessToken = result.AccessToken
	info.RefreshToken = result.RefreshToken
	info.UnionID = result.UnionID
	if scope != ScopeSnsapiUserinfo {
		return info, nil
	}
	var getuser = &resultUserDetail{}
	userGetParam := url.Values{}
	userGetParam.Add("access_token", result.AccessToken)
	userGetParam.Add("openid", result.OpenID)
	userGetParam.Add("lang", lang)
	req, err = apiGetUserInfo.NewJSONRequest(userGetParam, nil)
	if err != nil {
		return nil, err
	}
	resp, err = a.Clients.Fetch(req)
	if err != nil {
		return nil, err
	}
	err = resp.UnmarshalAsJSON(getuser)
	if err != nil {
		return nil, err
	}
	if getuser.Errcode != 0 {
		return nil, resp.NewAPICodeErr(getuser.Errcode)
	}

	info.Nickname = getuser.Nickname
	info.Sex = getuser.Sex
	info.Province = getuser.Province
	info.City = getuser.City
	info.Country = getuser.Country
	info.HeadimgURL = getuser.HeadimgURL
	info.Privilege = getuser.Privilege
	info.UnionID = getuser.UnionID
	return info, nil
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
