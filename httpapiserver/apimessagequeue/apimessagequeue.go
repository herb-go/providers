package apimessagequeue

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/herb-go/messagequeue"
	"github.com/herb-go/providers/httpapiserver"
)

type Config struct {
}
type Broke struct {
	Server   httpapiserver.Server
	Client   httpapiserver.Client
	recover  func()
	consumer func(*messagequeue.Message) messagequeue.ConsumerStatus
}

//Connect to brocker as producer
//Return any error if raised
func (b *Broke) Connect() error {
	return nil
}

//Disconnect stop producing and disconnect
//Return any error if raised
func (b *Broke) Disconnect() error {
	return nil
}

// Listen listen queue
//Return any error if raised
func (b *Broke) Listen() error {
	return b.Server.StartWithGuarder(func(w http.ResponseWriter, r *http.Request) {
		ms := []*messagequeue.Message{}
		body, err := r.GetBody()
		if err != nil {
			panic(err)
		}
		defer body.Close()
		data, err := ioutil.ReadAll(body)
		if err != nil {
			panic(err)
		}
		err = msgpack.Unmarshal(data, &ms)
		if err != nil {
			panic(err)
		}
		for k := range ms {
			_ = b.consumer(ms[k])
		}
	})
}

//Close close queue
//Return any error if raised
func (b *Broke) Close() error {
	return b.Server.Stop()
}

//SetRecover set recover
func (b *Broke) SetRecover(r func()) {
	b.recover = r
}

// ProduceMessages produce messages to broke
//Return sent result and any error if raised
func (b *Broke) ProduceMessages(bs ...[]byte) (sent []bool, err error) {
	ms := make([]*messagequeue.Message, len(bs))
	sent = make([]bool, len(bs))
	for k := range bs {
		ms[k] = messagequeue.NewMessage(bs[k])
		sent[k] = true
	}
	data, err := msgpack.Marshal(ms)
	if err != nil {
		return nil, err
	}
	resp, err := b.Client.FetchRequest(nil, data)
	if resp.StatusCode != 200 {
		return nil, resp
	}
	return sent, nil
}

//SetConsumer set message consumer
func (b *Broke) SetConsumer(c func(*messagequeue.Message) messagequeue.ConsumerStatus) {
	b.consumer = c
}
