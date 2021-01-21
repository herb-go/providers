package tencentminiprogramum

import (
	"encoding/json"
	"net/url"

	"github.com/herb-go/fetcher"
	"github.com/herb-go/providers/tencent/tencentminiprogram"
)

type ResultAPIError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type TemplateMessageMiniprogram struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

type WeappTemplateMessage struct {
	TemplateID      string `json:"template_id"`
	Page            string `json:"page"`
	FormID          string `json:"form_id"`
	EmphasisKeyword string `json:"emphasis_keyword"`
	Data            string `json:"data"`
}
type TemplateMessage struct {
	AppID       string                      `json:"appid"`
	TemlpateID  string                      `json:"template_id"`
	Miniprogram *TemplateMessageMiniprogram `json:"miniprogram"`
	URL         *string                     `json:"url"`
	Data        json.RawMessage             `json:"data"`
}

var APIUniformSend = tencentminiprogram.Server.EndPoint("POST", "cgi-bin/message/wxopen/template/uniform_send")

type Message struct {
	ToUser               string                `json:"touser"`
	WeappTemplateMessage *WeappTemplateMessage `json:"weapp_template_msg"`
	MpTemplateMsg        *TemplateMessage      `json:"mp_template_msg"`
}

func Send(app *tencentminiprogram.App, msg *Message) error {
	token, err := app.GetAccessToken()
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("access_token", token)
	result := &ResultAPIError{}
	preset := APIUniformSend.With(
		&app.Client,
		fetcher.Params(params),
		fetcher.JSONBody(msg),
	)
	resp, err := preset.FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if result.Errcode != 0 {
		return resp.NewAPICodeErr(result.Errcode)
	}
	return nil
}
