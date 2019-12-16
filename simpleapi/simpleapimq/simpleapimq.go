package simpleapimq

import (
	"io/ioutil"
	"net/http"

	"github.com/herb-go/providers/simpleapi"

	"github.com/herb-go/herb/middleware/misc"

	"github.com/herb-go/herb/middleware"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/herb-go/herb/service/httpservice/target"
	"github.com/herb-go/messagequeue"
)

type Config struct {
}
type Broke struct {
	Server   simpleapi.ServerConfig
	Client   *target.PlainPlan
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
func (b *Broke) register() error {
	app := middleware.New()
	app.Use(misc.MethodMiddleware("POST"))
	app.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
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
	return b.Server.Channel.Handle(app)
}

// Listen listen queue
//Return any error if raised
func (b *Broke) Listen() error {
	return b.Server.Channel.Start()
}

//Close close queue
//Return any error if raised
func (b *Broke) Close() error {
	return b.Server.Channel.Stop()
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
	resp, err := target.Do(b.Client.Doer, b.Client.Target, target.Body(data))
	if resp.StatusCode != 200 {
		return nil, target.NewError(resp)
	}
	return sent, nil
}

//SetConsumer set message consumer
func (b *Broke) SetConsumer(c func(*messagequeue.Message) messagequeue.ConsumerStatus) {
	b.consumer = c
}
