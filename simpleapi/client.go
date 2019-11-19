package simpleapi

import (
	"github.com/herb-go/herb/service/httpservice/fetcher"
)

type Client struct {
	EndPoint fetcher.TargetGetter
	Clients  fetcher.Client
}
