package redislistqueue

import (
	"github.com/gomodule/redigo/redis"
	"github.com/herb-go/herb/model/redis/redispool"
	"github.com/herb-go/messagequeue"
)

type Queue struct {
	*redispool.Config
	pool     *redispool.Pool
	Topic    string
	Timeout  int
	consumer func(*messagequeue.Message) messagequeue.ConsumerStatus
	recover  func()
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
func (q *Queue) Start() error {
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
	conn := q.pool.Get()
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

func QueueFactory(conf messagequeue.Config, prefix string) (messagequeue.Driver, error) {
	q := NewQueue()
	var err error
	err = conf.Get(prefix+"Name", &q.Topic)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"Timeout", &q.Timeout)
	if err != nil {
		return nil, err
	}
	if q.Timeout == 0 {
		q.Timeout = 30
	}
	err = conf.Get(prefix+"Network", &q.Config.Network)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"Address", &q.Config.Address)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"Password", &q.Config.Password)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"Db", &q.Config.Db)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"MaxIdle", &q.Config.MaxIdle)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"MaxAlive", &q.Config.MaxAlive)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"ConnectTimeoutInSecond", &q.Config.ConnectTimeoutInSecond)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"ReadTimeoutInSecond", &q.Config.ReadTimeoutInSecond)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"WriteTimeoutInSecond", &q.Config.WriteTimeoutInSecond)
	if err != nil {
		return nil, err
	}

	err = conf.Get(prefix+"IdleTimeoutInSecond", &q.Config.IdleTimeoutInSecond)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func init() {
	messagequeue.Register("redislist", QueueFactory)
}
