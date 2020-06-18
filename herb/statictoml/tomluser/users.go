package tomluser

import (
	"sync"

	"github.com/herb-go/uniqueid"

	"github.com/herb-go/providers/herb/statictoml"

	"github.com/herb-go/herb/user"
	"github.com/herb-go/member"
)

type Users struct {
	Source     statictoml.Source
	locker     sync.RWMutex
	uidmap     map[string]*User
	accountmap map[string][]*User
	idFactory  func() (string, error)
	HashMode   string
}

func NewUsers() *Users {
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
			continue
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

//SupportedStatus return supported status map
func (u *Users) SupportedStatus() map[member.Status]bool {
	return member.StatusMapMin
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

//PasswordChangeable return password changeable
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
		result[id] = user.Roles
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
	uid, err = u.accountToUID(account)
	if err != nil {
		return "", err
	}
	if uid != "" {
		return "", member.ErrAccountRegisterExists
	}
	return u.register(account)
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
	u.accountmap[account.Account] = append(u.accountmap[account.Account], accountuser)
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
		u.accountmap[a.Account] = append(u.accountmap[a.Keyword], user)
	}
}
