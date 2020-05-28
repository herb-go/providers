package hiredcache

import "github.com/herb-go/herb/cache"

type Driver struct {
	cache.Cacheable
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (d *Driver) SetGCErrHandler(f func(err error)) {

}
func New() *Driver {
	return &Driver{}
}
