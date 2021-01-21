package aliyunsms

import (
	"time"

	"github.com/herb-go/fetcher"
	"github.com/herb-go/providers/alibaba/aliyun"
)

type Result struct {
	BizId     string
	Code      string
	Message   string
	RequestId string
}

var Host = "https://dysmsapi.aliyuncs.com"

func NewSmsParams(accesskey *aliyun.AccessKey) aliyun.Params {
	p := aliyun.NewParams((accesskey))
	p.
		WithAction("SendSms").
		With("RegionId", "cn-hangzhou").
		WithVersion("2017-05-25").
		WithTimestamp(time.Now().Unix()).
		WithNonce()
	return p
}

type Message struct {
	PhoneNumbers    string
	SignName        string
	TemplateCode    string
	TemplateParam   string
	SmsUpExtendCode string
	OutID           string
}

func NewMessage() *Message {
	return &Message{}
}
func Send(accesskey *aliyun.AccessKey, msg *Message) (*Result, error) {
	p := NewSmsParams(accesskey)
	p.With("PhoneNumbers", msg.PhoneNumbers)
	p.With("SignName", msg.SignName)
	p.With("TemplateCode", msg.TemplateCode)
	p.With("TemplateParam", msg.TemplateParam)
	p.With("SmsUpExtendCode", msg.SmsUpExtendCode)
	p.With("OutId", msg.OutID)
	q := p.SignedQuery("GET", accesskey.AccessKeySecret)
	preset := fetcher.NewPreset().With(
		fetcher.SetDoer(&accesskey.Client),
		fetcher.URL(Host),
		fetcher.Params(q),
	)
	var result = &Result{}
	resp, err := preset.FetchAndParse(fetcher.Should200(fetcher.AsJSON(&result)))
	if err != nil {
		return nil, err
	}
	if result.Code != "OK" {
		return nil, resp.NewAPICodeErr(result.Code)
	}
	return result, nil
}
