package facebook

import (
	"net/http"

	"github.com/herb-go/fetch"
)

var Server = fetch.Server{
	Host:    "https://graph.facebook.com",
	Headers: http.Header{},
}

var apiAccessToken = Server.EndPoint("POST", "/v3.2/oauth/access_token")
var apiUser = Server.EndPoint("GET", "/me")

type ResultAPIAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type ResultAPIUser struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Name       string `json:"name"`
	ProfilePic string `json:"profile_pic"`
}
