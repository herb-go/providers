package tomlappkey

import (
	"github.com/herb-go/misc/generator"
	"github.com/herb-go/providers/herb/statictoml"
)

var DefaultKeycharList = generator.AlphanumericList
var DefaultMin = 32

type Config struct {
	Source      statictoml.Source
	KeycharList string
	KeyMin      int
	KeyMax      int
}

func (c *Config) Create() (*Applications, error) {
	apps := New()
	keygenerator := &generator.ListGenerator{}
	if c.KeycharList == "" {
		keygenerator.List = DefaultKeycharList
	} else {
		keygenerator.List = []byte(c.KeycharList)
	}
	if c.KeyMin == 0 {
		keygenerator.Min = DefaultMin
	} else {
		keygenerator.Min = c.KeyMin
	}
	keygenerator.Max = c.KeyMax
	apps.KeyGenerator = keygenerator.Generate
	return apps, nil
}
