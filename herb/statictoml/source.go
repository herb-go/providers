package statictoml

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

var DefaultFileMode = os.FileMode(0700)

type Source string

func (s Source) Save(v interface{}) error {
	buf := bytes.NewBuffer(nil)
	e := toml.NewEncoder(buf)
	err := e.Encode(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(string(s), buf.Bytes(), DefaultFileMode)
}

func (s Source) Load(v interface{}) error {
	data, err := ioutil.ReadFile(string(s))
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, v)
}
