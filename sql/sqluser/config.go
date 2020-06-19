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

func (c *Config) ApplyToUser(u *User) error {
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
	u.DB = database
	u.Flag = flag
	u.UIDGenerater = uniqueid.DefaultGenerator.GenerateID
	u.Tables.AccountMapperName = c.TableAccount
	u.Tables.PasswordMapperName = c.TablePassword
	u.Tables.UserMapperName = c.TableUser
	u.Tables.TokenMapperName = c.TableToken
	u.AddTablePrefix(c.Prefix)
	return nil
}
func (c *Config) Execute(s *member.Service) error {
	u := New(nil, nil, 0)
	err := c.ApplyToUser(u)
	if err != nil {
		return err
	}
	if c.TableAccount != "" {
		u.Account().Execute(s)
	}
	if c.TablePassword != "" {
		u.Password().Execute(s)
	}
	if c.TableUser != "" {
		u.User().Execute(s)
	}
	if c.TableToken != "" {
		u.Token().Execute(s)
	}
	return nil
}

var DirectiveFactory = func(loader func(v interface{}) error) (member.Directive, error) {
	c := &Config{}
	err := loader(c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
