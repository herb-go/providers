//Package redisluacache provides cache driver uses redis to store cache data.
//Using github.com/garyburd/redigo/redis as driver.
package redisluacache

import (
	"sync"
	"time"

	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/model/redis/redispool"
)

var defaultGcPeriod = 30 * time.Second
var defaultGcLimit = int64(100)
var defaultSepartor = string(0)
var tokenMask = cache.TokenMask
var tokenLength = 64
var flushLua = ` 
	if redis.call("HEXISTS",KEYS[2],KEYS[3])==1 then return 0 end
	local v=redis.call("GET",KEYS[1]);
	if (v==false) then v="" end;
    redis.call("HSET",KEYS[2],v,0);
	redis.call("SET",KEYS[1],KEYS[3]);
	return 1;
`
var gcLua = `
	redis.replicate_commands()
	local ks=redis.call("HKEYS",KEYS[1])
	if ks ==false then return end
	local k=ks[1]
	if k ==nil then return end
	local v=redis.call("HGET",KEYS[1],k)
	local r=redis.call("SCAN",v,"MATCH",KEYS[2]..KEYS[3]..KEYS[3]..k..KEYS[3].."*","COUNT",KEYS[4])
	for _,k in ipairs(r[2]) do 
    	redis.call('DEL', k) 
	end
	if r[1]=="0" then redis.call("HDEL",KEYS[1],k) end
`

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
	version, err := c.getVersionFromConn(conn)
	if err != nil {
		return err
	}
	c.version = version
	return nil
}
func (c *Cache) getKey(key string) string {
	c.versionLock.Lock()
	defer c.versionLock.Unlock()
	return c.name + c.Separtor + c.Separtor + c.version + c.Separtor + key
}

func (c *Cache) getVersionKey() string {
	return c.name + c.Separtor + "version" + c.Separtor
}
func (c *Cache) getUsedVersionsKey() string {
	return c.name + c.Separtor + "usedVersions" + c.Separtor

}
func (c *Cache) getVersionFromConn(conn redis.Conn) (string, error) {
	var version string
	vk := c.getVersionKey()
	version, err := redis.String(conn.Do("GET", vk))
	if err == redis.ErrNil {
		version = ""
	} else {
		return version, err
	}
	return version, nil
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *Cache) Flush() error {
	conn := c.Pool.Get()
	defer conn.Close()
	vk := c.getVersionKey()
	version, err := c.getVersionFromConn(conn)
	nv, err := cache.NewRandMaskedBytes(tokenMask, tokenLength, []byte(version))
	if err != nil {
		return err
	}
	vsk := c.getUsedVersionsKey()
	result, err := redis.Int64(conn.Do("EVAL", flushLua, 3, vk, vsk, string(nv)))
	if err != nil {
		return err
	}
	if result == 0 {
		return c.Flush()
	}
	return nil
}
func (c *Cache) gc() error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	vsk := c.getUsedVersionsKey()
	_, err = conn.Do("EVAL", gcLua, 4, vsk, c.name, c.Separtor, c.gcLimit)
	return err
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

