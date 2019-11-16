package httpapiserver

import (
	"net/url"

	"github.com/herb-go/fetch"
)

type Client struct {
	EndPoint *fetch.EndPoint
	Builder  *fetch.RequestBuilder
	Clients  fetch.Clients
}

func (c *Client) FetchRequest(params url.Values, body []byte) (*fetch.Result, error) {
	req, err := c.Builder.Apply(c.EndPoint.NewRequest(params, body))
	if err != nil {
		return nil, err
	}
	return c.Clients.Fetch(req)
}

func (c *Client) FetchJSONRequest(params url.Values, v interface{}) (*fetch.Result, error) {
	req, err := c.Builder.Apply(c.EndPoint.NewJSONRequest(params, v))
	if err != nil {
		return nil, err
	}
	return c.Clients.Fetch(req)
}

func (c *Client) FetchXMLRequest(params url.Values, v interface{}) (*fetch.Result, error) {
	req, err := c.Builder.Apply(c.EndPoint.NewXMLRequest(params, v))
	if err != nil {
		return nil, err
	}
	return c.Clients.Fetch(req)
}
