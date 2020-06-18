package hiredmember

import (
	"github.com/herb-go/member"
	"github.com/herb-go/providers/herb/overseers/memberdirectivefactoryoverseer"
)

type Directive struct {
	Name   string
	Config func(v interface{}) error `config:", lazyload"`
}

func (d *Directive) ApplyTo(s *member.Service) error {
	f := memberdirectivefactoryoverseer.GetMemberDirectiveFactoryByID(d.Name)
	directive, err := f(d.Config)
	if err != nil {
		return err
	}
	return directive.ApplyTo(s)
}

type Config struct {
	Directives []*Directive
}

func (c *Config) ApplyTo(s *member.Service) error {
	for k := range c.Directives {
		err := c.Directives[k].ApplyTo(s)
		if err != nil {
			return err
		}
	}
	return nil
}
