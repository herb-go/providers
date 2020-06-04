package hiredprotecter

import (
	"github.com/herb-go/herb/user/identifier"
)

type IdentifierConfig struct {
	Worker string
	Config func(interface{}) error
}

func (c *IdentifierConfig) CreateIdentifier() (identifier.Identifier, error) {
	return nil, nil
}
