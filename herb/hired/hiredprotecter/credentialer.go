package hiredprotecter

import "github.com/herb-go/herb/user/identifier/protecter"

type CredentialerConfig struct {
	Worker string
	Config func(interface{}) error
}

func (c *CredentialerConfig) CreateCredentialerr() (protecter.Credentialer, error) {
	return nil, nil
}
