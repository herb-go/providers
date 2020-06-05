package hireddb

import (
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/worker"
	"github.com/herb-go/providers/herb/overseers/cacheoverseer"
)

type HiredDB struct {
	PlainDB *db.PlainDB
}

func (d *HiredDB) ApplyTo(plaindb *db.PlainDB) error {
	db.Copy(d.PlainDB, plaindb)
	return nil
}
func New() *HiredDB {
	return &HiredDB{}
}
func Register() {
	db.Register("hireddb", func(c *db.Config) (db.Driver, error) {
		d := New()
		plaindb := dboverseer.GetDBByID(c.DataSource)
		if plaindb == nil {
			return nil, worker.ErrWorkerNotFound
		}
		d.PlainDB = plaindb
		return d, nil
	})
}

func init() {
	Register()
}
