package sqluser

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"crypto/rand"
	"crypto/sha256"

	"github.com/herb-go/app/member"
	"github.com/herb-go/herb/model-sql-datamapper"
	"github.com/herb-go/herb/model-sql-query"
	"github.com/herb-go/herb/user"
	"github.com/satori/go.uuid"
)

const (
	FlagEmpty        = 0
	FlagWithAccount  = 1
	FlagWithPassword = 2
	FlagWithToken    = 4
	FlagWithUser     = 8
)

const (
	UserStatusNormal = 0
	UserStatusBanned = 1
)

var RandomBytesLength = 32
var ErrHashMethodNotFound = errors.New("password hash method not found")

type HashFunc func(key string, salt string, password string) ([]byte, error)

var DefaultAccountTableName = "account"
var DefaultPasswordTableName = "password"
var DefaultTokenTableName = "token"
var DefaultUserTableName = "user"
var DefaultHashMethod = "sha256"
var HashFuncMap = map[string]HashFunc{
	"sha256": func(key string, salt string, password string) ([]byte, error) {
		var val = []byte(key + salt + password)
		var s256 = sha256.New()
		s256.Write(val)
		val = s256.Sum(nil)
		s256.Write(val)
		return []byte(hex.EncodeToString(s256.Sum(nil))), nil
	},
}

func New(db *sql.DB, prefix string, flag int) *User {
	return &User{
		DB: datamapper.NewDB(db, prefix),
		Tables: Tables{
			AccountTableName:  DefaultAccountTableName,
			PasswordTableName: DefaultPasswordTableName,
			TokenTableName:    DefaultTokenTableName,
			UserTableName:     DefaultUserTableName,
		},
		HashMethod:     DefaultHashMethod,
		UIDGenerater:   UUID,
		TokenGenerater: Timestamp,
		SaltGenerater:  RandomBytes,
		Flag:           flag,
	}
}

type Tables struct {
	AccountTableName  string
	PasswordTableName string
	TokenTableName    string
	UserTableName     string
}
type User struct {
	DB             datamapper.DB
	Tables         Tables
	Flag           int
	UIDGenerater   func() (string, error)
	TokenGenerater func() (string, error)
	SaltGenerater  func() (string, error)
	HashMethod     string
	PasswordKey    string
}

