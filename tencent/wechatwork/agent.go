package wechatwork

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/url"
	"sync"

	"github.com/herb-go/fetcher"
	"github.com/herb-go/remoteprocedure/fetcherapi/sharedrefresherapi"
)

type Agent struct {
	CorpID               string
	AgentID              int
	Secret               string
	RemoteRefresher      *fetcher.Server
	Client               fetcher.Client
	accessToken          string
	lock                 sync.Mutex
	accessTokenRefresher func(string) (string, error)
	accessTokenGetter    func() (string, error)
}

//RefreshShared refresh shared data.
//New data what different from old should be returned
func (a *Agent) RefreshShared(old []byte) ([]byte, error) {
	var t string
	var err error
	a.lock.Lock()
	defer a.lock.Unlock()
	oldtoken := string(old)
	t, err = a.loadAccessToken()
	if err != nil {
		return nil, err
	}
	if t != "" && oldtoken != t {
		return []byte(t), nil
	}
	if a.accessTokenRefresher == nil {
		t, err = a.GetAccessToken()
	} else {
		t, err = a.accessTokenRefresher(string(old))
	}
	if err != nil {
		return nil, err
	}
	return []byte(t), nil
}
func (a *Agent) loadAccessToken() (string, error) {
	if a.accessTokenGetter != nil {
		return a.accessTokenGetter()
	}
	return a.accessToken, nil
}

func (a *Agent) AccessToken() (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.loadAccessToken()
}
func (a *Agent) SetAccessTokenGetter(f func() (string, error)) {
	a.accessTokenGetter = f
}

func (a *Agent) SetAccessTokenRefresher(f func(string) (string, error)) {
	a.accessTokenRefresher = f
}
func (a *Agent) NewMessage() *Message {
	return &Message{
		AgentID: a.AgentID,
	}
}
func (a *Agent) SendMessage(b *Message) (*MessageResult, error) {
	result := &MessageResult{}
	if b.AgentID == 0 {
		b.AgentID = a.AgentID
	}
	err := a.CallJSONApiWithAccessToken(apiMessagePost, nil, b, result)
	return result, err
}

