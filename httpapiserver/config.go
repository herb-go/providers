package httpapiserver

import (
	"github.com/herb-go/herb/service/httpservice/apiserver"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type Config struct {
	apiserver.Option
	Guarder guarder.DirverConfigMap
}
