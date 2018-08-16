package ldapuser

import (
	"fmt"

	"github.com/herb-go/member"
	ldap "gopkg.in/ldap.v2"
)

type PasswordProvider struct {
	Config *Config
}

func (p *PasswordProvider) VerifyPassword(uid string, password string) (bool, error) {
	l, err := p.Config.Dial()
	defer l.Close()
	if err != nil {
		return false, err
	}
	err = l.Bind(fmt.Sprintf(p.Config.UserDN, uid), password)
	if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

//UpdatePassword update user password
//Return any error if raised
func (p *PasswordProvider) UpdatePassword(uid string, password string) error {
	l, err := p.Config.Dial()
	if err != nil {
		return err
	}
	defer l.Close()

	err = l.Bind(p.Config.BindDN, p.Config.BindPass)
	if err != nil {
		return err
	}
	passwordModifyRequest := ldap.NewPasswordModifyRequest(fmt.Sprintf(p.Config.UserDN, uid), "", password)
	_, err = l.PasswordModify(passwordModifyRequest)
	return err

}

func (p *PasswordProvider) InstallToMember(service *member.Service) {
	service.PasswordProvider = p
}
