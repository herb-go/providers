package tencentcloud

import (
	"github.com/herb-go/fetcher"
)

type App struct {
	AppID     string
	SecretID  string
	SecretKey string
	Client    fetcher.Client
}
