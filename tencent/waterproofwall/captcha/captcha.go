package waterproffwallcaptcha

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/session"
	"github.com/herb-go/herb/cache/session/captcha"
	"github.com/herb-go/providers/tencent/waterproofwall"
)

type Driver struct {
	waterproofwall.App
}

func (d *Driver) Name() string {
	return "waterproofwall"
}

type output struct {
	AppID string
}

func (d *Driver) MustCaptcha(s *session.Store, w http.ResponseWriter, r *http.Request, scene string, reset bool) {
	o := output{
		AppID: d.AppID,
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

var failResponseMap = map[string]bool{
	"0":   true,
	"104": true,
	"9":   true,
}

func (d *Driver) Verify(s *session.Store, r *http.Request, scene string, token string) (bool, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false, err
	}
	tokens := strings.SplitN(token, "\n", 2)
	if len(tokens) < 2 {
		return false, nil
	}
	if tokens[0] == "" || tokens[1] == "" {
		return false, nil
	}
	params := url.Values{}
	params.Add("aid", d.AppID)
	params.Add("AppSecretKey", d.AppSecretKey)
	params.Add("Ticket", tokens[0])
	params.Add("Randstr", tokens[1])
	params.Add("UserIP", ip)
	req, err := waterproofwall.ApiValidate.NewRequest(params, nil)
	if err != nil {
		return false, err
	}
	resp, err := d.Clients.Fetch(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, resp
	}
	result := &waterproofwall.ResultValidate{}
	err = resp.UnmarshalAsJSON(result)
	if err != nil {
		return false, err
	}
	if result.Response == "1" {
		return true, nil
	} else if failResponseMap[result.Response] {
		return false, nil
	}
	return false, resp.NewAPICodeErr(result.ErrMsg)
}

func init() {
	captcha.Register("tcaptcha", func(conf cache.Config, prefix string) (captcha.Driver, error) {
		var err error
		c := &Driver{}
		err = conf.Get(prefix+"AppID", &c.AppID)
		if err != nil {
			return nil, err
		}
		err = conf.Get(prefix+"AppSecretKey", &c.AppSecretKey)
		if err != nil {
			return nil, err
		}
		err = conf.Get(prefix+"Clients", &c.Clients)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
}
