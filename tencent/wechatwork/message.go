package wechatwork

import "strings"

type Message struct {
	ToUser   *string              `json:"touser"`
	ToParty  *string              `json:"toparty"`
	ToTag    *string              `json:"totag"`
	MsgType  string               `json:"msgtype"`
	AgentID  int                  `json:"agentid"`
	Safe     int                  `json:"safe"`
	Text     *MessageText         `json:"text"`
	Image    *MessageMedia        `json:"image"`
	Voice    *MessageMedia        `json:"voice"`
	File     *MessageMedia        `json:"file"`
	Video    *MessageVideo        `json:"video"`
	TextCard *BodyMessageTextCard `json:"textcard"`
}

func (p *Message) SetToUser(users ...string) {
	to := strings.Join(users, "|")
	p.ToUser = &to
}
func (p *Message) SetToParty(parties ...string) {
	to := strings.Join(parties, "|")
	p.ToParty = &to
}
func (p *Message) SetToTag(tags ...string) {
	to := strings.Join(tags, "|")
	p.ToTag = &to
}
func (p *Message) SetMsgType(MsgType string) {
	p.MsgType = MsgType
	switch MsgType {
	case "text":
		p.Text = &MessageText{}
	case "image":
		p.Image = &MessageMedia{}
	case "voice":
		p.Voice = &MessageMedia{}
	case "file":
		p.File = &MessageMedia{}
	case "video":
		p.Video = &MessageVideo{}
	case "textcard":
		p.TextCard = &BodyMessageTextCard{}
	}
}

type MessageText struct {
	Content string `json:"content"`
}
type MessageMedia struct {
	MediaID string `json:"media_id"`
}
type MessageVideo struct {
	MediaID     string  `json:"media_id"`
	Title       *string `json:""`
	Description *string `json:"description"`
}
type BodyMessageTextCard struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	Btntxt      *string `json:"btntxt"`
}
type MessageResult struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}
