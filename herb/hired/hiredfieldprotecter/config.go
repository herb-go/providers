package hiredfieldprotecter

import (
	"fmt"

	"github.com/herb-go/herbsecurity/authority/credential"
	"github.com/herb-go/httpinfomanager"
	"github.com/herb-go/protecter"
	actionoverseer "github.com/herb-go/providers/herb/overseers/actionoverseer"
	authenticatorfactoryoverseer "github.com/herb-go/providers/herb/overseers/authenticatorfactoryoverseer"
	"github.com/herb-go/worker"
)

type Config struct {
	Fields       map[string]*httpinfomanager.FieldName
	OnFailAction string
	AuthType     string
	AuthConfig   func(v interface{}) error `config:", lazyload"`
}

func (c *Config) ApplyTo(p *protecter.Protecter) error {
	var credentialers []protecter.Credentialer
	for k, v := range c.Fields {
		f, err := v.Field()
		if err != nil {
			return err
		}
		credentialers = append(credentialers, &Credentialer{
			credentialName: credential.Name(k),
			field:          f,
		})
	}
	p.Credentialers = credentialers
	if c.OnFailAction != "" {
		a := actionoverseer.GetActionByID(c.OnFailAction)
		if a == nil {
			return fmt.Errorf("%w (%s)", worker.ErrWorkerNotFound, c.OnFailAction)
		}
		p.OnFail = a
	}
	if c.AuthType != "" {
		authfactory := authenticatorfactoryoverseer.GetAuthenticatorFactoryByID(c.AuthType)
		if authfactory == nil {
			return fmt.Errorf("%w (%s)", worker.ErrWorkerNotFound, c.AuthType)
		}
		auth, err := authfactory.CreateAuthenticator(c.AuthConfig)
		if err != nil {
			return err
		}
		p.Authenticator = auth
	}
	return nil
}
