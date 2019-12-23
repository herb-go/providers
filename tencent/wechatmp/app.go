package wechatmp

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"

	"github.com/herb-go/fetch"
)

type App struct {
	AppID              string
	AppSecret          string
	Clients            fetch.Clients
	accessToken        string
	lock               sync.Mutex
	accessTokenGetter  func() (string, error)
	accessTokenCreator func() (string, error)
}

func (a *App) SetAccessTokenGetter(f func() (string, error)) {
	a.accessTokenGetter = f
}
func (a *App) SetAccessTokenCreator(f func() (string, error)) {
	a.accessTokenCreator = f
}
func (a *App) AccessToken() (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.accessTokenGetter != nil {
		return a.accessTokenGetter()
	}
	return a.accessToken, nil
}
func (a *App) GetAccessToken() (string, error) {
	params := url.Values{}
	params.Set("appid", a.AppID)
	params.Set("secret", a.AppSecret)
	params.Set("grant_type", "client_credential")
	req, err := apiToken.NewRequest(params, nil)
	if err != nil {
		return "", err
	}
	rep, err := a.Clients.Fetch(req)
	if err != nil {
		return "", err
	}
	if rep.StatusCode != http.StatusOK {
		return "", rep
	}
	result := &resultAccessToken{}
	err = rep.UnmarshalAsJSON(result)
	if err != nil {
		return "", err
	}
	if result.Errcode != 0 || result.Errmsg != "" || result.AccessToken == "" {
		return "", rep.NewAPICodeErr(result.Errcode)
	}
	return result.AccessToken, nil
}
func (a *App) GrantAccessToken() (string, error) {
	var token string
	var err error
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.accessTokenCreator == nil {
		token, err = a.GetAccessToken()
	} else {
		token, err = a.accessTokenCreator()
	}

	if err != nil {
		return "", err
	}
	a.accessToken = token
	return token, nil
}

func (a *App) callApiWithAccessToken(api *fetch.EndPoint, APIRequestBuilder func(accesstoken string) (*http.Request, error), v interface{}) error {
	var apierr ResultAPIError
	var err error
	token, err := a.AccessToken()
	if err != nil {
		return err
	}
	if token == "" {
		token, err = a.GrantAccessToken()
		if err != nil {
			return err
		}
	}

	req, err := APIRequestBuilder(token)
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
	if err != nil {
		return err
	}
	apierr = ResultAPIError{}
	err = resp.UnmarshalAsJSON(&apierr)
	if err != nil {
		return err
	}
	if apierr.Errcode != 0 {
		if apierr.Errcode == ApiErrAccessTokenOutOfDate || apierr.Errcode == ApiErrAccessTokenWrong || apierr.Errcode == ApiErrAccessTokenNotLast {
			token, err = a.GrantAccessToken()
			if err != nil {
				return err
			}
			req, err := APIRequestBuilder(token)
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
			apierr = ResultAPIError{}
			err = resp.UnmarshalAsJSON(&apierr)
			if err != nil {
				return err
			}
			if apierr.Errcode != 0 {
				return resp.NewAPICodeErr(apierr.Errcode)
			}
			return nil
		}
		return resp
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
