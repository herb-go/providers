package windowslive

import (
	"net/http"

	"github.com/herb-go/fetch"
)

var Server = fetch.Server{
	Host: "https://login.live.com",
	Headers: http.Header{
		"Content-type": []string{"application/x-www-form-urlencoded"},
	},
}
var APIServer = fetch.Server{
	Host: "https://apis.live.net",
}

var apiAccessToken = Server.EndPoint("POST", "/oauth20_token.srf")
var apiUser = APIServer.EndPoint("GET", "/v5.0/me")

type ResultAPIAccessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type ResultAPIUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Locale    string `json:"locale"`
}
