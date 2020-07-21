package tomlappkey

import (
	"github.com/herb-go/herbsecurity/authority"
	"github.com/herb-go/herbsecurity/authority/service/application"
)

type App struct {
	ID    string
	Owner string
	User  string
	Key   string
}

type Data struct {
	Apps []*App
}

func NewData(apps []*App) *Data {
	return &Data{
		Apps: apps,
	}
}
func ConvertToApp(v *application.Verified) *App {
	return &App{
		ID:    string(v.Authority),
		Owner: string(v.Principal),
		User:  string(v.Agent),
		Key:   string(v.Passphrase),
	}
}
func ConvertToApps(data map[string]*application.Verified) []*App {
	result := []*App{}
	for k := range data {
		result = append(result, ConvertToApp(data[k]))
	}
	return result
}
func ConvertFromApp(a *App) *application.Verified {
	v := application.NewVerified()
	v.Authority = authority.Authority(a.ID)
	v.Principal = authority.Principal(a.Owner)
	v.Agent = authority.Agent(a.User)
	v.Passphrase = authority.Passphrase(a.Key)
	return v
}
func ConvertFromApps(a []*App) map[string]*application.Verified {
	result := map[string]*application.Verified{}
	for k := range a {
		v := ConvertFromApp(a[k])
		result[a[k].ID] = v
	}
	return result
}
