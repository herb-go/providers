package apimessagequeue

import (
	"github.com/herb-go/fetch"
	"github.com/herb-go/herb/service/httpservice/apiserver"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type Broke struct {
	Guarder   guarder.Guarder
	Vistor    guarder.Visitor
	Server    fetch.Server
	Clients   fetch.Clients
	APIServer apiserver.Option
	recover   func()
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
	return b.APIServer.Server()
}

//Close close queue
//Return any error if raised
func (b *Broke) Close() error {

}

//SetRecover set recover
func (b *Broke) SetRecover(r func()) {
	b.recover = r
}

// ProduceMessages produce messages to broke
//Return sent result and any error if raised
func (b *Broke) ProduceMessages(...[]byte) (sent []bool, err error) {

}

//SetConsumer set message consumer
func (b *Broke) SetConsumer(func(*Message) ConsumerStatus) {

}
