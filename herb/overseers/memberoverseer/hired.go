package memberoverseer

import (
	"github.com/herb-go/member"
	"github.com/herb-go/worker"
)

var directivefactoryworker member.DirectiveFactory
var Team = worker.GetWorkerTeam(&directivefactoryworker)

func GetDirectiveFactoryByID(id string) member.DirectiveFactory {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*member.DirectiveFactory)
	if ok == false || c == nil {
		return nil
	}
	return *c
}
