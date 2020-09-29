package statictoml

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

type testResult map[string]string

func TestSource(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if tmpdir == "" || err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)
	err = ioutil.WriteFile(tmpdir+"/test.static.toml", []byte{}, 0700)
	if err != nil {
		panic(err)
	}
	s := Source(tmpdir + "/test.static.toml")
	abs, err := s.Abs()
	if err != nil {
		panic(err)
	}
	if abs != s {
		t.Fatal(abs)
	}
	r := &testResult{}
	err = s.Load(r)
	if err != nil {
		panic(err)
	}
	(*r)["test"] = "testvalue"
	err = s.Save(r)
	if err != nil {
		panic(err)
	}
	r2 := &testResult{}
	err = s.Load(r2)
	if err != nil {
		panic(err)
	}
	if len((*r2)) != 1 || (*r2)["test"] != "testvalue" {
		t.Fatal(r2)
	}
	wrongsource := Source(tmpdir + "/test")
	err = wrongsource.Verify()
	if errors.Unwrap(err) != ErrSuffixError {
		t.Fatal(err)
	}
	err = wrongsource.Save(r)
	if errors.Unwrap(err) != ErrSuffixError {
		t.Fatal(err)
	}
	_, err = os.Stat(string(wrongsource))
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	err = wrongsource.Load(r)
	if errors.Unwrap(err) != ErrSuffixError {
		t.Fatal(err)
	}
	_, err = os.Stat(string(wrongsource))
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestExample(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if tmpdir == "" || err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)
	example := Source(tmpdir + "/example.static.toml")
	v := map[string]interface{}{"test": "test"}
	err = example.Save(v)
	if err != nil {
		panic(err)
	}
	s := Source(tmpdir + "/test.static.toml")
	notexist := Source(tmpdir + "/notexist.static.toml")
	err = example.VerifyWithExample("")
	if err != nil {
		panic(err)
	}
	err = Source("wrong").VerifyWithExample("")
	if errors.Unwrap(err) != ErrSuffixError {
		t.Fatal(err)
	}
	err = s.VerifyWithExample("wrong")
	if errors.Unwrap(err) != ErrSuffixError {
		t.Fatal(err)
	}
	err = s.VerifyWithExample(notexist)
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	err = s.VerifyWithExample(example)
	if err != nil {
		t.Fatal(err)
	}
	err = s.VerifyWithExample(notexist)
	if err != nil {
		t.Fatal(err)
	}
	var result = map[string]interface{}{}
	err = s.Load(&result)
	if err != nil {
		t.Fatal(err)
	}
	if result["test"] != "test" {
		t.Fatal(result)
	}
}
