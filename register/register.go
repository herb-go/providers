package register

import (
	"fmt"
	"sync"
)

type RegisterType string

type Register struct {
	Type          RegisterType
	lock          sync.Mutex
	regiteredKeys map[string]bool
}

func (r *Register) duplicationError(key string) error {
	return &DuplicationError{
		Type: r.Type,
		Key:  key,
	}
}

func (r *Register) success(key string) {
	if Debug {
		fmt.Println("Register: \"" + key + " (type:\"" + string(r.Type) + "\") regstered")
	}
}

func (r *Register) RegisterKey(key string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if key == "" {
		return EmptyKeyError
	}
	_, ok := r.regiteredKeys[key]
	if ok {
		return r.duplicationError(key)
	}
	r.success(key)
	r.regiteredKeys[key] = true
	return nil
}

func (r *Register) UnregisterKey(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.regiteredKeys, key)
}

func New(registerType RegisterType) *Register {
	return &Register{
		Type:          registerType,
		regiteredKeys: map[string]bool{},
	}
}
