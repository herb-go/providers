package tomlmember

import (
	"github.com/herb-go/member"
	"github.com/herb-go/providers/herb/statictoml"
)

type Config struct {
	Source             statictoml.Source
	AsPasswordProvider bool
	AsStatusProvider   bool
	AsAccountsProvider bool
	AsRoleProvider     bool
}

func (c *Config) Load() (*Users, error) {
	u := newUsers()
	u.Source = c.Source
	data := []*User{}
	err := u.Source.Load(data)
	if err != nil {
		return nil, err
	}
	for k := range data {
		u.addUser(data[k])
	}
	return u, nil
}

func RegisterMemberDirective() {
	member.Register("statictoml", func(loader func(v interface{}) error) (member.Directive, error) {
		c := &Config{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		return member.DirectiveFunc(func(m *member.Service) error {
			u, err := c.Load()
			if err != nil {
				return err
			}
			if c.AsAccountsProvider {
				m.AccountsProvider = u
			}
			if c.AsPasswordProvider {
				m.PasswordProvider = u
			}
			if c.AsStatusProvider {
				m.StatusProvider = u
			}
			if c.AsRoleProvider {
				m.RoleProvider = u
			}
			return nil
		}), nil
	})
}

func init() {
	RegisterMemberDirective()
}
