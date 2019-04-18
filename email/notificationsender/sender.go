package notificationsender

import (
	"net/smtp"
	"strconv"

	"github.com/herb-go/notification"
	part "github.com/herb-go/notification/notificationpartmsgpack"
	"github.com/jordan-wright/email"
)

// RequiredFields reqired fields for email notification sender
var RequiredFields = [][]string{
	[]string{"title"},
	[]string{"summary", "text", "html"},
}

//Sender email sender struct
type Sender struct {
	// SenderName name for sender interface
	SenderName string
	// Host smtp host addr
	Host string
	// Prot smtp port addr
	Port int
	// Identity user identity(user account) for stmp arddr
	Identity string
	// Sender sender email name
	Sender string
	// email from address
	From string
	// Username email stmp user name
	Username string
	// Pasword email smtp password
	Password string
	// prefix for subject
	SubjectPrefix string
}

// Name name for sender
func (s *Sender) Name() string {
	return s.SenderName
}

//SendNotification send notification Instance
func (s *Sender) SendNotification(i *notification.NotificationInstance) error {
	n, err := notification.ValidatePartedNotificationInstanceWithFields(i, RequiredFields)
	if err != nil {
		return i.NewError(err)
	}
	if n == nil {
		return nil
	}
	msg := email.NewEmail()
	msg.Sender = s.Sender
	if s.From != "" {
		msg.From = s.From
	}
	msg.Sender = s.Sender
	title, err := part.NotificationPartTitle.Get(n)
	if err != nil {
		return i.NewError(err)
	}
	msg.Subject = s.SubjectPrefix + title
	text, err := part.NotificationPartText.Get(n)
	if err != nil {
		return i.NewError(err)
	}
	if text != "" {
		msg.Text = []byte(text)
	} else {
		text, err := part.NotificationPartSummary.Get(n)
		if err != nil {
			return i.NewError(err)
		}
		if text != "" {
			msg.Text = []byte(text)
		}
	}
	html, err := part.NotificationPartHtml.Get(n)
	if err != nil {
		return i.NewError(err)
	}
	if html != "" {
		msg.HTML = []byte(html)
	}
	msg.To = []string{i.Recipient}
	err = msg.Send(s.Host+":"+strconv.Itoa(s.Port), smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host))
	if err != nil {
		return i.NewError(err)
	}
	i.SetStatusSuccess()
	return nil
}
