package waterproofwall

import "github.com/herb-go/deprecated/fetch"

type App struct {
	AppID        string
	AppSecretKey string
	Clients      fetch.Clients
}
