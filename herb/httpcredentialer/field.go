package httpcredentialer

import (
	"net/http"

	"github.com/herb-go/herb/user/identifier"
	"github.com/herb-go/herb/user/identifier/protecter"
)

var headerFieldLoaderFactory = func(name string) func(r *http.Request) (identifier.CredentialData, error) {
	return func(r *http.Request) (identifier.CredentialData, error) {
		data := r.Header.Get(name)
		return identifier.CredentialData(data), nil
	}
}
var queryFieldLoaderFactory = func(name string) func(r *http.Request) (identifier.CredentialData, error) {
	return func(r *http.Request) (identifier.CredentialData, error) {
		data := r.URL.Query().Get(name)
		return identifier.CredentialData(data), nil
	}
}
var formFieldLoaderFactory = func(name string) func(r *http.Request) (identifier.CredentialData, error) {
	return func(r *http.Request) (identifier.CredentialData, error) {
		data := r.PostFormValue(name)
		return identifier.CredentialData(data), nil
	}
}
var cookieFieldLoaderFactory = func(name string) func(r *http.Request) (identifier.CredentialData, error) {
	return func(r *http.Request) (identifier.CredentialData, error) {
		c, err := r.Cookie(name)
		if err != nil {
			if err == http.ErrNoCookie {
				return identifier.CredentialData{}, nil
			}
			return nil, err
		}
		return identifier.CredentialData(c.Value), nil
	}
}

var fieldLoaderFactories = map[string]func(name string) func(r *http.Request) (identifier.CredentialData, error){
	"header": headerFieldLoaderFactory,
	"query":  queryFieldLoaderFactory,
	"form":   formFieldLoaderFactory,
	"cookie": cookieFieldLoaderFactory,
}

func NewFieldCredentialLoader(ctype identifier.CredentialType, source string, name string) (protecter.Credentialer, error) {
	loader, ok := fieldLoaderFactories[source]
	if ok == false {
		return nil, ErrUnknownSource
	}
	return &protecter.CredentialLoader{
		CredentialType: ctype,
		LoaderFunc:     loader(name),
	}, nil
}
