package tencentminiprogram

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

func (a *App) Login(code string) (*ResultUserInfo, error) {
	params := url.Values{}
	params.Set("appid", a.AppID)
	params.Set("secret", a.AppSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")
	preset := apiLogin.CloneWith(
		&a.Client,
		fetcher.Params(params),
	)
	result := &ResultUserInfo{}
	resp, err := preset.FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, err
	}
	if result.Errcode != 0 {
		return nil, resp.NewAPICodeErr(result.Errcode)
	}
	return result, nil
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
	return fetcher.ParamsBuilderFunc(func(params url.Values) error {
		params.Set("appid", a.AppID)
		params.Set("secret", a.AppSecret)
		params.Set("grant_type", "client_credential")
		return nil

	})
}

func (a *App) GetAccessToken() (string, error) {
	result := &resultAccessToken{}
	resp, err := APIToken.CloneWith(
		&a.Client,
		a.ClientCredentialBuilder()).
		FetchAndParse(fetcher.AsJSON(result))
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
