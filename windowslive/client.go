package windowslive

import (
	"net/http"
	"net/url"

	"github.com/herb-go/fetch"
)

type Client struct {
	ClientID     string
	ClientSecret string
	Clients      fetch.Clients
}

func (c *Client) GetAccessToken(code string, redirect_url string) (*ResultAPIAccessToken, error) {
	params := url.Values{}
	params.Set("client_id", c.ClientID)
	params.Set("client_secret", c.ClientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", redirect_url)
	params.Set("grant_type", "authorization_code")
	req, err := apiAccessToken.NewRequest(nil, []byte(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	rep, err := c.Clients.Fetch(req)
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
	return result, nil
}
func (c *Client) GetUser(accessToken string) (*ResultAPIUser, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)

	req, err := apiUser.NewRequest(params, nil)
	if err != nil {
		return nil, err
	}
	rep, err := c.Clients.Fetch(req)
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
	return result, nil
}
