package simpleapi

import (
	"net/url"

	"github.com/herb-go/fetch"
)

type Client struct {
	EndPoint *fetch.EndPoint
	Clients  *fetch.Clients
}

func (c *Client) FetchRequest(params url.Values, body []byte) (resp *fetch.Result, err error) {
	return c.Clients.FetchWithError(c.EndPoint.NewRequest(params, body))
}

func (c *Client) FetchJSONRequest(params url.Values, v interface{}) (resp *fetch.Result, err error) {
	return c.Clients.FetchWithError(c.EndPoint.NewJSONRequest(params, v))
}

func (c *Client) FetchXMLRequest(params url.Values, v interface{}) (resp *fetch.Result, err error) {
	return c.Clients.FetchWithError(c.EndPoint.NewXMLRequest(params, v))
}
