package ldapuser

import (
	"github.com/herb-go/member"
	ldap "gopkg.in/ldap.v2"
)

type PasswordProvider struct {
	Config *Config
}

func (p *PasswordProvider) VerifyPassword(uid string, password string) (bool, error) {
	l, err := p.Config.BindUser(uid, password)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return false, nil
		}
		return false, err
	}
	defer l.Close()
	return true, nil
}

//PasswordChangeable return password changeable
func (p *PasswordProvider) PasswordChangeable() bool {
	return true
}

//UpdatePassword update user password
//Return any error if raised
func (p *PasswordProvider) UpdatePassword(uid string, password string) error {
	return p.Config.UpdatePassword(uid, password)
}

func (p *PasswordProvider) InstallToMember(service *member.Service) {
	service.PasswordProvider = p
}
