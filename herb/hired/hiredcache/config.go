package hiredcache

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/providers/herb/overseers/cacheoverseer"
	"github.com/herb-go/worker"
)

type Config struct {
	ID     string
	Prefix string
	AsNode bool
}

func Register() {
	cache.Register("hiredcache", func(loader func(interface{}) error) (cache.Driver, error) {
		c := &Config{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		d := New()
		cacheproxy := cacheoverseer.GetCacheByID(c.ID)
		if cacheproxy == nil {
			return nil, worker.ErrWorkerNotFound
		}
		if c.Prefix == "" {
			d.Cacheable = cacheproxy
		} else {
			if c.AsNode {
				d.Cacheable = cache.NewNode(cacheproxy, c.Prefix)
			} else {
				d.Cacheable = cache.NewCollection(cacheproxy, c.Prefix, cache.DefaultTTL)
			}
		}
		return d, nil
	})
}

func init() {
	Register()
}
