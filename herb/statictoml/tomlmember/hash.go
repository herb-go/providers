package tomlmember

import (
	"crypto/md5"
	"crypto/sha256"
	"math/rand"
)

var saltlength = 8
var saltchars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getSalt(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result = result + string(saltchars[rand.Intn(len(saltchars))])
	}
	return result
}
func Hash(mode string, password string, user *User) (string, error) {
	switch mode {
	case "md5":
		data := md5.Sum([]byte(password + user.Salt))
		return string(data[:]), nil
	case "sha256":
		data := sha256.Sum256([]byte(password + user.Salt))
		return string(data[:]), nil
	}
	return password, nil
}
