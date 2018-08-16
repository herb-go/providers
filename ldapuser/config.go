package ldapuser

import (
	"gopkg.in/ldap.v2"
)

type Config struct {
	Net          string
	Addr         string
	UserDN       string
	BindDN       string
	BindPass     string
	SearchDN     string
	SearchFilter string
}

func (c *Config) PasswordProvider() *PasswordProvider {
	return &PasswordProvider{
		Config: c,
	}
}
func (c *Config) Dial() (*ldap.Conn, error) {
	return ldap.Dial(c.Net, c.Addr)
}

func (c *Config) DialBound() (*ldap.Conn, error) {
	l, err := c.Dial()
	if err != nil {
		return nil, err
	}
	err = l.Bind(c.BindDN, c.BindPass)
	if err != nil {
		return nil, err
	}
	return l, nil
}
