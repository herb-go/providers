package hiredprotecter

import (
	"github.com/herb-go/herb/user/identifier/protecter"
)

type Config struct {
	Credentialers []*CredentialerConfig
	Identifier    *IdentifierConfig
	OnFail        *HandlerConfig
}

func (c *Config) CreateProtecter() (*protecter.Protecter, error) {
	var err error
	p := protecter.New()
	credentialers := make([]protecter.Credentialer, len(c.Credentialers))
	for k := range c.Credentialers {
		credentialers[k], err = c.Credentialers[k].CreateCredentialerr()
		if err != nil {
			return nil, err
		}
	}
	if c.Identifier != nil {
		i, err := c.Identifier.CreateIdentifier()
		if err != nil {
			return nil, err
		}
		p.Identifier = i
	}
	if c.OnFail != nil {
		h, err := c.OnFail.CreateHandler()
		if err != nil {
			return nil, err
		}
		p.OnFail = h
	}
	return p, nil
}
