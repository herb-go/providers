package statictoml

import (
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
	if err != ErrSuffixError {
		t.Fatal(err)
	}
	err = wrongsource.Save(r)
	if err != ErrSuffixError {
		t.Fatal(err)
	}
	_, err = os.Stat(string(wrongsource))
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	err = wrongsource.Load(r)
	if err != ErrSuffixError {
		t.Fatal(err)
	}
	_, err = os.Stat(string(wrongsource))
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
}
