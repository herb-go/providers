package waterproofwall

import "github.com/herb-go/fetch"

type App struct {
	AppID        string
	AppSecretKey string
	Clients      fetch.Clients
}
