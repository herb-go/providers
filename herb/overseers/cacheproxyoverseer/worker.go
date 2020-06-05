package cacheproxyoverseer

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/worker"
)

var cacheproxyworker = &cache.Proxy{}
var Team = worker.GetWorkerTeam(&cacheproxyworker)

func GetCacheProxyByID(id string) *cache.Proxy {
	w := worker.FindWorker(Team, id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(**cache.Proxy)
	if ok == false || c == nil {
		return nil
	}
	return *c
}
