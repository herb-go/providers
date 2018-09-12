package vaptcha

import "github.com/herb-go/fetch"

type Cell struct {
	VID string
	Key string
}

type Config struct {
	*Cell
	Lang    string
	Type    string
	Clients fetch.Clients
}
