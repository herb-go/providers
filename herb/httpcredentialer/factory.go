package httpcredentialer

import (
	"github.com/herb-go/herb/user/identifier/protecter"
)

type Factory interface {
	CreateHTTPCredentialer(func(interface{}) error) (protecter.Credentialer, error)
}

type FactoryFunc func(func(interface{}) error) (protecter.Credentialer, error)

func (f FactoryFunc) CreateHTTPCredentialer(loader func(interface{}) error) (protecter.Credentialer, error) {
	return f(loader)
}
