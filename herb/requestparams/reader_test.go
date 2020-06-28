package requestparams_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/worker"

	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herb/user/httpuser"
	"github.com/herb-go/requestparams"
)

func read(name string, loader func(interface{}) error, r *http.Request) ([]byte, error) {
	f, err := requestparams.GetReaderFactory(name)
	if err != nil {
		return nil, err
	}
	reader, err := f.CreateReader(loader)
	if err != nil {
		return nil, err
	}
	return reader(r)
}

func readString(name string, loader func(interface{}) error, r *http.Request) (string, error) {
	b, err := read(name, loader, r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type testidentifier string

func (i testidentifier) IdentifyRequest(r *http.Request) (string, error) {
	return string(i), nil
}

func TestReader(t *testing.T) {
	var data string
	var err error
	var req *http.Request
	worker.Reset()
	defer worker.Reset()
	requestparams.Reset()
	defer requestparams.Reset()

	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	data, err = readString("test.notfound", nil, req)
	if data != "" || err == nil {
		t.Fatal(data, err)
	}
	commonconfig := &requestparams.CommonFieldConfig{
		Field: "test",
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	req.Header.Set("test", "test")
	data, err = readString("header", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	data, err = readString("query", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("POST", "http://127.0.0.1", bytes.NewBufferString("test=test"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	data, err = readString("form", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	router.GetParams(req).Set("test", "test")
	data, err = readString("router", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	data, err = readString("fixed", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	reqcookie := &http.Cookie{
		Name:  "test",
		Value: "test",
	}
	req.AddCookie(reqcookie)
	data, err = readString("cookie", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	data, err = readString("cookie", newLoader(commonconfig), req)
	if data != "" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	req.RemoteAddr = "127.0.0.1:8000"
	data, err = readString("ip", newLoader(commonconfig), req)
	if data != "127.0.0.1" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	data, err = readString("method", newLoader(commonconfig), req)
	if data != "GET" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = readString("path", newLoader(commonconfig), req)
	if data != "/test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = readString("host", newLoader(commonconfig), req)
	if data != "127.0.0.1" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	req.SetBasicAuth("testuser", "testpassword")
	data, err = readString("user", newLoader(commonconfig), req)
	if data != "testuser" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	req.SetBasicAuth("testuser", "testpassword")
	data, err = readString("password", newLoader(commonconfig), req)
	if data != "testpassword" || err != nil {
		t.Fatal(data, err)
	}
	identifier := httpuser.Identifier(testidentifier("test"))
	worker.Hire("test.identifier", &identifier)
	identifierconfig := &requestparams.WorkerConfig{ID: "test.identifier"}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = readString("identifier", newLoader(identifierconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	workernotfoundconfig := &requestparams.WorkerConfig{ID: "test.notfound"}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = readString("identifier", newLoader(workernotfoundconfig), req)
	if data != "" || err == nil {
		t.Fatal(data, err)
	}
	reader := requestparams.Reader(func(r *http.Request) ([]byte, error) {
		return []byte("test"), nil
	})
	worker.Hire("test.reader", &reader)
	hiredconfig := &requestparams.WorkerConfig{ID: "test.reader"}
	data, err = readString("hired", newLoader(hiredconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	data, err = readString("hired", newLoader(workernotfoundconfig), req)
	if data != "" || err == nil {
		t.Fatal(data, err)
	}
	worker.Hire("test.readerfactory", &requestparams.HostReaderFactory)
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = readString("test.readerfactory", newLoader(commonconfig), req)
	if data != "127.0.0.1" || err != nil {
		t.Fatal(data, err)
	}
}
