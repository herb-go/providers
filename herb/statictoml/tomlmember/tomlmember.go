package tomlmember

import (
	"sync"

	"github.com/herb-go/uniqueid"

	"github.com/herb-go/providers/herb/statictoml"

	"github.com/herb-go/herb/user/role"

	"github.com/herb-go/herb/user"
	"github.com/herb-go/member"
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

type Users struct {
	Source     statictoml.Source
	locker     sync.RWMutex
	uidmap     map[string]*User
	accountmap map[string][]*User
	idFactory  func() (string, error)
	HashMode   string
}

func newUsers() *Users {
	return &Users{
		uidmap:     map[string]*User{},
		accountmap: map[string][]*User{},
		idFactory:  uniqueid.DefaultGenerator.GenerateID,
		HashMode:   defaultUsersHashMode,
	}
}
func (u *Users) getAllUsers() *Data {
	data := NewData()
	data.Users = make([]*User, 0, len(u.uidmap))
	for k := range u.uidmap {
		data.Users = append(data.Users, u.uidmap[k])
	}
	return data
}
func (u *Users) save() error {
	return u.Source.Save(u.getAllUsers())
}

//Statuses return  status  map of given uid list.
//Return status  map and any error if raised.
func (u *Users) Statuses(uid ...string) (member.StatusMap, error) {
	u.locker.RLock()
	defer u.locker.RUnlock()
	m := member.StatusMap{}
	for _, id := range uid {
		user := u.uidmap[id]
		if user == nil {
			return nil, member.ErrUserNotFound
		}
		if user.Banned {
			m[id] = member.StatusBanned
		} else {
			m[id] = member.StatusNormal
		}
	}
	return m, u.save()
}

//SetStatus set user status.
//Return any error if raised.
func (u *Users) SetStatus(uid string, status member.Status) error {
	u.locker.Lock()
	defer u.locker.Unlock()
	if u.uidmap[uid] == nil {
		return member.ErrUserNotFound
	}
	u.uidmap[uid].Banned = !status.IsAvaliable()
	return u.save()
}

//VerifyPassword Verify user password.
//Return verify result and any error if raised
func (u *Users) VerifyPassword(uid string, password string) (bool, error) {
	u.locker.RLock()
	defer u.locker.RUnlock()
	user := u.uidmap[uid]
	if user == nil {
		return false, nil
	}
	return user.VerifyPassword(password)
}

func (u *Users) PasswordChangeable() bool {
	return true
}

//UpdatePassword update user password
//Return any error if raised
func (u *Users) UpdatePassword(uid string, password string) error {
	u.locker.Lock()
	defer u.locker.Unlock()
	user := u.uidmap[uid]
	if user == nil {
		return member.ErrUserNotFound
	}
	err := user.UpdatePassword(u.HashMode, password)
	if err != nil {
		return err
	}
	return u.save()
}

//Roles return role map of given uid list.
//Return role map and any error if raised.
func (u *Users) Roles(uid ...string) (*member.Roles, error) {
	u.locker.Lock()
	defer u.locker.Unlock()
	result := member.Roles{}
	for _, id := range uid {
		user := u.uidmap[id]
		if user == nil {
			continue
		}
		result[id] = &user.Roles
	}
	return &result, nil
}

//Accounts return account map of given uid list.
//Return account map and any error if raised.
func (u *Users) Accounts(uid ...string) (*member.Accounts, error) {
	u.locker.RLock()
	defer u.locker.RUnlock()
	a := member.Accounts{}
	for _, id := range uid {
		user := u.uidmap[id]
		if user == nil {
			continue
		}
		a[id] = user.Accounts
	}
	return &a, nil
}
func (u *Users) accountToUID(account *user.Account) (uid string, err error) {
	for _, user := range u.accountmap[account.Account] {
		for k := range user.Accounts {
			if user.Accounts[k].Equal(account) {
				return user.UID, nil
			}
		}
	}
	return "", nil
}

//AccountToUID query uid by user account.
//Return user id and any error if raised.
//Return empty string as userid if account not found.
func (u *Users) AccountToUID(account *user.Account) (uid string, err error) {
	u.locker.RLock()
	defer u.locker.RUnlock()
	return u.accountToUID(account)
}

func (u *Users) register(account *user.Account) (uid string, err error) {
	newuser := NewUser()
	id, err := u.idFactory()
	if err != nil {
		return "", err
	}
	newuser.UID = id
	newuser.Accounts = []*user.Account{account}
	u.addUser(newuser)
	err = u.save()
	if err != nil {
		return "", err
	}
	return newuser.UID, nil
}

//Register create new user with given account.
//Return created user id and any error if raised.
//Privoder should return ErrAccountRegisterExists if account is used.
func (u *Users) Register(account *user.Account) (uid string, err error) {
	u.locker.Lock()
	defer u.locker.Unlock()
	return u.Register(account)
}

//AccountToUIDOrRegister query uid by user account.Register user if account not found.
//Return user id and any error if raised.
func (u *Users) AccountToUIDOrRegister(account *user.Account) (uid string, registerd bool, err error) {
	u.locker.Lock()
	defer u.locker.Unlock()
	uid, err = u.accountToUID(account)
	if err != nil {
		return "", false, err
	}
	if uid != "" {
		return
	}
	uid, err = u.register(account)
	if err != nil {
		return "", false, err
	}
	return uid, true, nil
}

//BindAccount bind account to user.
//Return any error if raised.
//If account exists,user.ErrAccountBindingExists should be rasied.
func (u *Users) BindAccount(uid string, account *user.Account) error {
	u.locker.Lock()
	defer u.locker.Unlock()
	accountuser := u.uidmap[uid]
	if accountuser == nil {
		return member.ErrUserNotFound
	}
	accountid, err := u.accountToUID(account)
	if err != nil {
		return err
	}
	if accountid != "" {
		return user.ErrAccountBindingExists
	}
	accountuser.Accounts = append(accountuser.Accounts, account)
	u.accountmap[account.Keyword] = append(u.accountmap[account.Keyword], accountuser)
	return u.save()
}

//UnbindAccount unbind account from user.
//Return any error if raised.
//If account not exists,user.ErrAccountUnbindingNotExists should be rasied.
func (u *Users) UnbindAccount(uid string, account *user.Account) error {
	u.locker.Lock()
	defer u.locker.Unlock()
	accountid, err := u.accountToUID(account)
	if err != nil {
		return err
	}
	if accountid == "" || accountid != uid {
		return user.ErrAccountUnbindingNotExists
	}
	for k := range u.uidmap[accountid].Accounts {
		if u.uidmap[accountid].Accounts[k].Equal(account) {
			u.uidmap[accountid].Accounts = append(u.uidmap[accountid].Accounts[:k], u.uidmap[accountid].Accounts[k+1:]...)
			break
		}
	}
	for k := range u.accountmap[account.Account] {
		if u.accountmap[account.Account][k].UID == accountid {
			u.accountmap[account.Account] = append(u.accountmap[account.Account][:k], u.accountmap[account.Account][k+1:]...)
			break
		}
	}
	return u.save()
}
func (u *Users) addUser(user *User) {
	u.uidmap[user.UID] = user
	for _, a := range user.Accounts {
		u.accountmap[a.Keyword] = append(u.accountmap[a.Keyword], user)
	}
}
