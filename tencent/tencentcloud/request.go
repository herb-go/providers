package tencentcloud

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/herb-go/fetcher"
)

type Request struct {
	Action      string
	Region      string
	Version     string
	URL         *url.URL
	Method      string
	ContentType ContentType
	Header      http.Header
	Body        []byte
	Token       string
	Timestamp   int64
	Service     string
}

func (r *Request) CreateSignData(signed_headers []string) *SignData {
	data := &SignData{}
	r.Header.Set("Host", r.URL.Host)
	data.Timestamp = r.Timestamp
	data.Service = r.Service
	data.HTTPRequestMethod = strings.ToUpper(r.Method)
	if data.HTTPRequestMethod == "" {
		data.HTTPRequestMethod = "GET"
	}
	data.CanonicalURI = "/"
	if data.HTTPRequestMethod != "POST" {
		data.CanonicalQueryString = r.URL.Query().Encode()
	}
	var headers = make([]string, len(signed_headers))
	for k := range signed_headers {
		headers[k] = strings.ToLower(signed_headers[k])
	}
	headers = append(headers, "content-type", "host")
	sort.Strings(headers)
	data.SignedHeaders = strings.Join(headers, ";")
	for k := range headers {
		field := strings.ToLower(strings.TrimSpace(headers[k]))
		value := strings.ToLower(strings.TrimSpace(r.Header.Get(field)))
		data.CanonicalHeaders = data.CanonicalHeaders + fmt.Sprintf("%s:%s\n", field, value)
	}
	b := sha256.Sum256(r.Body)
	data.HashedRequestPayload = hex.EncodeToString(b[:])
	return data
}

func (r *Request) CreatePreset(a *App, signedHeaders []string) *fetcher.Preset {
	header := http.Header{
		"X-TC-Action":    []string{r.Action},
		"X-TC-Timestamp": []string{strconv.FormatInt(r.Timestamp, 10)},
		"X-TC-Version":   []string{r.Version},
		"Authorization":  []string{r.CreateSignData(signedHeaders).Authorization(a.SecretID, a.SecretKey)},
	}
	if r.Region != "" {
		header.Set("X-TC-Region", r.Region)
	}
	if r.Token != "" {
		header.Set("X-TC-Token", r.Token)
	}
	return fetcher.NewPreset().With(
		&a.Client,
		fetcher.ParsedURL(r.URL),
		fetcher.Header(header),
		fetcher.Method(r.Method),
		fetcher.Header(r.Header),
	)
}
func (r *Request) SetGET(url *url.URL) {
	r.Method = "GET"
	r.Header.Set("Content-Type", string(ContentTypeURLEncoded))
	r.URL = url
}

func (r *Request) MustSetPOSTJSON(url *url.URL, body interface{}) {
	r.Method = "POST"
	r.Header.Set("Content-Type", string(ContentTypeJSON))
	bs, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	r.Body = bs
	r.URL = url
}

func (r *Request) SetPOSTFormdata(url *url.URL, body []byte) {
	r.Method = "POST"
	r.Body = body
	r.URL = url
}

func NewRequest() *Request {
	return &Request{
		Timestamp: time.Now().Unix(),
		Header:    http.Header{},
	}
}
