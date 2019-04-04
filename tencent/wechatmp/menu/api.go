package menu

import (
	"github.com/herb-go/herbgo/util"
	"github.com/herb-go/providers/tencent/wechatmp"
)

func MustCreateMenu(App wechatmp.App, menu *Menu) {
	result := &wechatmp.ResultAPIError{}
	util.Must(App.CallJSONApiWithAccessToken(wechatmp.APIMenuCreate, nil, menu, result))
}
