package hiredmember

import (
	"github.com/herb-go/member"
)

type Directive struct {
	Name   string
	Config `config:", lazyload"`
}
type Config struct {
	Directives []*Directive
}

func (c *Config) ApplyTo(s *member.Service) error {
	for k := range c.Directives {

	}
}
