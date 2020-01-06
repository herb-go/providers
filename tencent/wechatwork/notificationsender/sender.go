package notificationsender

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/herb-go/notification"
	part "github.com/herb-go/notification/notificationpartmsgpack"
	"github.com/herb-go/providers/tencent/wechatwork"
)

type Sender struct {
	SenderName string
	Agent      wechatwork.Agent
}

func (s *Sender) Name() string {
	return s.SenderName
}

type NotificationData struct {
	AttachmentMediaID  string
	AttachmentType     string
	Attachment         []byte
	AttachmentFilename string
	Text               string
	Summary            string
	Title              string
	URL                string
	URLTitle           string
}
type messageInitFunc func(message *wechatwork.Message, data *NotificationData) (bool, error)

func (s *Sender) InitMessage(message *wechatwork.Message, data *NotificationData, inits ...messageInitFunc) (bool, error) {
	for _, v := range inits {
		result, err := v(message, data)
		if err != nil || result == true {
			return result, err
		}
	}
	return false, nil
}
func (s *Sender) InitTextCardMessage(message *wechatwork.Message, data *NotificationData) (bool, error) {
	if data.URL == "" && (data.Title == "" || data.Summary == "") {
		return false, nil
	}
	message.SetMsgType("textcard")
	message.TextCard = &wechatwork.BodyMessageTextCard{
		Title:       data.Title,
		URL:         data.URL,
		Description: data.Summary,
	}
	if data.URLTitle != "" {
		message.TextCard.Btntxt = &data.URLTitle
	}
	return true, nil
}
func (s *Sender) InitTextMessage(message *wechatwork.Message, data *NotificationData) (bool, error) {
	if data.Text == "" && data.Summary == "" && data.Title == "" {
		return false, nil
	}
	message.SetMsgType("text")
	message.Text = &wechatwork.MessageText{}
	if data.Text != "" {
		message.Text.Content = data.Text
	} else if data.Summary != "" {
		message.Text.Content = data.Summary
	} else {
		message.Text.Content = data.Title
	}
	return true, nil
}

func (s *Sender) InitMediaMessage(message *wechatwork.Message, data *NotificationData) (bool, error) {
	if data.AttachmentType == "" || data.AttachmentMediaID == "" {
		return false, nil
	}
	message.SetMsgType(data.AttachmentType)
	media := &wechatwork.MessageMedia{
		MediaID: data.AttachmentMediaID,
	}
	switch data.AttachmentType {
	case "image":
		message.Image = media
	case "voice":
		message.Voice = media
	case "video":
		message.Video = &wechatwork.MessageVideo{
			MediaID: data.AttachmentMediaID,
		}
		if data.Title != "" {
			message.Video.Title = &data.Title
		}
		if data.Summary != "" {
			message.Video.Description = &data.Summary
		}
	case "file":
		message.File = media
	}
	return true, nil
}
func (s *Sender) GetData(
	n *notification.PartedNotification,
	data *NotificationData,
	getters ...func(n *notification.PartedNotification, data *NotificationData) error) error {
	for _, v := range getters {
		err := v(n, data)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Sender) GetURLData(n *notification.PartedNotification, data *NotificationData) error {
	var err error
	data.URL, err = part.NotificationPartURL.Get(n)
	if err != nil {
		return err
	}
	data.URLTitle, err = part.NotificationPartURLTitle.Get(n)
	if err != nil {
		return err
	}
	return nil
}
func (s *Sender) GetTextData(n *notification.PartedNotification, data *NotificationData) error {
	var err error
	data.Text, err = part.NotificationPartText.Get(n)
	if err != nil {
		return err
	}
	data.Title, err = part.NotificationPartTitle.Get(n)
	if err != nil {
		return err
	}
	data.Summary, err = part.NotificationPartSummary.Get(n)
	if err != nil {
		return err
	}
	return nil
}
func (s *Sender) GetMediaData(n *notification.PartedNotification, data *NotificationData) error {
	var err error
	data.AttachmentMediaID, err = NotificationPartAttachmentMediaID.Get(n)
	if err != nil {
		return err
	}
	data.AttachmentType, err = NotificationPartAttachmentMediaType.Get(n)
	if err != nil {
		return err
	}
	if data.AttachmentMediaID == "" || data.AttachmentType == "" {
		data.Attachment, err = part.NotificationPartAttachment.Get(n)
		if err != nil {
			return err
		}
		if len(data.Attachment) == 0 {
			url, err := part.NotificationPartAttachmentURL.Get(n)
			if err != nil {
				return err
			}
			if url != "" {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return err
				}
				resp, err := s.Agent.Client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				data.Attachment, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
			}
		}
		data.AttachmentFilename, err = part.NotificationPartAttachmentFilename.Get(n)
		if err != nil {
			return err
		}
		if len(data.Attachment) > 0 && data.AttachmentFilename != "" {
			switch strings.ToLower(path.Ext(data.AttachmentFilename)) {
			case ".png", ".jpg", ".jpeg":
				data.AttachmentType = "image"
			case ".amr":
				data.AttachmentType = "voice"
			case ".mp4":
				data.AttachmentType = "video"
			default:
				data.AttachmentType = "file"
			}
			data.AttachmentMediaID, err = s.Agent.MediaUpload(data.AttachmentType, data.AttachmentFilename, bytes.NewBuffer(data.Attachment))
			if err != nil {
				return err
			}
			err = NotificationPartAttachmentMediaID.Set(n, data.AttachmentMediaID)
			if err != nil {
				return err
			}
			err = NotificationPartAttachmentMediaType.Set(n, data.AttachmentType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Sender) SendNotification(i *notification.NotificationInstance) error {
	var data = &NotificationData{}
	n, err := notification.ValidatePartedNotificationInstanceWithFields(i, nil)
	if err != nil {
		return i.NewError(err)
	}
	if n == nil {
		return nil
	}
	i.Notification.LockNotification()
	defer i.Notification.UnlockNotification()
	err = s.GetData(n, data,
		s.GetMediaData,
		s.GetTextData,
		s.GetURLData,
	)
	if err != nil {
		return i.NewError(err)
	}
	message := s.Agent.NewMessage()
	message.SetToUser(i.Recipient)
	ok, err := s.InitMessage(message, data,
		s.InitTextCardMessage,
		s.InitMediaMessage,
		s.InitTextMessage,
	)
	if err != nil {
		return i.NewError(err)
	}
	if ok == false {
		i.SetStatusUnsupported()
		return nil
	}
	result, err := s.Agent.SendMessage(message)
	if err != nil {
		return i.NewError(err)
	}
	if result.Errcode != 0 {
		return i.NewError(wechatwork.NewResultError(result.Errcode, result.Errmsg))
	}
	i.SetStatusSuccess()
	return nil
}
