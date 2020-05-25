package sqluser

import (
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/member"
	"github.com/herb-go/uniqueid"
)

type Config struct {
	Database      *db.Config
	TableAccount  string
	TablePassword string
	TableToken    string
	TableUser     string
	Prefix        string
}

func (c *Config) ApplyTo(s *member.Service) error {
	var err error
	database := db.New()
	err = c.Database.ApplyTo(database)
	if err != nil {
		return err
	}
	var flag int
	if c.TableAccount != "" {
		flag = flag | FlagWithAccount
	}
	if c.TablePassword != "" {
		flag = flag | FlagWithPassword
	}
	if c.TableAccount != "" {
		flag = flag | FlagWithAccount
	}
	if c.TableToken != "" {
		flag = flag | FlagWithToken
	}
	u := New(database, uniqueid.DefaultGenerator.GenerateID, flag)
	u.Tables.AccountMapperName = c.TableAccount
	u.Tables.PasswordMapperName = c.TablePassword
	u.Tables.UserMapperName = c.TableUser
	u.Tables.TokenMapperName = c.TableToken
	u.AddTablePrefix(c.Prefix)
	if c.TableAccount != "" {
		s.Install(u.Account())
	}
	if c.TablePassword != "" {
		s.Install(u.Password())
	}
	if c.TableUser != "" {
		s.Install(u.User())
	}
	if c.TableToken != "" {
		s.Install(u.Token())
	}
	return nil
}

func Register() {
	member.Register("sqluser", func(loader func(v interface{}) error) (member.Directive, error) {
		d := &Config{}
		err := loader(d)
		if err != nil {
			return nil, err
		}
		return d, nil
	})
}

func init() {
	Register()
}
