package tencentcloud

import "github.com/herb-go/deprecated/fetch"

type App struct {
	AppID     string
	SecretID  string
	SecretKey string
	Clients   fetch.Clients
}
