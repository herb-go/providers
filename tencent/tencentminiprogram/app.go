package tencentminiprogram

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/herb-go/deprecated/fetch"
)

type APP struct {
	AppID     string
	AppSecret string
	Clients   fetch.Clients
}

func (a *APP) Login(code string) (*ResultUserInfo, error) {
	params := url.Values{}
	params.Set("appid", a.AppID)
	params.Set("secret", a.AppSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")
	req, err := apiLogin.NewRequest(params, nil)
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
	result := &ResultUserInfo{}
	err = json.Unmarshal(rep.BodyContent, result)
	if err != nil {
		return nil, err
	}
	if result.Errcode != 0 {
		return nil, rep.NewAPICodeErr(result.Errcode)
	}
	return result, nil
}
