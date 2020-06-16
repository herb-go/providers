package tomlmember

import (
	"github.com/herb-go/herb/user"
	"github.com/herb-go/herb/user/role"
)

var defaultUsersHashMode = "sha256"

type User struct {
	UID      string
	Password string
	HashMode string
	Salt     string
	Accounts []*user.Account
	Banned   bool
	Roles    role.Roles
}

func (u *User) Clone() *User {
	newuser := NewUser()
	newuser.UID = u.UID
	newuser.HashMode = u.HashMode
	newuser.Salt = u.Salt
	newuser.Accounts = make([]*user.Account, len(newuser.Accounts))
	copy(newuser.Accounts, u.Accounts)
	newuser.Banned = u.Banned
	newuser.Roles = make(role.Roles, len(u.Roles))
	copy(newuser.Roles, u.Roles)
	return u
}
func (u *User) SetTo(newuser *User) {
	newuser.UID = u.UID
	newuser.HashMode = u.HashMode
	newuser.Salt = u.Salt
	newuser.Accounts = u.Accounts
	newuser.Banned = u.Banned
	newuser.Roles = u.Roles
}

func (u *User) VerifyPassword(password string) (bool, error) {
	if u.Password == "" {
		return false, nil
	}
	hashed, err := Hash(u.HashMode, password, u)
	if err != nil {
		return false, err
	}
	return hashed == u.Password, nil
}
func (u *User) UpdatePassword(hashmode string, password string) error {
	newuser := u.Clone()
	newuser.HashMode = hashmode
	newuser.Salt = getSalt(saltlength)
	hashed, err := Hash(hashmode, password, newuser)
	if err != nil {
		return err
	}
	newuser.Password = hashed
	newuser.SetTo(u)
	return nil
}
func NewUser() *User {
	return &User{}
}
