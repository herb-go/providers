package httpapiserver

import (
	"net/http"

	"github.com/herb-go/herb/middleware"

	"github.com/herb-go/herb/service/httpservice/apiserver"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type Config struct {
	apiserver.Option
	Guarder guarder.DirverConfigMap
}

type Server struct {
	apiserver.Option
	Middlewares *middleware.Middlewares
	Action      func(w http.ResponseWriter, r *http.Request)
}

func (s *Server) StartWithMiddlewares(h func(w http.ResponseWriter, r *http.Request)) error {
	return s.Start(s.Middlewares.App(h).ServeHTTP)
}
