package simpleapi

import (
	"github.com/herb-go/herb/service/httpservice/fetcher"
)

type Client struct {
	Target fetcher.TargetGetter
	Client fetcher.Client
}
