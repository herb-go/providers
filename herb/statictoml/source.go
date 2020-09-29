package statictoml

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

//Suffix static toml path suffix
var Suffix = ".static.toml"

//ErrSuffixError error raised when source path suffix error
var ErrSuffixError = errors.New("statictoml: toml source must end with " + Suffix)

//FileMode source data file mode.
var FileMode = os.FileMode(0700)

//Source static toml data source.
type Source string

//Verify  check source path validity
func (s Source) Verify() error {
	if !strings.HasSuffix(string(s), Suffix) {
		return fmt.Errorf("%w (%s)", ErrSuffixError, s)
	}
	return nil
}
func (s Source) Abs() (Source, error) {
	p, err := filepath.Abs(string(s))
	if err != nil {
		return "", err
	}
	return Source(p), nil
}

//Save value to source.
func (s Source) Save(v interface{}) error {
	err := s.Verify()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	e := toml.NewEncoder(buf)
	err = e.Encode(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(string(s), buf.Bytes(), FileMode)
}

//Load value  from source.
func (s Source) Load(v interface{}) error {
	err := s.Verify()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(string(s))
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, v)
}

func (s Source) VerifyWithExample(example Source) error {
	err := s.Verify()
	if err != nil {
		return err
	}
	if example == "" {
		return nil
	}
	err = example.Verify()
	if err != nil {
		return err
	}
	_, err = os.Stat(string(s))
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		data, err := ioutil.ReadFile(string(example))
		if err != nil {
			return err
		}
		return ioutil.WriteFile(string(s), data, FileMode)
	}
	return err
}
