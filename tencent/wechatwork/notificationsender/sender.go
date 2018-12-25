package notificationsender

import (
	"github.com/herb-go/herb/notification"
	part "github.com/herb-go/herb/notification/notificationpartmsgpack"
	"github.com/herb-go/providers/tencent/wechatwork"
)

var RequiredFields = [][]string{
	[]string{"summary", "text", "html"},
}

type Sender struct {
	SenderName string
	Agent      wechatwork.Agent
}

func (s *Sender) Name() string {
	return s.SenderName
}

func (s *Sender) SendNotification(i *notification.NotificationInstance) error {
	n, err := notification.ValidatePartedNotificationInstanceWithFields(i, RequiredFields)
	if err != nil {
		return i.NewError(err)
	}
	if n == nil {
		return nil
	}
	message := s.Agent.NewMessage()
	text, err := part.NotificationPartText.Get(n)
	if err != nil {
		return i.NewError(err)
	}
	message.SetToUser(i.Recipient)
	message.SetMsgType("text")
	message.Text.Content = text
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
