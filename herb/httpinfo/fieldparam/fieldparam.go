package fieldparam

import (
	"net/http"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/responsecache"
)

type param struct {
	field httpinfo.Field
}

func (p *param) GetParam(r *http.Request) (param string, success bool) {
	data, ok, err := p.field.LoadInfo(r)
	if err != nil {
		panic(err)
	}
	return string(data), ok
}

func Wrap(f httpinfo.Field) responsecache.Param {
	return &param{
		field: f,
	}
}

func WrapFields(f ...httpinfo.Field) []responsecache.Param {
	result := make([]responsecache.Param, len(f))
	for k := range f {
		result[k] = Wrap(f[k])
	}
	return result
}
