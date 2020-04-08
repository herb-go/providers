//Package sqlcache provides cache driver uses sqlite or mysql to store cache data.
//Using database/sql as driver.
//You should create data table with sql file in "sql" folder first.
package sqlcache

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/herb-go/herb/model/sql/db"

	"github.com/herb-go/herb/cache"
)

const modelSet = 0
const modelUpdate = 1

var defaultGCPeriod = 5 * time.Minute
var tokenMask = cache.TokenMask
var defaultGcLimit = int64(100)

//Cache The sql cache Driver.
type Cache struct {
	cache.DriverUtil
	DB           *sql.DB
	table        string
	name         string
	ticker       *time.Ticker
	quit         chan int
	gcErrHandler func(err error)
	gcLimit      int64
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	c.gcErrHandler = f
	return
}
func (c *Cache) start() error {
	err := c.gc()
	return err
}
func (c *Cache) getVersionTx(tx *sql.Tx) ([]byte, error) {
	var version []byte
	stmt, err := tx.Prepare(`Select version from ` + c.table + ` WHERE cache_key="" AND cache_name = ?`)
	if err != nil {
		return version, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(c.name)
	err = row.Scan(&version)
	if err == sql.ErrNoRows {
		return []byte{}, nil
	}
	return version, err
}

func (c *Cache) gc() error {
	var keys []string

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}

	stmtExpired, err := tx.Prepare(`Select cache_key FROM ` + c.table + ` Where cache_name = ? AND expired > -1  AND expired < ? limit ?`)
	if err != nil {
		return err
	}
	defer stmtExpired.Close()

	rows, err := stmtExpired.Query(c.name, time.Now().Unix(), c.gcLimit)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}
	stmtVersionWrong, err := tx.Prepare(`Select cache_key FROM ` + c.table + ` Where cache_name = ? AND version != ? limit ?`)
	if err != nil {
		return err
	}
	defer stmtVersionWrong.Close()
	rows2, err := stmtVersionWrong.Query(c.name, version, c.gcLimit)
	if err != nil {
		return err
	}
	defer rows2.Close()
	for rows2.Next() {
		var key string
		err = rows2.Scan(&key)
		if err != nil {
			return err
		}
		keys = append(keys, key)

	}

	stmt, err := tx.Prepare(`DELETE FROM ` + c.table + ` Where cache_name=? and cache_key = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range keys {
		_, err = stmt.Exec(c.name, v)
		if err != nil {
			return err
		}
	}
	if err == nil {
		tx.Commit()
	}
	return err
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var v int64
	tx, err := c.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(`Select cache_value from ` + c.table + ` WHERE  expired > ? AND cache_name =? AND cache_key = ? AND version=?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	r := stmt.QueryRow(time.Now().Unix(), c.name, key, version)
	var j string
	err = r.Scan(&j)
	if err == sql.ErrNoRows {
		v = 0
	} else if err != nil {
		return 0, err
	} else {
		err = json.Unmarshal([]byte(j), &v)
		if err != nil {
			v = 0
		}
	}
	v = v + increment
	val, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	var expired int64

	expired = time.Now().Add(ttl).Unix()

	stmtset, err := tx.Prepare(`update ` + c.table + ` set
	 cache_value=?,
	 version=?,
	 expired=?
	 Where cache_name=? 
	 and cache_key=?
	 `)

	defer stmtset.Close()
	row, err := stmtset.Exec(
		val,
		version,
		expired,
		c.name,
		key)
	if err != nil {
		return 0, err
	}
	affected, err := row.RowsAffected()
	if err != nil {
		return v, err
	}
	if affected == 0 {
		stmt2, err := tx.Prepare(`insert into ` + c.table + ` (cache_name,cache_key,cache_value,version,expired) values (?,?,?,?,?)`)
		if err != nil {
			return v, err
		}
		defer stmt2.Close()
		_, err = stmt2.Exec(c.name, key, string(val), version, expired)
	}
	if err != nil {
		return v, err
	}
	tx.Commit()
	return v, nil
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	bs, err := c.Util().Marshaler.Marshal(v)
	if err != nil {
		return nil
	}
	return c.SetBytesValue(key, bs, ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
	bs, err := c.GetBytesValue(key)
	if err != nil {
		return 0, err
	}
	err = c.Util().Marshaler.Unmarshal(bs, &v)
	return v, err
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	return c.Del(key)
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bs []byte, ttl time.Duration) error {
	return c.doSet(key, bs, ttl, modelSet)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bs []byte, ttl time.Duration) error {
	return c.doSet(key, bs, ttl, modelUpdate)
}
func (c *Cache) doSet(key string, bs []byte, ttl time.Duration, mode int) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	var expired int64

	expired = time.Now().Add(ttl).Unix()

	stmt, err := tx.Prepare(`update ` + c.table + ` set
	 cache_value=?,
	 version=?,
	 expired=?
	 Where cache_name=? 
	 and cache_key=?
	 `)

	defer stmt.Close()
	r, err := stmt.Exec(
		bs,
		version,
		expired,
		c.name,
		key)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 && mode != modelUpdate {
		stmt2, err := tx.Prepare(`insert into ` + c.table + ` (cache_name,cache_key,cache_value,version,expired) values (?,?,?,?,?)`)
		if err != nil {
			return err
		}
		defer stmt2.Close()
		_, err = stmt2.Exec(c.name, key, bs, version, expired)
	}
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {

	tx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`Select cache_value from ` + c.table + ` WHERE  expired > ? AND cache_name =? AND cache_key = ? AND version=?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r := stmt.QueryRow(time.Now().Unix(), c.name, key, version)
	bs := []byte{}
	err = r.Scan(&bs)
	if err == sql.ErrNoRows {
		return nil, cache.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return bs, err
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *Cache) Flush() error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	newversion, err := cache.NewRandMaskedBytes(tokenMask, 16, version)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`update ` + c.table + ` set
	 cache_value=?,
	 version=?,
	 expired=?
	 Where cache_name=? and cache_key=""
	 `)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(
		"",
		string(newversion),
		-1,
		c.name)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		stmt2, err := tx.Prepare(`insert into ` + c.table + ` (cache_name,cache_value,version,expired,cache_key) 
		values (?,?,?,?,"")`)
		if err != nil {
			return err
		}
		defer stmt2.Close()
		_, err = stmt2.Exec(c.name, string(newversion), newversion, -1)

	}
	if err != nil {
		return err
	}
	tx.Commit()
	err = c.gc()
	if err != nil {
		return err
	}
	return err
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`DELETE FROM ` + c.table + ` WHERE cache_name= ? and cache_key = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.name, key)
	if err == nil {
		tx.Commit()
	}
	return err
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var data = make(map[string][]byte, len(keys))
	tx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	version, err := c.getVersionTx(tx)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`Select cache_value FROM ` + c.table + ` WHERE  expired > ?  AND cache_name =? AND cache_key = ? AND version=? `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	for k := range keys {
		var b []byte
		err := stmt.QueryRow(time.Now().Unix(), c.name, keys[k], version).Scan(&b)
		if err == sql.ErrNoRows {
		} else if err != nil {
			return nil, err
		} else {
			data[keys[k]] = b
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return data, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`update ` + c.table + ` set
		cache_value=?,
		version=?,
		expired=?
		Where cache_name=? 
		and cache_key=?
		`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmt2, err := tx.Prepare(`insert into ` + c.table + ` (cache_name,cache_key,cache_value,version,expired) values (?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt2.Close()
	var expired int64
	expired = time.Now().Add(ttl).Unix()
	for k := range data {
		r, err := stmt.Exec(
			data[k],
			version,
			expired,
			c.name,
			k)
		if err != nil {
			return err
		}
		affected, err := r.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {

			_, err = stmt2.Exec(c.name, k, data[k], version, expired)
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

//Expire set cache value expire duration by given key and ttl
func (c *Cache) Expire(key string, ttl time.Duration) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	var expired int64

	expired = time.Now().Add(ttl).Unix()

	stmt, err := tx.Prepare(`update ` + c.table + ` set
	 expired=?
	 Where cache_name=? 
	 and version=?
	 and cache_key=?
	 `)

	defer stmt.Close()
	_, err = stmt.Exec(
		expired,
		c.name,
		version,
		key)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	var expired int64

	expired = time.Now().Add(ttl).Unix()

	stmt, err := tx.Prepare(`update ` + c.table + ` set
	 expired=?
	 Where cache_name=? 
	 and version=?
	 and cache_key=?
	 `)

	defer stmt.Close()
	_, err = stmt.Exec(
		expired,
		c.name,
		version,
		key)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	err := c.gc()
	if err != nil {
		return nil
	}
	close(c.quit)
	return c.DB.Close()
}

//Config Cache driver config.
type Config struct {
	db.Config
	//Database table name.
	Table string
	//Database cache name.
	Name string
	//Period of gc.Default value is 5 minute.
	GCPeriod int64
	//Max delete limit in every gc call.Default value is 100.
	GCLimit int64
}

//Create create new cache driver .
//Return driver created and any error if raised.
func (cf *Config) Create() (cache.Driver, error) {
	var err error
	cache := Cache{}
	d := db.New()
	err = cf.Config.ApplyTo(d)
	if err != nil {
		return &cache, err
	}
	cache.DB = d.DB()
	cache.table = cf.Table
	cache.quit = make(chan int)
	period := time.Duration(cf.GCPeriod)
	if period == 0 {
		period = defaultGCPeriod
	}
	cache.ticker = time.NewTicker(period)
	gcLimit := cf.GCLimit
	if gcLimit == 0 {
		gcLimit = defaultGcLimit
	}
	cache.gcLimit = gcLimit
	cache.name = cf.Name
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
	err = cache.start()
	if err != nil {
		return &cache, err
	}
	return &cache, nil
}

func init() {
	cache.Register("sqlcache", func(loader func(interface{}) error) (cache.Driver, error) {
		var err error
		c := &Config{}
		err = loader(c)
		if err != nil {
			return nil, err
		}
		return c.Create()
	})
}
