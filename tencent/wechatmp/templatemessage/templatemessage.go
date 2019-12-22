package templatemessage

import "github.com/herb-go/providers/tencent/wechatmp"

func GetAllPrivateTemplate(App *wechatmp.App) (*wechatmp.AllPrivateTemplateResult, error) {
	result := &wechatmp.AllPrivateTemplateResult{}
	err := App.CallJSONApiWithAccessToken(wechatmp.APIGetAllPrivateTemplate, nil, nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil

}
