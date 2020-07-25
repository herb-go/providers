package tomluser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/herb-go/herb/user"

	"github.com/herb-go/herbsecurity/authorize/role"

	"github.com/herb-go/member"
	"github.com/herb-go/providers/herb/statictoml"
)

func TestConfig(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if tmpdir == "" || err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	source := path.Join(tmpdir, "test.static.toml")
	err = ioutil.WriteFile(source, []byte{}, 0700)
	if err != nil {
		panic(err)
	}
	data := NewData()
	acc := user.NewAccount()
	acc.Account = "testaccount"
	acc.Keyword = "testkeyword"
	accnotexist := user.NewAccount()
	acc.Account = "testaccountnotexist"
	acc.Keyword = "testkeyword"
	acctobind := user.NewAccount()
	acctobind.Account = "acctobind"
	acctobind.Keyword = "testkeyword"
	acctoregister := user.NewAccount()
	acctoregister.Account = "acctoregister"
	acctoregister.Keyword = "testkeyword"

	u := NewUser()
	u.UID = "testuid"
	u.Password = "password"
	u.Accounts = append(u.Accounts, acc)
	u.Roles.Append(role.NewRole("admin"))
	usertobind := NewUser()
	usertobind.UID = "usertobind"
	userbanned := NewUser()
	userbanned.UID = "userbanned"
	userbanned.Banned = true
	data.Users = append(data.Users, u, usertobind, userbanned)
	err = statictoml.Source(source).Save(data)
	if err != nil {
		panic(err)
	}
	c := &Config{
		Source:             statictoml.Source(source),
		AsPasswordProvider: true,
		AsStatusProvider:   true,
		AsAccountsProvider: true,
		AsRoleProvider:     true,
	}
	m := member.New()
	err = c.Execute(m)
	if err != nil {
		panic(err)
	}
	as := member.NewAccountsStore()
	err = m.Accounts().Load(as, "testuid", "testuidnotexists")
	if err != nil {
		panic(err)
	}
	accresult := as.Get("testuid")
	if accresult == nil || !accresult.Exists(acc) {
		t.Fatal(accresult)
	}
	uid, err := m.Accounts().AccountToUID(acc)
	if err != nil {
		panic(err)
	}
	if uid != "testuid" {
		t.Fatal(uid)
	}
	uid, err = m.Accounts().AccountToUID(accnotexist)
	if err != nil {
		panic(err)
	}
	if uid != "" {
		t.Fatal(uid)
	}
	uid, err = m.Accounts().Register(acc)
	if uid != "" || err != member.ErrAccountRegisterExists {
		t.Fatal(uid, err)
	}
	uid, err = m.Accounts().Register(acctoregister)
	if uid == "" || err != nil {
		t.Fatal(uid, err)
	}

	uid, ok, err := m.Accounts().AccountToUIDOrRegister(acc)
	if uid != "testuid" || ok != false || err != nil {
		t.Fatal(uid, ok, err)
	}
	uid, ok, err = m.Accounts().AccountToUIDOrRegister(accnotexist)
	if uid == "" || ok != true || err != nil {
		t.Fatal(uid, ok, err)
	}
	err = m.Accounts().BindAccount("uidnotexist", acctobind)
	if err != member.ErrUserNotFound {
		t.Fatal(err)
	}
	err = m.Accounts().BindAccount("uidnotexist", acctobind)
	if err != member.ErrUserNotFound {
		t.Fatal(err)
	}

	err = m.Accounts().BindAccount(usertobind.UID, acc)
	if err != user.ErrAccountBindingExists {
		t.Fatal(err)
	}
	uid, err = m.Accounts().AccountToUID(acctobind)
	if err != nil {
		panic(err)
	}
	if uid != "" {
		t.Fatal(uid)
	}

	err = m.Accounts().BindAccount(usertobind.UID, acctobind)
	if err != nil {
		t.Fatal(err)
	}
	uid, err = m.Accounts().AccountToUID(acctobind)
	if err != nil {
		panic(err)
	}
	if uid != usertobind.UID {
		t.Fatal(uid)
	}
	err = m.Accounts().UnbindAccount("", acctobind)
	if err != user.ErrAccountUnbindingNotExists {
		t.Fatal(err)
	}
	err = m.Accounts().UnbindAccount("uidnotexists", acctobind)
	if err != user.ErrAccountUnbindingNotExists {
		t.Fatal(err)
	}
	err = m.Accounts().UnbindAccount(usertobind.UID, acctobind)
	if err != nil {
		t.Fatal(err)
	}
	uid, err = m.Accounts().AccountToUID(acctobind)
	if err != nil {
		panic(err)
	}
	if uid != "" {
		t.Fatal(uid)
	}
	rs := member.NewRolesStore()
	err = m.Roles().Load(rs, u.UID, "uidnotexists")
	if err != nil {
		panic(err)
	}
	uroles := rs.Get(u.UID)
	ok, err = uroles.Authorize(role.NewPlainRoles("rolenotexists"))
	if ok == true || err != nil {
		t.Fatal(ok, err)
	}
	ok, err = uroles.Authorize(role.NewPlainRoles("admin"))
	if ok == false || err != nil {
		t.Fatal(ok, err)
	}
	ok = m.Password().PasswordChangeable()
	if !ok {
		t.Fatal(ok)
	}
	ok, err = m.Password().VerifyPassword("usernotexist", "password")
	if ok || err != nil {
		t.Fatal(ok, err)
	}
	ok, err = m.Password().VerifyPassword(u.UID, "password")
	if !ok || err != nil {
		t.Fatal(ok, err)
	}
	err = m.Password().UpdatePassword(u.UID, "newpassword")
	if err != nil {
		panic(err)
	}
	ok, err = m.Password().VerifyPassword(u.UID, "password")
	if ok || err != nil {
		t.Fatal(ok, err)
	}
	ok, err = m.Password().VerifyPassword(u.UID, "newpassword")
	if !ok || err != nil {
		t.Fatal(ok, err)
	}
	err = m.Password().UpdatePassword("usernotexist", "newpassword")
	if err != member.ErrUserNotFound {
		t.Fatal(err)
	}
	ss := member.NewStatusStore()
	err = m.Status().Load(ss, u.UID, "notexsits", userbanned.UID)
	if err != nil {
		t.Fatal(err)
	}
	if status := ss.Get(u.UID); *status != member.StatusNormal {
		t.Fatal(status)
	}
	if status := ss.Get("notexists"); status != nil {
		t.Fatal(status)
	}
	if status := ss.Get(userbanned.UID); *status != member.StatusBanned {
		t.Fatal(status)
	}
	err = m.Status().SetStatus(userbanned.UID, member.StatusNormal)
	if err != nil {
		t.Fatal(err)
	}
	ss = member.NewStatusStore()
	err = m.Status().Load(ss, userbanned.UID)
	if err != nil {
		t.Fatal(err)
	}
	if status := ss.Get(userbanned.UID); *status != member.StatusNormal {
		t.Fatal(status)
	}
	err = m.Status().SetStatus("notexists", member.StatusNormal)
	if err != member.ErrUserNotFound {
		t.Fatal(err)
	}
	configbytes, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	loader := json.NewDecoder(bytes.NewBuffer(configbytes)).Decode
	m = member.New()
	d, err := DirectiveFactory(loader)
	if err != nil {
		t.Fatal(err)
	}
	err = d.Execute(m)
	if err != nil {
		t.Fatal(err)
	}
	uid, err = m.Accounts().AccountToUID(acctoregister)
	if err != nil {
		panic(err)
	}
	if uid == "" {
		t.Fatal(uid)
	}
	ok, err = m.Password().VerifyPassword(u.UID, "newpassword")
	if !ok || err != nil {
		t.Fatal(ok, err)
	}
	ss = member.NewStatusStore()
	err = m.Status().Load(ss, userbanned.UID)
	if err != nil {
		t.Fatal(err)
	}
	if status := ss.Get(userbanned.UID); *status != member.StatusNormal {
		t.Fatal(status)
	}
}
