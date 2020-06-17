package tomlmember

import (
	"github.com/herb-go/member"
	"github.com/herb-go/providers/herb/statictoml"
)

type Data struct {
	Users []*User
}

func NewData() *Data {
	return &Data{}
}

type Config struct {
	Source             statictoml.Source
	AsPasswordProvider bool
	AsStatusProvider   bool
	AsAccountsProvider bool
	AsRoleProvider     bool
	HashMode           string
}

func (c *Config) Load() (*Users, error) {
	u := newUsers()
	u.Source = c.Source
	data := NewData()
	err := u.Source.Load(data)
	if err != nil {
		return nil, err
	}
	for k := range data.Users {
		u.addUser(data.Users[k])
	}
	return u, nil
}

var DirectiveFactory = func(loader func(v interface{}) error) (member.Directive, error) {
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
}
