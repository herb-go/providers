package identifieroverseer

import (
	"github.com/herb-go/herb/user/httpuser"
	"github.com/herb-go/worker"
)

var identifierworker httpuser.Identifier
var Team = worker.GetWorkerTeam(&identifierworker)

func GetIdentifierByID(id string) httpuser.Identifier {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*httpuser.Identifier)
	if ok == false || c == nil {
		return nil
	}
	return *c
}
