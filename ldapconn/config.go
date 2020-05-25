package ldapconn

import (
	"fmt"

	"gopkg.in/ldap.v2"
)

type Config struct {
	Net          string
	Addr         string
	UserPattern  string
	BindDN       string
	BindPass     string
	SearchDN     string
	SearchFilter string
	GroupDN      string
	GroupIDField string
	GroupFilter  string
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
func (c *Config) BindUser(uid, password string) (*ldap.Conn, error) {
	uid = ldap.EscapeFilter(uid)
	l, err := c.Dial()
	if err != nil {
		return nil, err
	}
	err = l.Bind(fmt.Sprintf(c.UserPattern, uid), password)
	return l, err
}
