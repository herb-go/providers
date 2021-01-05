package wechatwork

import "strings"

const MsgTypeText = "text"
const MsgTypeImage = "image"
const MsgTypeVoice = "voice"
const MsgTypeVideo = "video"
const MsgTypeFile = "file"
const MsgTypeNews = "news"

type Message struct {
	ToUser  *string       `json:"touser"`
	ToParty *string       `json:"toparty"`
	ToTag   *string       `json:"totag"`
	MsgType string        `json:"msgtype"`
	AgentID int           `json:"agentid"`
	Safe    int           `json:"safe"`
	Text    *MessageText  `json:"text"`
	Image   *MessageMedia `json:"image"`
	Voice   *MessageMedia `json:"voice"`
	File    *MessageMedia `json:"file"`
	Video   *MessageVideo `json:"video"`
	News    *MessageNews  `json:"news"`
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
	case "news":
		p.News = &MessageNews{}
	case "video":
		p.Video = &MessageVideo{}
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
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type MessageNews struct {
	Articles []*MessageArticle `json:"articles"`
}
type MessageArticle struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	URL         *string `json:"url"`
	PicURL      *string `json:"picurl"`
}
type MessageResult struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

func NewMessage() *Message {
	return &Message{}
}
