package httpapiserver

import (
	"github.com/herb-go/fetch"
)

type Client struct {
	EndPoint *fetch.EndPoint
	Builder  *fetch.RequestBuilder
	Clients  fetch.Clients
}
