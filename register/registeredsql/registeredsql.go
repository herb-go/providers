package registeredsql

import (
	"sync"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/providers/register"
)

type objectType = db.Database

var Type = register.RegisterType("sql")

var registeredKeys = register.New(Type)
var registerlock = sync.Mutex{}
var registered = make(map[string]objectType)

func Register(key string, object objectType) error {
	registerlock.Lock()
	defer registerlock.Unlock()
	err := registeredKeys.RegisterKey(key)
	if err != nil {
		return err
	}
	registered[key] = object
	return nil
}

func Unregister(key string) {
	registerlock.Lock()
	defer registerlock.Unlock()
	registeredKeys.UnregisterKey(key)
	delete(registered, key)
}

func Get(key string) (object objectType, err error) {
	var ok bool
	registerlock.Lock()
	defer registerlock.Unlock()
	object, ok = registered[key]
	if ok == false {
		err = &register.NotRegsiteredError{
			Type: Type,
			Key:  key,
		}
		return
	}
	return object, nil
}

func MustGet(key string) objectType {
	object, err := Get(key)
	if err != nil {
		panic(err)
	}
	return object
}
