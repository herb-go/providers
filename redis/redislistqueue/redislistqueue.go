package redislistqueue

import (
	"github.com/gomodule/redigo/redis"
	"github.com/herb-go/herb/model/redis/redispool"
	"github.com/herb-go/messagequeue"
)

type Queue struct {
	*redispool.Config
	pool         *redispool.Pool
	producerpool *redispool.Pool
	Topic        string
	Timeout      int
	consumer     func(*messagequeue.Message) messagequeue.ConsumerStatus
	recover      func()
}

func (q *Queue) SetRecover(r func()) {
	q.recover = r
}

func (q *Queue) brpop() {
	conn := q.pool.Get()
	defer conn.Close()
	r, err := redis.ByteSlices(conn.Do("BRPOP", q.Topic, q.Timeout))
	if err == redis.ErrNil {
		return
	}
	if err != nil {
		panic(err)
	}
	q.consumer(messagequeue.NewMessage(r[1]))
}
func (q *Queue) pull() {
	defer q.recover()
	for {
		q.brpop()
	}
}

//Connect to brocker as producer
func (q *Queue) Connect() error {
	q.producerpool = redispool.New()
	err := q.Config.ApplyTo(q.producerpool)
	if err != nil {
		return err
	}
	q.producerpool.Open()
	return nil
}

//Disconnect stop producing and disconnect
func (q *Queue) Disconnect() error {
	return q.producerpool.Close()
}

func (q *Queue) Listen() error {
	q.pool = redispool.New()
	err := q.Config.ApplyTo(q.pool)
	if err != nil {
		return err
	}
	q.pool.Open()

	go q.pull()
	return nil
}
func (q *Queue) Close() error {
	return q.pool.Close()
}
func (q *Queue) ProduceMessages(messages ...[]byte) (sent []bool, err error) {
	sent = make([]bool, len(messages))
	conn := q.producerpool.Get()
	defer conn.Close()
	for k := range messages {
		_, err := conn.Do("LPUSH", q.Topic, messages[k])
		if err != nil {
			return sent, err
		}
		sent[k] = true
	}
	return sent, nil
}
func (q *Queue) SetConsumer(c func(*messagequeue.Message) messagequeue.ConsumerStatus) {
	q.consumer = c
}

func NewQueue() *Queue {
	return &Queue{
		recover: func() {},
		Config:  redispool.NewConfig(),
	}
}

func QueueFactory(loader func(interface{}) error) (messagequeue.Driver, error) {
	q := NewQueue()
	var err error
	err = loader(q)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func init() {
	messagequeue.Register("redislist", QueueFactory)
}
