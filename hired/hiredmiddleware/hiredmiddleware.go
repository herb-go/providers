package hiredmiddleware

import (
	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/middlewarefactory"
	"github.com/herb-go/worker"
	"github.com/herb-go/worker/overseers/middlewareoverseer"
)

type Config struct {
	ID string
}

var NewFactory = func() middlewarefactory.Factory {
	return func(loader func(v interface{}) error) (middleware.Middleware, error) {
		c := &Config{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		m := middlewareoverseer.GetMiddlewareByID(c.ID)
		if m == nil {
			return nil, worker.ErrWorkerNotFound
		}
		return m, nil
	}
}
