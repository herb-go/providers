package wechatmp

import (
	"net/url"
	"sync"

	"github.com/herb-go/fetcher"
)

type App struct {
	AppID              string
	AppSecret          string
	Client             fetcher.Client
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

func (a *App) ClientCredentialBuilder() fetcher.Command {
	return fetcher.CommandFunc(func(f *fetcher.Fetcher) error {
		params := f.URL.Query()
		params.Set("appid", a.AppID)
		params.Set("secret", a.AppSecret)
		params.Set("grant_type", "client_credential")
		return nil
	})
}
func (a *App) AuthorizationCodeBuilder(code string) fetcher.Command {
	return fetcher.CommandFunc(func(f *fetcher.Fetcher) error {
		params := f.URL.Query()
		params.Set("appid", a.AppID)
		params.Set("secret", a.AppSecret)
		params.Set("grant_type", "authorization_code")
		params.Set("code", code)
		return nil
	})
}

func (a *App) GetAccessToken() (string, error) {
	result := &resultAccessToken{}
	resp, err := fetcher.DoAndParse(
		&a.Client,
		APIToken.With(a.ClientCredentialBuilder()),
		fetcher.Should200(fetcher.AsJSON(result)),
	)
	if err != nil {
		return "", err
	}
	if result.Errcode != 0 || result.Errmsg != "" || result.AccessToken == "" {
		return "", resp.NewAPICodeErr(result.Errcode)
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

func (a *App) callApiWithAccessToken(api *fetcher.Preset, APIPresetBuilder func(accesstoken string) (*fetcher.Preset, error), v interface{}) error {
	var apierr = &ResultAPIError{}
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
	preset, err := APIPresetBuilder(token)
	if err != nil {
		return err
	}
	resp, err := fetcher.DoAndParse(&a.Client, preset, fetcher.Should200(fetcher.AsJSON(apierr)))
	if err != nil {
		return err
	}
	if !apierr.IsOK() {
		if apierr.IsAccessTokenError() {
			token, err = a.GrantAccessToken()
			if err != nil {
				return err
			}
			apierr = &ResultAPIError{}
			resp, err = fetcher.DoAndParse(&a.Client, preset, fetcher.Should200(fetcher.AsJSON(apierr)))
			if err != nil {
				return err
			}
			if !apierr.IsOK() {
				return resp.NewAPICodeErr(apierr.Errcode)
			}
		} else {
			return resp
		}
	}
	return fetcher.AsJSON(v).Parse(resp)
}

func (a *App) CallJSONApiWithAccessToken(api *fetcher.Preset, params url.Values, body interface{}, v interface{}) error {
	jsonAPIPresetBuilder := func(accesstoken string) (*fetcher.Preset, error) {
		return api.With(fetcher.SetQuery("access_token", accesstoken), fetcher.JSONBody(body)), nil
	}
	return a.callApiWithAccessToken(api, jsonAPIPresetBuilder, v)
}

func (a *App) GetUserInfo(code string, scope string, lang string) (*Userinfo, error) {
	var info = &Userinfo{}
	if code == "" {
		return nil, nil
	}
	var result = &resultOauthToken{}
	resp, err := fetcher.DoAndParse(
		&a.Client,
		APIOauth2AccessToken.With(a.AuthorizationCodeBuilder(code)),
		fetcher.Should200(fetcher.AsJSON(result)),
	)
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
	resp, err = fetcher.DoAndParse(
		&a.Client,
		APIGetUserInfo.With(
			fetcher.SetQuery("access_token", result.AccessToken),
			fetcher.SetQuery("openid", result.OpenID),
			fetcher.SetQuery("lang", lang),
		),
		fetcher.Should200(fetcher.AsJSON(getuser)),
	)
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
