package httpapiserver

import (
	"github.com/herb-go/fetch"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type Client struct {
	fetch.Server
	Channel string
	Clients fetch.Clients
	Vistor  *guarder.Visitor
}

func (c *Client) Builder() (*fetch.RequestBuilder, error) {

	return fetch.NewRequestBuilder(c.Vistor), nil
}
