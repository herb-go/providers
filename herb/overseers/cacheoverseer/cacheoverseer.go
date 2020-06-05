package cacheoverseer

import "github.com/herb-go/worker"

//Config overseer config struct
type Config struct {
}

//Apply apply config to overseer
func (c *Config) Apply(o *worker.PlainOverseer) error {
	o.WithIntroduction("Cache workers")
	return nil
}

//New create new config
func New() *Config {
	return &Config{}
}
