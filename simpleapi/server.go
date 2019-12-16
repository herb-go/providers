package simpleapi

import (
	"github.com/herb-go/herb/service/httpservice/channel"
	"github.com/herb-go/herb/service/httpservice/guarder"
)

type ServerConfig struct {
	channel.Channel
	Guarder guarder.DriverConfig
}
