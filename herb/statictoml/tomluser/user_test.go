package tomluser

import (
	"testing"

	"github.com/herb-go/herb/user"
	"github.com/herb-go/herb/user/role"
)

func TestUser(t *testing.T) {
	u := NewUser()
	u.UID = "uid"
	u.Password = "password"
	u.HashMode = "hash"
	u.Salt = "salt"
	u.Accounts = []*user.Account{}
	u.Banned = false
	u.Roles = role.NewRoles()
	u2 := u.Clone()
	u2.UID = "uid2"
	u2.Password = "password2"
	u2.HashMode = "hash2"
	u2.Salt = "salt2"
	u2.Accounts = []*user.Account{user.NewAccount()}
	u2.Banned = true
	u2.Roles = role.NewRoles("role2")
	if u.UID == u2.UID ||
		u.Password == u2.Password ||
		u.HashMode == u2.HashMode ||
		u.Salt == u2.Salt ||
		len(u.Accounts) == len(u2.Accounts) ||
		u.Banned == u2.Banned ||
		len(*u2.Roles) == len(*u.Roles) {
		t.Fatal(u, u2)
	}
	u2.SetTo(u)
	if u.UID != u2.UID || u.UID != "uid2" ||
		u.Password != u2.Password || u.Password != "password2" ||
		u.HashMode != u2.HashMode || u.HashMode != "hash2" ||
		u.Salt != u2.Salt || u.Salt != "salt2" ||
		len(u.Accounts) != len(u2.Accounts) || len(u.Accounts) != 1 ||
		u.Banned != u2.Banned || u.Banned != true ||
		len(*u2.Roles) != len(*u.Roles) || len(*u2.Roles) != 1 {
		t.Fatal(u, u2)
	}
}
func testNewPassword(t *testing.T, u *User, hashmode string) {
	err := u.UpdatePassword(hashmode, "newpassword")
	if err != nil {
		panic(err)
	}
	result, err := u.VerifyPassword("newpassword")
	if err != nil {
		panic(err)
	}
	if !result {
		t.Fatal(result)
	}
	result, err = u.VerifyPassword("password")
	if err != nil {
		panic(err)
	}
	if result {
		t.Fatal(result)
	}
	result, err = u.VerifyPassword("wrongpassword")
	if err != nil {
		panic(err)
	}
	if result {
		t.Fatal(result)
	}
	u.UpdatePassword(hashmode, "newpassword")
	if err != nil {
		panic(err)
	}
}
func TestPassword(t *testing.T) {
	u := NewUser()
	result, err := u.VerifyPassword("")
	if err != nil {
		panic(err)
	}
	if result {
		t.Fatal(result)
	}

	testNewPassword(t, u, "sha256")
	if len(u.Password) != 64 {
		t.Fatal(u.Password)
	}
	testNewPassword(t, u, "md5")
	if len(u.Password) != 32 {
		t.Fatal(u.Password)
	}
	testNewPassword(t, u, "")
	if u.Password != "newpassword" {
		t.Fatal(u)
	}
}
