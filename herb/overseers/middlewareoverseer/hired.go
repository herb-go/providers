package middlewareoverseer

import (
	"net/http"

	"github.com/herb-go/worker"
)

var middlewareworker = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}
var Team = worker.GetWorkerTeam(&middlewareworker)

func GetMiddlewareByID(id string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w := worker.FindWorker(Team, id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc))
	if ok == false || c == nil {
		return nil
	}
	return *c
}
