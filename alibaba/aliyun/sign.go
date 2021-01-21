package aliyun

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var replacer = strings.NewReplacer(
	"+", "%20",
	"*", "%2A",
	"%7E", "~",
)

func HMAC_Sha1(key []byte, msg []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

func escape(value string) string {
	escaped := url.QueryEscape((value))
	return replacer.Replace(escaped)
}

type SignParam struct {
	Name  string
	Value string
}

type SignParams []*SignParam

// Len is the number of elements in the collection.
func (sp SignParams) Len() int {
	return len(sp)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (sp SignParams) Less(i, j int) bool {
	return sp[i].Name < sp[j].Name
}

// Swap swaps the elements with indexes i and j.
func (sp SignParams) Swap(i, j int) {
	p := sp[j]
	sp[j] = sp[i]
	sp[i] = p
}
func (sp SignParams) CanonicalizedQueryString() string {
	sort.Sort(sp)
	query := make([]string, 0, len(sp))
	for _, v := range sp {
		query = append(query, escape(v.Name)+"="+escape(v.Value))
	}
	result := strings.Join(query, "&")
	return result
}
func (sp SignParams) ToSign(HTTPMethod string) string {
	qs := sp.CanonicalizedQueryString()
	tosign := HTTPMethod + "&" + escape("/") + "&" + escape(qs)
	return tosign
}

type Params map[string]string

func (p Params) With(name string, value string) Params {
	p[name] = value
	return p
}
func (p Params) WithBool(name string, value bool) Params {
	if value {
		p[name] = "true"
	} else {
		p[name] = "false"
	}
	return p
}

func (p Params) WithInt(name string, value int) Params {
	p[name] = strconv.Itoa(value)
	return p
}
func (p Params) WithAction(value string) Params {
	p["Action"] = value
	return p
}
func (p Params) WithAccessKeyId(value string) Params {
	p["AccessKeyId"] = value
	return p
}

func (p Params) WithNonce() Params {
	var data = make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	p["SignatureNonce"] = base64.StdEncoding.EncodeToString(data)
	return p
}
func (p Params) WithTimestamp(ts int64) Params {
	p["Timestamp"] = time.Unix(ts, 0).UTC().Format(time.RFC3339)
	return p
}
func (p Params) WithVersion(value string) Params {
	p["Version"] = value
	return p
}

func (p Params) ToSignParams() SignParams {
	sp := make(SignParams, 0, len(p))
	for k := range p {
		if p[k] != "" {
			sp = append(sp, &SignParam{Name: k, Value: p[k]})
		}
	}
	return sp
}

func (p Params) Sign(HTTPMethod string, key string) string {
	sp := p.ToSignParams()
	tosign := sp.ToSign(HTTPMethod)
	return base64.StdEncoding.EncodeToString(HMAC_Sha1([]byte(key+"&"), []byte(tosign)))
}

func (p Params) SignedQuery(HTTPMethod string, key string) url.Values {
	signed := p.Sign(HTTPMethod, key)
	query := p.ToSignParams().CanonicalizedQueryString()

	v, err := url.ParseQuery("Signature=" + signed + "&" + query)
	if err != nil {
		panic(err)
	}
	return v
}

func NewParams(AccessKey *AccessKey) Params {
	return Params{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"Format":           "json",
		"AccessKeyId":      AccessKey.AccessKeyID,
	}
}
