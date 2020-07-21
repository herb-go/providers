package tomlappkey

import (
	"github.com/herb-go/herbsecurity/authority/credential"
	"github.com/herb-go/misc/generator"
	"github.com/herb-go/protecter/authenticator"
	"github.com/herb-go/protecter/authenticator/appsecret"
	"github.com/herb-go/providers/herb/statictoml"
)

var DefaultKeycharList = generator.AlphanumericList
var DefaultMin = 32

type Config struct {
	Source      statictoml.Source
	KeycharList string
	KeyMin      int
	KeyMax      int
}

func (c *Config) Create() (*Applications, error) {
	apps := New()
	keygenerator := &generator.ListGenerator{}
	if c.KeycharList == "" {
		keygenerator.List = DefaultKeycharList
	} else {
		keygenerator.List = []byte(c.KeycharList)
	}
	if c.KeyMin == 0 {
		keygenerator.Min = DefaultMin
	} else {
		keygenerator.Min = c.KeyMin
	}
	keygenerator.Max = c.KeyMax
	apps.KeyGenerator = keygenerator.Generate
	apps.Source = c.Source
	err := apps.load()
	if err != nil {
		return nil, err
	}
	return apps, nil
}
func CreateAuthenticator(loader func(interface{}) error) (credential.Authenticator, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	apps, err := c.Create()
	if err != nil {
		return nil, err
	}
	a := appsecret.New()
	a.Loader = apps
	return a, nil
}

var AuthenticatorFactory authenticator.AuthenticatorFactory = authenticator.AuthenticatorFactoryFunc(CreateAuthenticator)

func NewConfig() *Config {
	return &Config{}
}
