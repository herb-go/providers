package redislistqueue_test

import (
	"bytes"
	"container/list"
	"testing"
	"time"

	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"

	"github.com/herb-go/messagequeue"
)

func newTestBroker() *messagequeue.Broker {
	b := messagequeue.NewBroker()
	c := messagequeue.NewOptionConfig()
	err := loader.LoadConfig("json", []byte(testConfig), c)
	if err != nil {
		panic(err)
	}
	err = c.ApplyTo(b)
	if err != nil {
		panic(err)
	}
	return b
}

func TestBroker(t *testing.T) {
	b := newTestBroker()
	err := b.Listen()
	if err != nil {
		t.Fatal(err)
	}
	err = b.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := b.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = b.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()
	testchan := make(chan *messagequeue.Message, 100)
	b.SetConsumer(messagequeue.NewChanConsumer(testchan))
	messages := [][]byte{}
	unreceived := list.New()
	for i := byte(0); i < 5; i++ {
		messages = append(messages, []byte{i})
		unreceived.PushBack([]byte{i})
	}
	sent, err := b.ProduceMessages(messages...)
	if err != nil {
		t.Fatal(err)
	}
	for k := range sent {
		if sent[k] == false {
			t.Fatal(k)
		}
	}
	time.Sleep(1 * time.Second)
	if len(testchan) != 5 {
		t.Fatal(len(testchan))
	}
	if unreceived.Len() != 5 {
		t.Fatal(unreceived.Len())
	}
	for i := byte(0); i < 5; i++ {
		m := <-testchan
		e := unreceived.Front()
		for e != nil {
			if bytes.Compare(e.Value.([]byte), m.Data) == 0 {
				unreceived.Remove(e)
				break
			}
			e = e.Next()
		}
	}
	if unreceived.Len() != 0 {
		t.Fatal(unreceived)
	}
}
