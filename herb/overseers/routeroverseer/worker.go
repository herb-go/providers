package routeroverseer

import (
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/worker"
)

var routerworker *router.Factory
var Team = worker.GetWorkerTeam(&routerworker)

func GetRouterByID(id string) *router.Factory {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(**router.Factory)
	if ok == false || c == nil {
		return nil
	}
	return *c
}