func RandomBytes() (string, error) {
	var bytes = make([]byte, RandomBytesLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func UUID() (string, error) {
	return uuid.NewV1().String(), nil
}
func Timestamp() (string, error) {
	return strconv.FormatInt(time.Now().UnixNano(), 10), nil
}
func (u *User) HasFlag(flag int) bool {
	return u.Flag&flag != 0
}
func (u *User) AccountTableName() string {
	return u.DB.TableName(u.Tables.AccountTableName)
}
func (u *User) PasswordTableName() string {
	return u.DB.TableName(u.Tables.PasswordTableName)
}
func (u *User) TokenTableName() string {
	return u.DB.TableName(u.Tables.TokenTableName)
}
func (u *User) UserTableName() string {
	return u.DB.TableName(u.Tables.UserTableName)
}
func (u *User) Account() *AccountDataMapper {
	return &AccountDataMapper{
		DataMapper: datamapper.New(u.DB, u.Tables.AccountTableName),
		User:       u,
	}
}

func (u *User) Password() *PasswordDataMapper {
	return &PasswordDataMapper{
		DataMapper: datamapper.New(u.DB, u.Tables.PasswordTableName),
		User:       u,
	}
}
func (u *User) Token() *TokenDataMapper {
	return &TokenDataMapper{
		DataMapper: datamapper.New(u.DB, u.Tables.TokenTableName),
		User:       u,
	}
}
func (u *User) User() *UserDataMapper {
	return &UserDataMapper{
		DataMapper: datamapper.New(u.DB, u.Tables.UserTableName),
		User:       u,
	}
}

type AccountDataMapper struct {
	datamapper.DataMapper
	User    *User
	Service *member.Service
}

func (a *AccountDataMapper) InstallToMember(service *member.Service) {
	service.AccountsProvider = a
	a.Service = service
}

func (a *AccountDataMapper) Unbind(uid string, account user.UserAccount) error {
	tx, err := a.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	Delete := query.NewDelete(a.DBTableName())
	Delete.Where.Condition = query.And(
		query.Equal("account.uid", uid),
		query.Equal("account.keyword", account.Keyword),
		query.Equal("account.account", account.Account),
	)
	_, err = Delete.Query().Exec(tx)
	if err != nil {
		return err
	}
	return tx.Commit()

}

func (a *AccountDataMapper) Bind(uid string, account user.UserAccount) error {
	tx, err := a.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var u = ""
	Select := query.NewSelect()
	Select.Select.Add("account.uid")
	Select.From.AddAlias("account", a.DBTableName())
	Select.Where.Condition = query.And(
		query.Equal("keyword", account.Keyword),
		query.Equal("account", account.Account),
	)
	row := Select.QueryRow(a.DB())
	err = row.Scan(&u)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return user.ErrAccountBindExists

	}

	var CreatedTime = time.Now().Unix()
	Insert := query.NewInsert(a.DBTableName())
	Insert.Insert.
		Add("uid", uid).
		Add("keyword", account.Keyword).
		Add("account", account.Account).
		Add("created_time", CreatedTime)
	_, err = Insert.Query().Exec(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (a *AccountDataMapper) FindOrInsert(UIDGenerater func() (string, error), account user.UserAccount) (string, error) {
	var result = AccountModel{}
	tx, err := a.DB().Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	Select := query.NewSelect()
	Select.From.AddAlias("account", a.DBTableName())
	Select.Select.Add("account.uid", "account.keyword", "account.account", "account.created_time")
	Select.Where.Condition = query.And(
		query.Equal("account.keyword", account.Keyword),
		query.Equal("account.account", account.Account),
	)
	row := Select.QueryRow(a.DB())
	err = Select.Result().
		Bind("account.uid", &result.UID).
		Bind("account.keyword", &result.Keyword).
		Bind("account.account", &result.Account).
		Bind("account.created_time", &result.CreatedTime).
		ScanFrom(row)
	if err == nil {
		return result.UID, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}
	uid, err := UIDGenerater()
	var CreatedTime = time.Now().Unix()
	Insert := query.NewInsert(a.DBTableName())
	Insert.Insert.
		Add("uid", uid).
		Add("keyword", account.Keyword).
		Add("account", account.Account).
		Add("created_time", CreatedTime)
	_, err = Insert.Query().Exec(tx)
	if err != nil {
		return "", err
	}
	if a.User.HasFlag(FlagWithUser) {
		Insert := query.NewInsert(a.User.UserTableName())
		Insert.Insert.
			Add("uid", uid).
			Add("status", UserStatusNormal).
			Add("created_time", CreatedTime).
			Add("updated_time", CreatedTime)
		_, err = Insert.Query().Exec(tx)
		if err != nil {
			return "", err
		}
	}
	return uid, tx.Commit()
}
func (a *AccountDataMapper) Insert(uid string, keyword string, account string) error {
	tx, err := a.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var u = ""
	Select := query.NewSelect()
	Select.Select.Add("uid")
	Select.From.Add(a.DBTableName())
	Select.Where.Condition = query.And(
		query.Equal("keyword", keyword),
		query.Equal("account", account),
	)
	row := Select.QueryRow(a.DB())
	err = row.Scan(&u)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return member.ErrAccountRegisterExists
	}
	var CreatedTime = time.Now().Unix()
	Insert := query.NewInsert(a.DBTableName())
	Insert.Insert.
		Add("uid", uid).
		Add("keyword", keyword).
		Add("account", account).
		Add("created_time", CreatedTime)
	_, err = Insert.Query().Exec(tx)
	if err != nil {
		return err
	}
	if a.User.HasFlag(FlagWithUser) {
		Insert := query.NewInsert(a.User.UserTableName())
		Insert.Insert.
			Add("uid", uid).
			Add("status", UserStatusNormal).
			Add("created_time", CreatedTime).
			Add("updated_time", CreatedTime)
		_, err = Insert.Query().Exec(tx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
func (a *AccountDataMapper) Find(keyword string, account string) (AccountModel, error) {
	var result = AccountModel{}
	if keyword == "" || account == "" {
		return result, sql.ErrNoRows
	}
	Select := query.NewSelect()
	Select.Select.Add("uid", "keyword", "account", "created_time")
	Select.From.Add(a.DBTableName())
	Select.Where.Condition = query.And(
		query.Equal("keyword", keyword),
		query.Equal("account", account),
	)
	row := Select.QueryRow(a.DB())
	err := Select.Result().
		Bind("uid", &result.UID).
		Bind("keyword", &result.Keyword).
		Bind("account", &result.Account).
		Bind("created_time", &result.CreatedTime).
		ScanFrom(row)
	return result, err
}
func (a *AccountDataMapper) FindAllByUID(uids ...string) ([]AccountModel, error) {
	var result = []AccountModel{}
	if len(uids) == 0 {
		return result, nil
	}
	Select := query.NewSelect()
	Select.Select.Add("account.uid", "account.keyword", "account.account")
	Select.From.AddAlias("account", a.DBTableName())
	Select.Where.Condition = query.In("account.uid", uids)
	rows, err := Select.QueryRows(a.DB())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := AccountModel{}
		err := Select.Result().
			Bind("account.uid", &v.UID).
			Bind("account.keyword", &v.Keyword).
			Bind("account.account", &v.Account).
			ScanFrom(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
func (a *AccountDataMapper) Accounts(uid ...string) (member.Accounts, error) {
	models, err := a.FindAllByUID(uid...)
	if err != nil {
		return nil, err
	}
	result := member.Accounts{}
	for _, v := range models {
		if result[v.UID] == nil {
			result[v.UID] = user.UserAccounts{}
		}
		account := user.UserAccount{Keyword: v.Keyword, Account: v.Account}
		result[v.UID] = append(result[v.UID], account)
	}
	return result, nil
}
func (a *AccountDataMapper) AccountToUID(account user.UserAccount) (uid string, err error) {
	model, err := a.Find(account.Keyword, account.Account)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return model.UID, err
}
func (a *AccountDataMapper) Register(account user.UserAccount) (uid string, err error) {
	uid, err = a.User.UIDGenerater()
	if err != nil {
		return
	}
	err = a.Insert(uid, account.Keyword, account.Account)
	return
}
func (a *AccountDataMapper) AccountToUIDOrRegister(account user.UserAccount) (uid string, err error) {
	return a.FindOrInsert(a.User.UIDGenerater, account)
}
func (a *AccountDataMapper) BindAccounts(uid string, account user.UserAccount) error {
	return a.Bind(uid, account)
}
func (a *AccountDataMapper) UnbindAccounts(uid string, account user.UserAccount) error {
	return a.Unbind(uid, account)
}

type AccountModel struct {
	UID         string
	Keyword     string
	Account     string
	CreatedTime int64
}
type PasswordDataMapper struct {
	datamapper.DataMapper
	User    *User
	Service *member.Service
}

func (p *PasswordDataMapper) InstallToMember(service *member.Service) {
	service.PasswordProvider = p
	p.Service = service
}
func (p *PasswordDataMapper) Find(uid string) (PasswordModel, error) {
	var result = PasswordModel{}
	if uid == "" {
		return result, sql.ErrNoRows
	}
	Select := query.NewSelect()
	Select.Select.Add("password.hash_method", "password.salt", "password.password", "password.updated_time")
	Select.From.AddAlias("password", p.DBTableName())
	Select.Where.Condition = query.Equal("uid", uid)
	q := Select.Query()
	row := p.DB().QueryRow(q.QueryCommand(), q.QueryArgs()...)
	result.UID = uid
	args := Select.Result().
		Bind("password.hash_method", &result.HashMethod).
		Bind("password.salt", &result.Salt).
		Bind("password.password", &result.Password).
		Bind("password.updated_time", &result.UpdatedTime).
		Args()
	err := row.Scan(args...)
	return result, err
}
func (p *PasswordDataMapper) InsertOrUpdate(model *PasswordModel) error {
	tx, err := p.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	Update := query.NewUpdate(p.DBTableName())
	Update.Update.
		Add("hash_method", model.HashMethod).
		Add("salt", model.Salt).
		Add("password", model.Password).
		Add("updated_time", model.UpdatedTime)
	Update.Where.Condition = query.Equal("uid", model.UID)
	r, err := Update.Query().Exec(tx)

	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 0 {
		return tx.Commit()
	}
	Insert := query.NewInsert(p.DBTableName())
	Insert.Insert.
		Add("uid", model.UID).
		Add("hash_method", model.HashMethod).
		Add("salt", model.Salt).
		Add("password", model.Password).
		Add("updated_time", model.UpdatedTime)
	_, err = Insert.Query().Exec(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (p *PasswordDataMapper) VerifyPassword(uid string, password string) (bool, error) {
	model, err := p.Find(uid)
	if err == sql.ErrNoRows {
		return false, member.ErrUserNotFound
	}
	if err != nil {
		return false, err
	}
	hash := HashFuncMap[model.HashMethod]
	if hash == nil {
		return false, ErrHashMethodNotFound
	}
	hashed, err := hash(p.User.PasswordKey, model.Salt, password)
	if err != nil {
		return false, err
	}
	return bytes.Compare(hashed, model.Password) == 0, nil
}
func (p *PasswordDataMapper) UpdatePassword(uid string, password string) error {
	salt, err := p.User.SaltGenerater()
	if err != nil {
		return err
	}
	hash := HashFuncMap[p.User.HashMethod]
	if hash == nil {
		return ErrHashMethodNotFound
	}
	hashed, err := hash(p.User.PasswordKey, salt, password)
	if err != nil {
		return err
	}
	model := &PasswordModel{
		UID:         uid,
		HashMethod:  p.User.HashMethod,
		Salt:        salt,
		Password:    hashed,
		UpdatedTime: time.Now().Unix(),
	}
	return p.InsertOrUpdate(model)
}

type PasswordModel struct {
	UID         string
	HashMethod  string
	Salt        string
	Password    []byte
	UpdatedTime int64
}

type TokenDataMapper struct {
	datamapper.DataMapper
	User    *User
	Service *member.Service
}

func (t *TokenDataMapper) InstallToMember(service *member.Service) {
	service.TokenProvider = t
	t.Service = service
}

func (t *TokenDataMapper) InsertOrUpdate(uid string, token string) error {
	tx, err := t.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var CreatedTime = time.Now().Unix()
	Update := query.NewUpdate(t.DBTableName())
	Update.Update.
		Add("token", token).
		Add("updated_time", CreatedTime)
	Update.Where.Condition = query.Equal("uid", uid)
	r, err := Update.Query().Exec(tx)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 0 {
		return tx.Commit()
	}
	Insert := query.NewInsert(t.DBTableName())
	Insert.Insert.
		Add("uid", uid).
		Add("token", token).
		Add("updated_time", CreatedTime)
	_, err = Insert.Query().Exec(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (t *TokenDataMapper) FindAllByUID(uids ...string) ([]TokenModel, error) {
	var result = []TokenModel{}
	if len(uids) == 0 {
		return result, nil
	}
	Select := query.NewSelect()
	Select.Select.Add("token.uid", "token.token")
	Select.From.AddAlias("token", t.DBTableName())
	Select.Where.Condition = query.In("token.uid", uids)
	rows, err := Select.QueryRows(t.DB())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := TokenModel{}
		err = rows.Scan(&v.UID, &v.Token)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
func (t *TokenDataMapper) Tokens(uid ...string) (member.Tokens, error) {
	models, err := t.FindAllByUID(uid...)
	if err != nil {
		return nil, err
	}
	result := member.Tokens{}
	for _, v := range models {
		result[v.UID] = v.Token
	}
	return result, nil

}
func (t *TokenDataMapper) Revoke(uid string) (string, error) {
	token, err := t.User.TokenGenerater()
	if err != nil {
		return "", err
	}
	return token, t.InsertOrUpdate(uid, token)
}

type TokenModel struct {
	UID         string
	Token       string
	UpdatedTime string
}
type UserDataMapper struct {
	datamapper.DataMapper
	User    *User
	Service *member.Service
}

func (u *UserDataMapper) InstallToMember(service *member.Service) {
	service.BannedProvider = u
	u.Service = service
}

func (u *UserDataMapper) FindAllByUID(uids ...string) ([]UserModel, error) {
	var result = []UserModel{}
	if len(uids) == 0 {
		return result, nil
	}
	Select := query.NewSelect()
	Select.Select.Add("user.uid", "user.status")
	Select.From.AddAlias("user", u.DBTableName())
	Select.Where.Condition = query.In("user.uid", uids)
	rows, err := Select.QueryRows(u.DB())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := UserModel{}
		err = rows.Scan(&v.UID, &v.Status)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func (u *UserDataMapper) InsertOrUpdate(uid string, status int) error {
	tx, err := u.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var CreatedTime = time.Now().Unix()
	Update := query.NewUpdate(u.DBTableName())
	Update.Update.
		Add("status", status).
		Add("updated_time", CreatedTime)
	Update.Where.Condition = query.Equal("uid", uid)
	r, err := Update.Query().Exec(tx)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 0 {
		return tx.Commit()
	}
	Insert := query.NewInsert(u.DBTableName())
	Insert.Insert.
		Add("uid", uid).
		Add("status", status).
		Add("updated_time", CreatedTime).
		Add("created_time", CreatedTime)
	_, err = Insert.Query().Exec(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (u *UserDataMapper) Banned(uid ...string) (member.BannedMap, error) {
	models, err := u.FindAllByUID(uid...)
	if err != nil {
		return nil, err
	}
	result := member.BannedMap{}
	for _, v := range models {
		result[v.UID] = (v.Status == UserStatusBanned)
	}
	return result, nil
}
func (u *UserDataMapper) Ban(uid string, banned bool) error {
	var status int
	if banned {
		status = UserStatusBanned
	} else {
		status = UserStatusNormal
	}
	return u.InsertOrUpdate(uid, status)

}

type UserModel struct {
	UID         string
	CreatedTime int64
	UpdateTIme  int64
	Status      int
}
