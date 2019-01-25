package uuidgenerator

import (
	"github.com/herb-go/providers/register/registeredidgenerator"
	"github.com/satori/go.uuid"
)

var UuidV1 = func() (string, error) {
	unid, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	return unid.String(), nil
}

func init() {
	err := registeredidgenerator.Register("uuidv1", UuidV1)
	if err != nil {
		panic(err)
	}
}
