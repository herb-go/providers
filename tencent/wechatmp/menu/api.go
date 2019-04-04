package menu

import (
	"github.com/herb-go/providers/tencent/wechatmp"
)

func CreateMenu(App *wechatmp.App, menu *Menu) error {
	result := &wechatmp.ResultAPIError{}
	return App.CallJSONApiWithAccessToken(wechatmp.APIMenuCreate, nil, menu, result)
}
