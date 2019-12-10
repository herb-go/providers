package cos

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/herb-go/herb/file/store"
	"github.com/herb-go/providers/tencent/tencentcloud"
)

func HmacSha1(key string, dst string) (string, error) {
	mac := hmac.New(sha1.New, []byte(key))
	_, err := mac.Write([]byte(dst))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

const timeHalfLife = 5 * time.Minute

type CloudObjectStorage struct {
	App    tencentcloud.App
	Region string
	Bucket string
}

func (s *CloudObjectStorage) Authorization(method string, requrl url.URL, header http.Header, body io.Reader) (*http.Request, error) {
	start := time.Now().Add(-timeHalfLife).Unix()
	end := time.Now().Add(timeHalfLife).Unix()
	signTime := fmt.Sprintf("%d;%d", start, end)
	if header == nil {
		header = http.Header{}
	}
	header.Set("Host", requrl.Host)
	headerList := []string{}
	headers := []string{}
	for k := range header {
		headerList = append(headerList, strings.ToLower(k))
		headers = append(headers, strings.ToLower(k)+"="+url.QueryEscape(header.Get(k)))
	}
	sort.Strings(headerList)
	sort.Strings(headers)
	urlParamList := []string{}
	httpParams := []string{}
	q := requrl.Query()
	for k := range q {
		urlParamList = append(urlParamList, strings.ToLower(k))
		httpParams = append(httpParams, strings.ToLower(k+"="+q.Get(k)))
	}
	sort.Strings(urlParamList)
	sort.Strings(httpParams)
	httpString := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		strings.ToLower(method),
		requrl.RequestURI(),
		strings.Join(httpParams, "&"),
		strings.Join(headers, "&"),
	)
	signKey, err := HmacSha1(s.App.SecretKey, signTime)
	if err != nil {
		return nil, err
	}
	h := sha1.New()
	h.Write([]byte(httpString))
	sha1edHTTPString := hex.EncodeToString(h.Sum(nil))
	stringToSign := fmt.Sprintf("sha1\n%s\n%s\n", signTime, sha1edHTTPString)
	sign, err := HmacSha1(signKey, stringToSign)
	if err != nil {
		return nil, err
	}
	authorization := fmt.Sprintf(
		"q-sign-algorithm=sha1&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=%s&q-url-param-list=%s&q-signature=%s",
		s.App.SecretID,
		signTime,
		signTime,
		strings.Join(headerList, "&"),
		strings.Join(urlParamList, "&"),
		sign,
	)
	u := requrl.String()
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	header.Set("Authorization", authorization)
	for k := range header {
		req.Header.Set(k, header.Get(k))
	}
	return req, nil
}
func (s *CloudObjectStorage) API() string {
	return fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", s.Bucket, s.App.AppID, s.Region)
}
func (s *CloudObjectStorage) Save(filename string, reader io.Reader) (string, int64, error) {
	body := &bytes.Buffer{}
	size, err := io.Copy(body, reader)
	apiurl, err := url.Parse(s.API() + path.Join("/", filename))
	if err != nil {
		return "", 0, err
	}
	req, err := s.Authorization("PUT", *apiurl, nil, body)
	if err != nil {
		return "", 0, err
	}
	resp, err := s.App.Clients.Fetch(req)
	if err != nil {
		return "", 0, err
	}
	if resp.StatusCode != 200 {
		return "", 0, errors.New(resp.Status + string(resp.BodyContent))
	}
	return filename, size, nil
}
func (s *CloudObjectStorage) Load(id string) (io.ReadCloser, error) {
	apiurl, err := url.Parse(s.API() + path.Join("/", id))
	if err != nil {
		return nil, err
	}
	req, err := s.Authorization("GET", *apiurl, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.App.Clients.Client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, store.NewNotExistsError(id)
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(resp.Status + string(content))
	}
	return resp.Body, nil
}
func (s *CloudObjectStorage) Remove(id string) error {
	apiurl, err := url.Parse(s.API() + path.Join("/", id))
	if err != nil {
		return err
	}
	req, err := s.Authorization("DELETE", *apiurl, nil, nil)
	if err != nil {
		return err
	}
	resp, err := s.App.Clients.Fetch(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return errors.New(resp.Status + string(resp.BodyContent))
	}
	return err
}
func (s *CloudObjectStorage) URL(id string) (string, error) {
	return s.API() + path.Join("/", id), nil
}

func register() {
	store.Register("tencentcos", func(loader func(interface{}) error) (store.Driver, error) {
		c := &CloudObjectStorage{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}
