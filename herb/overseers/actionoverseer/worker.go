package actionoverseer

import (
	"github.com/herb-go/herb/middleware/action"
	"github.com/herb-go/worker"
)

var actionworker = action.New(nil)
var Team = worker.GetWorkerTeam(&actionworker)

func GetActionByID(id string) *action.Action {
	w := worker.FindWorker(Team, id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(**action.Action)
	if ok == false || c == nil {
		return nil
	}
	return *c
}