func (a *Agent) ClientCredentialBuilder() fetcher.Command {
	return fetcher.ParamsBuilderFunc(func(params url.Values) error {
		params.Set("corpid", a.CorpID)
		params.Set("corpsecret", a.Secret)
		return nil
	})
}
func (a *Agent) getRemoteToken() (string, error) {
	t, err := a.loadAccessToken()
	if err != nil {
		return "", err
	}
	data, err := sharedrefresherapi.FetcherRefreshShared(a.RemoteRefresher, []byte(t))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (a *Agent) GetAccessToken() (string, error) {
	if !a.RemoteRefresher.IsEmpty() {
		return a.getRemoteToken()
	}
	result := &resultAccessToken{}
	resp, err := fetcher.DoAndParse(
		&a.Client,
		apiGetToken.With(a.ClientCredentialBuilder()),
		fetcher.Should200(fetcher.AsJSON(result)),
	)
	if err != nil {
		return "", err
	}
	if result.Errcode != 0 || result.Errmsg == "" || result.AccessToken == "" {
		return "", resp.NewAPICodeErr(result.Errcode)
	}
	return result.AccessToken, nil
}
func (a *Agent) GrantAccessToken() (string, error) {
	var token string
	var err error
	a.lock.Lock()
	defer a.lock.Unlock()
	token, err = a.loadAccessToken()
	if err != nil {
		return "", err
	}
	if a.accessTokenRefresher == nil {
		token, err = a.GetAccessToken()
	} else {
		token, err = a.accessTokenRefresher(token)
	}

	if err != nil {
		return "", err
	}
	a.accessToken = token
	return token, nil
}

func (a *Agent) CallJSONApiWithAccessToken(api *fetcher.Preset, params url.Values, body interface{}, v interface{}) error {
	jsonAPIRequestBuilder := func(accesstoken string) (*fetcher.Preset, error) {
		return api.With(fetcher.Params(params), fetcher.SetQuery("access_token", accesstoken), fetcher.JSONBody(body)), nil
	}
	return a.callApiWithAccessToken(api, jsonAPIRequestBuilder, v)
}
func (a *Agent) UploadApiWithAccessToken(api *fetcher.Preset, params url.Values, filename string, body io.Reader, v interface{}) error {
	jsonAPIRequestBuilder := func(accesstoken string) (*fetcher.Preset, error) {
		buffer := bytes.NewBuffer([]byte{})
		w := multipart.NewWriter(buffer)
		w.WriteField("type", params.Get("type"))
		defer w.Close()
		filewriter, err := w.CreateFormFile("media", filename)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(filewriter, body)
		if err != nil {
			return nil, err
		}
		err = w.Close()
		if err != nil {
			return nil, err
		}
		contenttype := w.FormDataContentType()
		return api.With(fetcher.SetQuery("access_token", accesstoken), fetcher.SetHeader("Content-Type", contenttype), fetcher.Body(buffer)), nil
	}
	return a.callApiWithAccessToken(api, jsonAPIRequestBuilder, v)
}
func (a *Agent) callApiWithAccessToken(api *fetcher.Preset, APIPresetBuilder func(accesstoken string) (*fetcher.Preset, error), v interface{}) error {
	var apierr = &resultAPIError{}
	var err error
	token, err := a.AccessToken()
	if err != nil {
		return err
	}

	if token == "" {
		token, err = a.GrantAccessToken()
		if err != nil {
			return err
		}
	}

	preset, err := APIPresetBuilder(token)
	if err != nil {
		return err
	}

	resp, err := fetcher.DoAndParse(&a.Client, preset, fetcher.Should200(fetcher.AsJSON(apierr)))
	if err != nil {
		return err
	}
	if !apierr.IsOK() {
		if apierr.IsAccessTokenError() {
			token, err = a.GrantAccessToken()
			if err != nil {
				return err
			}
			apierr = &resultAPIError{}
			resp, err = fetcher.DoAndParse(&a.Client, preset, fetcher.Should200(fetcher.AsJSON(apierr)))
			if err != nil {
				return err
			}
			if !apierr.IsOK() {
				return resp.NewAPICodeErr(apierr.Errcode)
			}

		} else {
			return resp.NewAPICodeErr(apierr.Errcode)
		}
	}

	return fetcher.AsJSON(v).Parse(resp)
}

type Userinfo struct {
	UserID     string
	Name       string
	Mobile     string
	Email      string
	Gender     string
	Avatar     string
	Department []int
}

func (a *Agent) GetUserInfo(code string) (*Userinfo, error) {
	var info = &Userinfo{}
	if code == "" {
		return nil, nil
	}
	var result = &resultUserInfo{}
	params := url.Values{}
	params.Set("code", code)
	err := a.CallJSONApiWithAccessToken(apiGetUserInfo, params, nil, result)
	if err != nil {
		return nil, err
	}
	if result.UserID == "" {
		return nil, nil
	}
	var getuser = &resultUserGet{}
	userGetParam := url.Values{}
	userGetParam.Add("userid", result.UserID)
	err = a.CallJSONApiWithAccessToken(apiUserGet, userGetParam, nil, getuser)
	if err != nil {
		if fetcher.CompareAPIErrCode(err, APIErrUserUnaccessible) || fetcher.CompareAPIErrCode(err, APIErrNoPrivilege) {
			return nil, nil
		}
		return nil, err
	}
	info.UserID = result.UserID
	info.Avatar = getuser.Avatar
	info.Email = getuser.Email
	info.Gender = getuser.Gender
	info.Mobile = getuser.Mobile
	info.Name = getuser.Name
	info.Department = getuser.Department
	return info, nil
}

func (a *Agent) GetDepartmentList(id string) (*[]DepartmentInfo, error) {
	params := url.Values{}
	if id != "" {
		params.Set("id", id)
	}
	var result = &resultDepartmentList{}
	err := a.CallJSONApiWithAccessToken(apiDepartmentList, params, nil, result)
	if err != nil {
		return nil, err
	}
	return result.Department, nil
}

func (a *Agent) MediaUpload(mediatype MediaType, filename string, body io.Reader) (string, error) {
	params := url.Values{}
	params.Set("type", string(mediatype))
	result := &resultMediaUpload{}
	err := a.UploadApiWithAccessToken(apiMediaUpload, params, filename, body, result)
	if err != nil {
		return "", err
	}
	return result.MediaID, nil
}
