package httpapiserver

import (
	"net/http"

	"github.com/herb-go/herb/user/httpuser"

	"github.com/herb-go/herb/service/httpservice/apiserver"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type Config struct {
	apiserver.Option
	Guarder guarder.DirverConfigMap
}

type Server struct {
	apiserver.Option
	Guarder *guarder.Guarder
}

func (s *Server) StartWithGuarder(h func(w http.ResponseWriter, r *http.Request)) error {
	m := httpuser.LoginRequiredMiddleware(s.Guarder, nil)
	return s.Start(func(w http.ResponseWriter, r *http.Request) {
		m(w, r, h)
	})
}
