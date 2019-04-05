package menu

import (
	"github.com/herb-go/providers/tencent/wechatmp"
)

func CreateMenu(App *wechatmp.App, menu *Menu) error {
	result := &wechatmp.ResultAPIError{}
	return App.CallJSONApiWithAccessToken(wechatmp.APIMenuCreate, nil, menu, result)
}

func GetMenu(App *wechatmp.App) (*MenuResult, error) {
	menu := NewMenuResult()
	result := &wechatmp.ResultAPIError{}
	err := App.CallJSONApiWithAccessToken(wechatmp.APIMenuGet, nil, menu, result)
	return menu, err
}
