package facebook

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/herb-go/deprecated/fetch"
)

type App struct {
	ID      string
	Key     string
	Clients fetch.Clients
}

var DefaultFields = []string{
	"id",
	"email",
	"first_name",
	"middle_name",
	"last_name",
	"gender",
	"name",
	"profile_pic",
}

func (a *App) GetAccessToken(code string, redirect_url string) (*ResultAPIAccessToken, error) {
	params := url.Values{}
	params.Set("client_id", a.ID)
	params.Set("client_secret", a.Key)
	params.Set("code", code)
	params.Set("redirect_uri", redirect_url)
	req, err := apiAccessToken.NewRequest(nil, []byte(params.Encode()))
	if err != nil {
		return nil, err
	}
	// req.Header.Set("Accept", "application/json")
	rep, err := a.Clients.Fetch(req)
	if err != nil {
		return nil, err
	}
	if rep.StatusCode != http.StatusOK {
		return nil, rep
	}
	result := &ResultAPIAccessToken{}
	err = rep.UnmarshalAsJSON(result)
	if err != nil {
		return nil, err
	}
	if result.AccessToken == "" {
		return nil, rep
	}
	return result, nil
}
func (a *App) GetUser(accessToken string) (*ResultAPIUser, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("field", strings.Join(DefaultFields, ","))
	req, err := apiUser.NewRequest(params, nil)
	if err != nil {
		return nil, err
	}
	rep, err := a.Clients.Fetch(req)
	if err != nil {
		return nil, err
	}
	if rep.StatusCode != http.StatusOK {
		return nil, rep
	}
	result := &ResultAPIUser{}
	err = rep.UnmarshalAsJSON(result)
	if err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, rep
	}
	return result, nil
}
