package wechatwork

import "strings"

const MsgTypeText = "text"
const MsgTypeImage = "image"
const MsgTypeVoice = "voice"
const MsgTypeVideo = "video"
const MsgTypeFile = "file"
const MsgTypeNews = "news"
const MsgTypeTextcard = "textcard"
const MsgTypeMarkdown = "markdown"
const MsgTypeMPNews = "mpnews"
const MsgTypeTaskcard = "taskcard"

type Message struct {
	ToUser   *string          `json:"touser"`
	ToParty  *string          `json:"toparty"`
	ToTag    *string          `json:"totag"`
	MsgType  string           `json:"msgtype"`
	AgentID  int              `json:"agentid"`
	Safe     int              `json:"safe"`
	Text     *MessageText     `json:"text"`
	Image    *MessageMedia    `json:"image"`
	Voice    *MessageMedia    `json:"voice"`
	File     *MessageMedia    `json:"file"`
	Video    *MessageVideo    `json:"video"`
	News     *MessageNews     `json:"news"`
	MPNews   *MessageMPNews   `json:"mpnews"`
	Textcard *MessageTextcard `json:"textcard"`
	Taskcard *MessageTaskcard `json:"taskcard"`
	Markdown *MessageMarkdown `json:"markdown"`
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
	case MsgTypeText:
		p.Text = &MessageText{}
	case MsgTypeImage:
		p.Image = &MessageMedia{}
	case MsgTypeVoice:
		p.Voice = &MessageMedia{}
	case MsgTypeFile:
		p.File = &MessageMedia{}
	case MsgTypeNews:
		p.News = &MessageNews{}
	case MsgTypeVideo:
		p.Video = &MessageVideo{}
	case MsgTypeTextcard:
		p.Textcard = &MessageTextcard{}
	case MsgTypeTaskcard:
		p.Taskcard = &MessageTaskcard{}
	case MsgTypeMarkdown:
		p.Markdown = &MessageMarkdown{}
	case MsgTypeMPNews:
		p.MPNews = &MessageMPNews{}
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

type MessageMPNews struct {
	Articles []*MessageMPArticle `json:"articles"`
}
type MessageMPArticle struct {
	Title            string  `json:"title"`
	ThumbMediaID     string  `json:"thumb_media_id"`
	ContentSourceUrl *string `json:"content_source_url"`
	Content          *string `json:"content"`
	Digest           *string `json:"digest"`
}

type MessageTextcard struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	Btntxt      *string `json:"btntxt"`
}

type MessageTaskcard struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	URL         *string               `json:"url"`
	TaskID      string                `json:"task_id"`
	Btn         []*MessageTaskcardBtn `json:"btn"`
}
type MessageTaskcardBtn struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	ReplaceName string `json:"replace_name"`
	Color       string `json:"color"`
	IsBold      bool   `json:"is_bold"`
}
type MessageMarkdown struct {
	Content string `json:"content"`
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
