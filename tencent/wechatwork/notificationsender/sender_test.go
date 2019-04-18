package notificationsender

import (
	"sync"
	"testing"
	"time"

	"github.com/herb-go/notification"
	part "github.com/herb-go/notification/notificationpartmsgpack"
)

var testSender = &Sender{}
var TestGroup = sync.WaitGroup{}

func Test(t *testing.T) {
	var err error
	testSender.Agent.CorpID = TestCorpID
	testSender.Agent.AgentID = TestAgentID
	testSender.Agent.Secret = TestSecret
	notification.DefaultService.RegisterSender(notification.NotificationTypeDefault, testSender)
	notification.DefaultService.SetRecover(func() {
		r := recover()
		if r != nil {
			t.Fatal(r)
		}
	})
	err = notification.DefaultService.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer notification.DefaultService.Stop()
	tests := []func(){
		SendTextMessage,
		SendTextCard,
		SendAttachment,
		SendAttachmentURL,
	}
	for _, v := range tests {
		TestGroup.Add(1)
		go v()
	}
	TestGroup.Wait()
}
func SendAttachmentURL() {
	defer TestGroup.Done()
	ni, err := notification.NewPartedNotificationWithID()
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartAttachmentFilename.Set(ni, "test.png")
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartAttachmentURL.Set(ni, TestPictureURL)
	if err != nil {
		panic(err)
	}
	err = ni.SetNotificationRecipient(TestRecipient)
	if err != nil {
		panic(err)
	}
	notification.Notify(ni)
	time.Sleep(10 * time.Second)

}
func SendAttachment() {
	defer TestGroup.Done()
	ni, err := notification.NewPartedNotificationWithID()
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartAttachmentFilename.Set(ni, "test.txt")
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartAttachment.Set(ni, []byte("test file content"))
	if err != nil {
		panic(err)
	}
	err = ni.SetNotificationRecipient(TestRecipient)
	if err != nil {
		panic(err)
	}
	notification.Notify(ni)
	time.Sleep(5 * time.Second)

}
func SendTextCard() {
	defer TestGroup.Done()
	ni, err := notification.NewPartedNotificationWithID()
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartTitle.Set(ni, "卡片")
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartSummary.Set(ni, "testnewapi")
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartURLTitle.Set(ni, "点击查看")
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartURL.Set(ni, "https://blog.jarlyyn.com")
	if err != nil {
		panic(err)
	}
	err = ni.SetNotificationRecipient(TestRecipient)
	if err != nil {
		panic(err)
	}
	notification.Notify(ni)
	time.Sleep(5 * time.Second)

}
func SendTextMessage() {
	var err error
	defer TestGroup.Done()
	ni, err := notification.NewPartedNotificationWithID()
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	err = part.NotificationPartText.Set(ni, "testnewapi")
	if err != nil {
		panic(err)
	}
	err = ni.SetNotificationRecipient(TestRecipient)
	if err != nil {
		panic(err)
	}
	notification.Notify(ni)
	time.Sleep(5 * time.Second)
}

func init() {
}
