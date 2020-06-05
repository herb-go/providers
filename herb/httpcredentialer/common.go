package httpcredentialer

import (
	"net"
	"net/http"

	"github.com/herb-go/herb/user/identifier"
	"github.com/herb-go/herb/user/identifier/protecter"
)

type FieldConfig struct {
	Source string
	Name   string
}

type CustomFieldConfig struct {
	CredentialType identifier.CredentialType
	FieldConfig
}

func NewFieldFactory(ctype identifier.CredentialType) Factory {
	return FactoryFunc(func(loader func(interface{}) error) (protecter.Credentialer, error) {
		c := &FieldConfig{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		return NewFieldCredentialLoader(ctype, c.Source, c.Name)
	})
}

var CustomFactory Factory = FactoryFunc(func(loader func(interface{}) error) (protecter.Credentialer, error) {
	c := &CustomFieldConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return NewFieldCredentialLoader(c.CredentialType, c.Source, c.Name)
})

var iploader = func(r *http.Request) (identifier.CredentialData, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return nil, err
	}
	return identifier.CredentialData(host), nil

}
var NewIPFactory = func(ctype identifier.CredentialType) Factory {
	return FactoryFunc(func(loader func(interface{}) error) (protecter.Credentialer, error) {
		l := &protecter.CredentialLoader{
			CredentialType: ctype,
			LoaderFunc:     iploader,
		}
		return l, nil
	})
}