func (c *Cache) setVersion(newVersion string) {
	c.versionLock.Lock()
	c.version = newVersion
	c.versionLock.Unlock()

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
	var version string
	var v int64
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	_, err = conn.Do("MULTI")
	if err != nil {
		return v, err
	}
	vk := c.getVersionKey()
	_, err = conn.Do("GET", vk)
	if err != nil {
		return 0, err
	}
	_, err = conn.Do("INCRBY", k, increment)
	if err != nil {
		return v, err
	}

	_, err = conn.Do("EXPIRE", k, int64(ttl/time.Second))
	if err != nil {
		return v, err
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return v, err
	}
	values, err = redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else if err != nil {
		return 0, err
	}
	if version != c.version {
		c.version = version
		_, err = conn.Do("DEL", k)
		if err != nil {
			return 0, err
		}
		return c.IncrCounter(key, increment, ttl)
	}
	_, err = redis.Scan(values, &v)
	return v, err
}
func (c *Cache) doSet(key string, bytes []byte, ttl time.Duration, mode int) error {
	var err error
	var version string
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	_, err = conn.Do("MULTI")
	if err != nil {
		return err
	}
	vk := c.getVersionKey()
	_, err = conn.Do("GET", vk)
	if err != nil {
		return err
	}

	if mode == modeUpdate {
		_, err = conn.Do("SET", k, bytes, "EX", int64(ttl/time.Second), "XX")

	} else {
		_, err = conn.Do("SET", k, bytes, "EX", int64(ttl/time.Second))

	}

	if err != nil {
		return err
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return err
	}
	_, err = redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else if err != nil {
		return err
	}
	if version != c.version {
		c.version = version
		_, err = conn.Do("DEL", k)
		if err != nil {
			return err
		}
		return c.doSet(key, bytes, ttl, mode)
	}
	return nil
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

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	var bytes []byte
	var version string
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	v := c.getVersionKey()
	values, err := redis.Values((conn.Do("MGET", v, k)))
	b, err := redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else {
		if err != nil {
			return bytes, err
		}
	}
	c.versionLock.Lock()
	if version != c.version {
		c.version = version
		c.versionLock.Unlock()
		return c.GetBytesValue(key)
	}
	c.versionLock.Unlock()
	_, err = redis.Scan(b, &bytes)
	if err == redis.ErrNil || bytes == nil {
		return bytes, cache.ErrNotFound
	} else if err != nil {
		return bytes, nil
	}
	return bytes, err
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var data map[string][]byte
	var version string
	var args = make([]interface{}, len(keys)+1)
	conn := c.Pool.Get()
	defer conn.Close()
	v := c.getVersionKey()
	args[0] = v
	for k := range keys {
		args[k+1] = c.getKey(keys[k])
	}
	values, err := redis.Values((conn.Do("MGET", args...)))
	values, err = redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else {
		if err != nil {
			return nil, err
		}
	}
	c.versionLock.Lock()
	if version != c.version {
		c.version = version
		c.versionLock.Unlock()
		return c.MGetBytesValue(keys...)
	}
	c.versionLock.Unlock()
	data = make(map[string][]byte, len(keys))
	for k := range keys {
		var bytes []byte
		values, err = redis.Scan(values, &bytes)
		if err == redis.ErrNil || bytes == nil {
		} else if err != nil {
			return nil, nil
		} else {
			data[keys[k]] = bytes
		}
	}

	return data, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) (err error) {
	var version string
	conn := c.Pool.Get()
	defer conn.Close()
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	vk := c.getVersionKey()
	err = conn.Send("GET", vk)
	if err != nil {
		return err
	}
	var ttlInSecond = int64(ttl / time.Second)
	for k := range data {

		err = conn.Send("SET", c.getKey(k), data[k], "EX", ttlInSecond)

		if err != nil {
			return err
		}
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return err
	}
	_, err = redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else if err != nil {
		return err
	}
	if version != c.version {
		for k := range data {
			err = conn.Send("DEL", c.getKey(k))
			if err != nil {
				return err
			}
		}
		c.version = version
		return c.MSetBytesValue(data, ttl)
	}
	return nil
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
	c.gcErrHandler = f
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
	period := time.Duration(c.GCPeriod)
	if period == 0 {
		period = defaultGcPeriod
	}
	cache.ticker = time.NewTicker(period)
	gcLimit := c.GCLimit
	if gcLimit == 0 {
		gcLimit = defaultGcLimit
	}
	cache.gcLimit = gcLimit
	go func() {
		for {
			select {
			case <-cache.ticker.C:
				err := cache.gc()
				if err != nil {
					if cache.gcErrHandler != nil {
						cache.gcErrHandler(err)
					}
				}
			case <-cache.quit:
				cache.ticker.Stop()
				return
			}
		}

	}()
	err := cache.start()
	if err != nil {
		return &cache, err
	}
	return &cache, nil
}
func init() {
	cache.Register("redisluacache", func(loader func(interface{}) error) (cache.Driver, error) {
		var err error
		c := &Config{}
		err = loader(c)
		if err != nil {
			return nil, err
		}
		return c.Create()
	})
}
