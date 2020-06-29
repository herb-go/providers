package memberdirectivefactoryoverseer

import (
	"github.com/herb-go/member"
	"github.com/herb-go/worker"
)

var factoryworker func(loader func(v interface{}) error) (member.Directive, error)
var Team = worker.GetWorkerTeam(&factoryworker)

func GetMemberDirectiveFactoryByID(id string) func(loader func(v interface{}) error) (member.Directive, error) {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*func(loader func(v interface{}) error) (member.Directive, error))
	if ok == false || c == nil {
		return nil
	}
	return *c
}
