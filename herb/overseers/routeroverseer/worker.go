package routeroverseer

import (
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/worker"
)

var routerworker router.Router
var Team = worker.GetWorkerTeam(&routerworker)

func GetRouterByID(id string) router.Router {
	w := worker.FindWorker(Team, id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*router.Router)
	if ok == false || c == nil {
		return nil
	}
	return *c
}
