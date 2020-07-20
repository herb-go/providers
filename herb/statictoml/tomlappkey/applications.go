package tomlappkey

import (
	"sync"

	"github.com/herb-go/uniqueid"

	"github.com/herb-go/herbsecurity/authority"
	"github.com/herb-go/herbsecurity/authority/service/application"
	"github.com/herb-go/providers/herb/statictoml"
)

type Applications struct {
	Source       statictoml.Source
	locker       sync.Mutex
	Data         map[authority.Authority]*application.Verified
	KeyGenerator func() ([]byte, error)
	IDGenerator  func() (string, error)
}

func (apps *Applications) save() error {
	return apps.Source.Save(apps.Data)
}
func (apps *Applications) CreateApplication(p authority.Principal, a authority.Agent) (*application.Verified, error) {
	var err error
	if p == "" {
		return nil, authority.ErrEmptyPrincipal
	}
	apps.locker.Lock()
	defer apps.locker.Unlock()
	id, err := apps.IDGenerator()
	if err != nil {
		return nil, err
	}
	pass, err := apps.KeyGenerator()
	if err != nil {
		return nil, err
	}
	app := application.NewVerified()
	app.Principal = p
	app.Agent = a

	app.Authority = authority.Authority(id)
	app.Passphrase = authority.Passphrase(pass)
	apps.Data[app.Authority] = app
	err = apps.save()
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (apps *Applications) RegenerateApplication(p authority.Principal, a authority.Authority) error {
	apps.locker.Lock()
	defer apps.locker.Unlock()
	if a == "" {
		return authority.ErrEmptyAuthority
	}
	app := apps.Data[a]
	if app == nil || app.Principal != p {
		return authority.ErrNotFound
	}
	pass, err := apps.KeyGenerator()
	if err != nil {
		return err
	}
	app.Passphrase = authority.Passphrase(pass)
	apps.Data[app.Authority] = app
	return apps.save()
}

func (apps *Applications) RevokeApplication(p authority.Principal, a authority.Authority) error {
	apps.locker.Lock()
	defer apps.locker.Unlock()
	if a == "" {
		return authority.ErrEmptyAuthority
	}
	app := apps.Data[a]
	if app == nil || app.Principal != p {
		return authority.ErrNotFound
	}
	delete(apps.Data, a)
	return apps.save()
}

func (apps *Applications) LoadApplication(a authority.Authority) (*application.Verified, error) {
	apps.locker.Lock()
	defer apps.locker.Unlock()
	if a == "" {
		return nil, authority.ErrEmptyAuthority
	}
	return apps.Data[a], nil
}

func New() *Applications {
	return &Applications{
		Data:        map[authority.Authority]*application.Verified{},
		IDGenerator: uniqueid.GenerateID,
	}
}
