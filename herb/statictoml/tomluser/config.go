package tomluser

import (
	"sync"

	"github.com/herb-go/member"
	"github.com/herb-go/providers/herb/statictoml"
)

var locker sync.Mutex
var registered = map[statictoml.Source]*Users{}

func Flush() {
	locker.Lock()
	locker.Unlock()
	registered = map[statictoml.Source]*Users{}
}

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
	locker.Lock()
	locker.Unlock()
	source, err := c.Source.Abs()
	if err != nil {
		return nil, err
	}
	u, ok := registered[source]
	if ok && u != nil {
		return u, nil
	}
	u = NewUsers()
	u.Source = c.Source
	data := NewData()
	err = u.Source.Load(data)
	if err != nil {
		return nil, err
	}
	for k := range data.Users {
		u.addUser(data.Users[k])
	}
	return u, nil
}
func (c *Config) Execute(m *member.Service) error {
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
}

var DirectiveFactory = func(loader func(v interface{}) error) (member.Directive, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
