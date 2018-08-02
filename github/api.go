package github

import (
	"net/http"

	"github.com/herb-go/fetch"
)

var Server = fetch.Server{
	Host: "https://github.com",
	Headers: http.Header{
		"Accept": []string{"application/json"},
	},
}
var APIServer = fetch.Server{
	Host: "https://api.github.com",
	Headers: http.Header{
		"Accept": []string{"application/json"},
	},
}

var apiAccessToken = Server.EndPoint("POST", "/login/oauth/access_token")
var apiUser = APIServer.EndPoint("GET", "/user")

type ResultAPIAccessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type ResultAPIUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Blog      string `json:"blog"`
	Location  string `json:"location"`
	Email     string `json:"email"`
}
