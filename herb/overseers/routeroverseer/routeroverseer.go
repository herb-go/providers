package routeroverseer

import (
	"github.com/herb-go/herb/middleware/middlewarefactory"
	"github.com/herb-go/worker"
)

//Config overseer config struct
type Config struct {
}

//Apply apply config to overseer
func (c *Config) Apply(o *worker.PlainOverseer) error {
	o.WithIntroduction("HTTP Router workers")
	o.WithTrainFunc(func(w []*worker.Worker) error {
		for _, v := range w {
			router := GetRouterByID(v.Name)
			if router == nil {
				continue
			}
			t := worker.GetTranning(v.Name)
			if t == nil {
				continue
			}
			config := &middlewarefactory.ConfigList{}
			err := t.TranningPlan(config)
			if err != nil {
				return err
			}
			m, err := config.Middleware(middlewarefactory.DefaultContext)
			if err != nil {
				return err
			}
			router.Middlewares().Use(m)
		}
		return nil
	})
	return nil
}

//New create new config
func New() *Config {
	return &Config{}
}
