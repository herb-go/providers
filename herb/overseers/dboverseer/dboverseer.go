package dboverseer

import "github.com/herb-go/worker"

//Config overseer config struct
type Config struct {
}

//ApplyTo apply config to overseer
func (c *Config) ApplyTo(o *worker.PlainOverseer) error {
	o.WithIntroduction("Database workers")
	return nil
}

//New create new config
func New() *Config {
	return &Config{}
}
