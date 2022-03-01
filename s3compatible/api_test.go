package s3compatible

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	var content = []byte("testdata")
	api := New()
	err := testConfig.ApplyTo(api)
	if err != nil {
		panic(err)
	}
	err = api.Remove(context.TODO(), testBucket, testFilename)
	if err != nil && !IsHTTPError(err, 404) {
		panic(err)
	}
	defer api.Remove(context.TODO(), testBucket, testFilename)
	info, err := api.Info(context.TODO(), testBucket, testFilename)
	if info != nil || err == nil || !IsHTTPError(err, 404) {
		t.Fatal(info, err)
	}
	err = api.Save(context.TODO(), testBucket, testFilename, bytes.NewBuffer(content))
	if err != nil {
		panic(err)
	}
	info, err = api.Info(context.TODO(), testBucket, testFilename)
	if info == nil || err != nil || info.Size != int64(len(content)) {
		t.Fatal(info, err)
	}
	var body = bytes.NewBuffer(nil)
	n, err := api.Load(context.TODO(), testBucket, testFilename, body)
	if n != 8 || err != nil {
		t.Fatal(n, err)
	}
	if body.String() != string(content) {
		t.Fatal(body.String())
	}
	err = api.Remove(context.TODO(), testBucket, testFilename)
	if err != nil {
		panic(err)
	}
	info, err = api.Info(context.TODO(), testBucket, testFilename)
	if info != nil || err == nil || !IsHTTPError(err, 404) {
		t.Fatal(info, err)
	}
	url, err := api.PresignPutObject(context.TODO(), testBucket, testFilename, time.Hour)
	if url == "" || err != nil {
		t.Fatal(url, err)
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(content))
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	url, err = api.PresignGetObject(context.TODO(), testBucket, testFilename, time.Hour)
	if url == "" || err != nil {
		t.Fatal(url, err)
	}
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if !bytes.Equal(data, content) {
		t.Fatal(data)
	}
	info, err = api.Info(context.TODO(), testBucket, testFilename)
	if info == nil || err != nil || info.Size != int64(len(content)) {
		t.Fatal(info, err)
	}
}
