package registeredidgenerator

import (
	"sync"

	"github.com/herb-go/providers/register"
)

type objectType = func() (string, error)

var Type = register.RegisterType("idgenerator")

var registeredKeys = register.New(Type)
var registerlock = sync.Mutex{}
var registered = make(map[string]objectType)

var DefaultKey string

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

func MustRegister(key string, object objectType) {
	err := Register(key, object)
	if err != nil {
		panic(err)
	}
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
	if key == "" {
		key = DefaultKey
	}
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
