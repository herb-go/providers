package simpleapi

import "github.com/herb-go/herb/service/httpservice/target"

type Client struct {
	Target target.Target
	Doer   target.Doer
}
