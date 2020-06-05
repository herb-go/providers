package cacheproxyoverseer

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/worker"
)

//Config overseer config struct
type Config struct {
}

//Apply apply config to overseer
func (c *Config) Apply(o *worker.PlainOverseer) error {
	o.WithIntroduction("Cacheproxy workers")
	o.WithTrainFunc(func(w []*worker.Worker) error {
		for _, v := range w {
			proxy := GetCacheProxyByID(v.Name)
			if proxy == nil {
				continue
			}
			t := worker.GetTranning(v.Name)
			if t == nil {
				continue
			}
			config := &cache.OptionConfig{}
			err := t.TranningPlan(config)
			if err != nil {
				return err
			}
			proxycache := cache.New()
			err = config.ApplyTo(proxycache)
			if err != nil {
				return err
			}
			proxy.Cacheable = proxycache
		}
		return nil
	})
	return nil
}

//New create new config
func New() *Config {
	return &Config{}
}
