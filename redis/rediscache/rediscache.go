//Package rediscache provides cache driver uses redis to store cache data.
//Using github.com/garyburd/redigo/redis as driver.
package rediscache

import (
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/model/redis/redispool"
)

var defaultSepartor = string(0)

const modeSet = 0
const modeUpdate = 1

//Cache The redis cache Driver.
type Cache struct {
	cache.DriverUtil
	Pool           *redis.Pool //Redis pool.
	ticker         *time.Ticker
	name           string
	quit           chan int
	gcErrHandler   func(err error)
	gcLimit        int64
	network        string
	address        string
	password       string
	version        string
	versionLock    sync.Mutex
	db             int
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	Separtor       string //Separtor in redis key.
}

func (c *Cache) start() error {
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}
func (c *Cache) getKey(key string) string {
	return c.name + c.Separtor + key
}

//Flush Flush not supported.
func (c *Cache) Flush() error {
	return cache.ErrFeatureNotSupported
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	return c.Pool.Close()
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	k := c.getKey(key)
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	return err
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	val := strconv.FormatInt(v, 10)
	return c.SetBytesValue(key, []byte(val), ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return v, err
	}
	return strconv.ParseInt(string(bytes), 10, 64)
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	k := c.getKey(key)
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	return err
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var err error
	var v int64
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)

	v, err = redis.Int64(conn.Do("INCRBY", k, increment))
	if err != nil {
		return v, err
	}

	_, err = conn.Do("EXPIRE", k, int64(ttl/time.Second))
	if err != nil {
		return v, err
	}

	return v, err
}
func (c *Cache) doSet(key string, bytes []byte, ttl time.Duration, mode int) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)

	if mode == modeUpdate {
		_, err = conn.Do("SET", k, bytes, "EX", int64(ttl/time.Second), "XX")

	} else {
		_, err = conn.Do("SET", k, bytes, "EX", int64(ttl/time.Second))

	}
	return err
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	return c.doSet(key, bytes, ttl, modeSet)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	return c.doSet(key, bytes, ttl, modeUpdate)
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var data = make(map[string][]byte, len(keys))
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	for key := range keys {
		k := c.getKey(keys[key])
		err := (conn.Send("GET", k))
		if err != nil {
			return nil, err
		}
	}

	err = conn.Flush()
	if err != nil {
		return nil, err
	}
	for key := range keys {
		bs, err := redis.Bytes((conn.Receive()))
		if err == redis.ErrNil {
			continue
		}
		if err != nil {
			return nil, err
		}
		data[keys[key]] = bs
	}

	return data, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	var ttlInSecond = int64(ttl / time.Second)
	for key := range data {
		k := c.getKey(key)

		err = conn.Send("SET", k, data[key], "EX", ttlInSecond)

		if err != nil {
			return err
		}
	}
	return conn.Flush()

}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	var bs []byte
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	bs, err := redis.Bytes((conn.Do("GET", k)))
	if err == redis.ErrNil {
		return nil, cache.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if bs == nil {
		return nil, cache.ErrNotFound
	}
	return bs, err
}

//Expire set cache value expire duration by given key and ttl
func (c *Cache) Expire(key string, ttl time.Duration) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)

	_, err = conn.Do("EXPIRE", k, int64(ttl/time.Second))

	return err
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)

	_, err = conn.Do("EXPIRE", k, int64(ttl/time.Second))

	return err
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	return
}

//Config Cache driver config.
type Config struct {
	redispool.Config
	GCPeriod int64 //Period of gc.Default value is 30 second.
	GCLimit  int64 //Max delete limit in every gc call.Default value is 100.
}

//Create create new cache driver.
//Return driver created and any error if raised.
func (c *Config) Create() (cache.Driver, error) {

	cache := Cache{}
	p := redispool.New()
	c.Config.ApplyTo(p)
	cache.Pool = p.Open()
	cache.quit = make(chan int)
	err := cache.start()
	if err != nil {
		return &cache, err
	}
	return &cache, nil
}
func init() {
	cache.Register("rediscache", func(loader func(interface{}) error) (cache.Driver, error) {
		var err error
		c := &Config{}
		err = loader(c)
		if err != nil {
			return nil, err
		}
		return c.Create()
	})
}
