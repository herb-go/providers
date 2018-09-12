package vaptchacaptcha

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"

	"github.com/herb-go/providers/vaptcha"
)

type CaptchaDriver struct {
	Config *vaptcha.Config
	Type   string
}

func (d *CaptchaDriver) Name() string {
	return "vaptcha"
}

type output struct {
	VID   string
	Type  string
	Scene string
}

func (d *CaptchaDriver) MustCaptcha(scene string, reset bool, w http.ResponseWriter, r *http.Request) {
	o := output{
		VID:   d.Config.Key,
		Type:  d.Type,
		Scene: scene,
	}
	if o.Type == "" {
		o.Type = "click"
	}
	bs, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bs)
	if err != nil {
		panic(err)
	}
}
func (d *CaptchaDriver) Verify(r *http.Request, scene string, token string) (bool, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false, err
	}
	params := url.Values{}
	params.Set("id", d.Config.VID)
	params.Set("secretkey", d.Config.Key)
	params.Set("scene", scene)
	params.Set("token", token)
	params.Set("ip", ip)
	req, err := vaptcha.ApiValidate.NewRequest(nil, []byte(params.Encode()))
	if err != nil {
		return false, err
	}
	resp, err := d.Config.Clients.Fetch(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, resp
	}
	result := &vaptcha.ResultValidate{}
	err = resp.UnmarshalAsJSON(result)
	if err != nil {
		return false, err
	}
	if result.Success != 1 || result.Msg != "" {
		if result.Msg == vaptcha.MsgTokenError || result.Msg == vaptcha.MsgTokenExpired {
			return false, nil
		}
		return false, resp.NewAPICodeErr(result.Msg)
	}
	return true, nil
}
