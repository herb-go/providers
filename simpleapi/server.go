package simpleapi

import (
	"github.com/herb-go/herb/service/httpservice/apiserver"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type ServerConfig struct {
	apiserver.Channel
	Guarder guarder.DirverConfigMap
}
