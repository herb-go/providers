package register

import (
	"fmt"
	"sync"
)

var Debug = false

type RegisterType string
type DuplicationError struct {
	Type RegisterType
	Key  string
}

func (e *DuplicationError) Error() string {
	return "register error: \"" + e.Key + "\" has registered to type \"" + string(e.Type) + "\""
}

func IsDuplicationError(err error) bool {
	_, ok := err.(*DuplicationError)
	return ok
}

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
