package httpapiserver

import (
	"github.com/herb-go/fetch"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type Client struct {
	fetch.Server
	Channel string
	Clients fetch.Clients
	Vistor  guarder.DirverConfigMap
}

func (c *Client) Builder() (*fetch.RequestBuilder, error) {
	var err error
	v := guarder.NewVisitor()
	err = c.Vistor.ApplyToVisitor(v)
	if err != nil {
		return nil, err
	}
	return fetch.NewRequestBuilder(v), nil
}
