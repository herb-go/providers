package qrcode

import "github.com/herb-go/providers/tencent/wechatmp"

type QRCodeScene struct {
	SceneID  *int    `json:"scene_id"`
	SceneStr *string `json:"scene_str"`
}
type QRCodeActionInfo struct {
	Scene QRCodeScene `json:"scene"`
}
type QRCodeConfig struct {
	ExpireSeconds *int64           `json:"expire_seconds"`
	ActionName    string           `json:"action_name"`
	ActionInfo    QRCodeActionInfo `json:"action_info"`
}

func NewQRCodeConfig() *QRCodeConfig {
	return &QRCodeConfig{}
}
func CreateQRCode(App *wechatmp.App, c *QRCodeConfig) (*wechatmp.ResultQRCodeCreate, error) {
	result := &wechatmp.ResultQRCodeCreate{}
	err := App.CallJSONApiWithAccessToken(wechatmp.APIQRCodeCreate, nil, c, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CreateLimitStrScene(App *wechatmp.App, code string) (*wechatmp.ResultQRCodeCreate, error) {
	str := code
	c := NewQRCodeConfig()
	c.ActionName = "QR_LIMIT_STR_SCENE"
	c.ActionInfo.Scene.SceneStr = &str
	return CreateQRCode(App, c)
}

func CreateLimitScene(App *wechatmp.App, code int) (*wechatmp.ResultQRCodeCreate, error) {
	int := code
	c := NewQRCodeConfig()
	c.ActionName = "QR_LIMIT_SCENE"
	c.ActionInfo.Scene.SceneID = &int
	return CreateQRCode(App, c)
}

func CreateStrScene(App *wechatmp.App, code string, expireSeconds int64) (*wechatmp.ResultQRCodeCreate, error) {
	str := code
	ex := expireSeconds
	c := NewQRCodeConfig()
	c.ExpireSeconds = &ex
	c.ActionName = "QR_LIMIT_STR_SCENE"
	c.ActionInfo.Scene.SceneStr = &str
	return CreateQRCode(App, c)
}

func CreateScene(App *wechatmp.App, code int, expireSeconds int64) (*wechatmp.ResultQRCodeCreate, error) {
	int := code
	ex := expireSeconds
	c := NewQRCodeConfig()
	c.ExpireSeconds = &ex
	c.ActionName = "QR_LIMIT_SCENE"
	c.ActionInfo.Scene.SceneID = &int
	return CreateQRCode(App, c)
}
