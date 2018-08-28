package tencentminiprogram

import (
	"github.com/herb-go/fetch"
)

func IsErrorInvalidCode(err error) bool {
	return fetch.GetAPIErrCode(err) == "40029"
}
