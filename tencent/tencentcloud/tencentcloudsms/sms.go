package tencentcloudsms

import (
	"net/url"
	"strconv"

	"github.com/herb-go/fetcher"
	"github.com/herb-go/providers/tencent/tencentcloud"
)

type Sms struct {
	SdkAppid string
	*tencentcloud.App
}
type SendStatusSet struct {
	SerialNo       string
	PhoneNumber    string
	Fee            int
	SessionContext string
	Code           string
	Message        string
	IsoCode        string
}
type Response struct {
	SendStatusSet []SendStatusSet
	tencentcloud.BaseResponse
}

type Result struct {
	Response Response
}

func (s *Sms) Send(msg *Message) (*Result, error) {
	r := NewSMSRequest()
	q := r.URL.Query()
	for k, v := range msg.PhoneNumber {
		q.Set("PhoneNumberSet."+strconv.Itoa(k), v)
	}
	q.Set("TemplateID", msg.TemplateID)
	q.Set("SmsSdkAppid", s.SdkAppid)
	if msg.Sign != "" {
		q.Set("Sign", msg.Sign)
	}
	for k, v := range msg.TemplateParam {
		q.Set("TemplateParamSet."+strconv.Itoa(k), v)
	}
	if msg.ExtendCode != "" {
		q.Set("ExtendCode", msg.ExtendCode)
	}
	if msg.SessionContext != "" {
		q.Set("SessionContext", msg.SessionContext)
	}
	if msg.SenderID != "" {
		q.Set("SenderId", msg.SenderID)
	}
	r.URL.RawQuery = q.Encode()
	result := &Result{}
	_, err := r.CreatePreset(s.App, nil).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, err
	}
	err = result.Response.CodeError()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewSMSRequest() *tencentcloud.Request {
	r := tencentcloud.NewRequest()
	u, err := url.Parse("https://sms.tencentcloudapi.com")
	if err != nil {
		panic(err)
	}
	r.SetGET(u)
	r.Action = "SendSms"
	r.Version = "2019-07-11"
	r.Service = "sms"
	return r
}

type Message struct {
	TemplateID     string
	Sign           string
	PhoneNumber    []string
	TemplateParam  []string
	SessionContext string
	ExtendCode     string
	SenderID       string
}

func NewMessage() *Message {
	return &Message{}
}
