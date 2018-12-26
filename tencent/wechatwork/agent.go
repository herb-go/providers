package wechatwork

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"

	"github.com/herb-go/fetch"
)

type Agent struct {
	CorpID      string
	AgentID     int
	Secret      string
	Clients     fetch.Clients
	accessToken string
	lock        sync.Mutex
}

func (a *Agent) AccessToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.accessToken
}
func (a *Agent) NewMessage() *Message {
	return &Message{
		AgentID: a.AgentID,
	}
}
func (a *Agent) SendMessage(b *Message) (*MessageResult, error) {
	result := &MessageResult{}
	err := a.CallJSONApiWithAccessToken(apiMessagePost, nil, b, result)
	return result, err
}

func (a *Agent) GrantAccessToken() error {
	params := url.Values{}
	params.Set("corpid", a.CorpID)
	params.Set("corpsecret", a.Secret)
	req, err := apiGetToken.NewRequest(params, nil)
	if err != nil {
		return err
	}
	rep, err := a.Clients.Fetch(req)
	if err != nil {
		return err
	}
	if rep.StatusCode != http.StatusOK {
		return rep
	}
	result := &resultAccessToken{}
	err = rep.UnmarshalAsJSON(result)
	if err != nil {
		return err
	}
	if result.Errcode != 0 || result.Errmsg == "" || result.AccessToken == "" {
		return rep.NewAPICodeErr(result.Errcode)
	}
	a.accessToken = result.AccessToken
	return nil
}

func (a *Agent) CallJSONApiWithAccessToken(api *fetch.EndPoint, params url.Values, body interface{}, v interface{}) error {
	jsonAPIRequestBuilder := func(accesstoken string) (*http.Request, error) {
		p := url.Values{}
		if params != nil {
			for k, vs := range params {
				for _, v := range vs {
					p.Add(k, v)
				}
			}
		}
		p.Set("access_token", accesstoken)
		return api.NewJSONRequest(p, body)
	}
	return a.callApiWithAccessToken(api, jsonAPIRequestBuilder, v)
}
func (a *Agent) UploadApiWithAccessToken(api *fetch.EndPoint, params url.Values, filename string, body io.Reader, v interface{}) error {
	jsonAPIRequestBuilder := func(accesstoken string) (*http.Request, error) {
		p := url.Values{}
		if params != nil {
			for k, vs := range params {
				for _, v := range vs {
					p.Add(k, v)
				}
			}
		}
		p.Set("access_token", accesstoken)
		buffer := bytes.NewBuffer([]byte{})
		w := multipart.NewWriter(buffer)
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
		req, err := api.NewRequest(p, buffer.Bytes())
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contenttype)
		return req, nil
	}
	return a.callApiWithAccessToken(api, jsonAPIRequestBuilder, v)
}
func (a *Agent) callApiWithAccessToken(api *fetch.EndPoint, APIRequestBuilder func(accesstoken string) (*http.Request, error), v interface{}) error {
	var apierr resultAPIError
	var err error
	if a.AccessToken() == "" {
		err := a.GrantAccessToken()
		if err != nil {
			return err
		}
	}

	req, err := APIRequestBuilder(a.AccessToken())
	if err != nil {
		return err
	}
	resp, err := a.Clients.Fetch(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return resp
	}
	apierr = resultAPIError{}
	err = resp.UnmarshalAsJSON(&apierr)
	if err != nil {
		return err
	}
	if fetch.CompareAPIErrCode(err, ApiErrAccessTokenOutOfDate) || fetch.CompareAPIErrCode(err, ApiErrAccessTokenWrong) {
		err := a.GrantAccessToken()
		if err != nil {
			return err
		}
		req, err := APIRequestBuilder(a.AccessToken())
		if err != nil {
			return err
		}
		resp, err := a.Clients.Fetch(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return resp
		}
		apierr = resultAPIError{}
		err = resp.UnmarshalAsJSON(&apierr)
		if err != nil {
			return err
		}
	}
	if apierr.Errcode != 0 {
		return resp.NewAPICodeErr(apierr.Errcode)
	}
	return resp.UnmarshalAsJSON(&v)
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
		if fetch.CompareAPIErrCode(err, ApiErrUserUnaccessible) || fetch.CompareAPIErrCode(err, ApiErrNoPrivilege) {
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

func (a *Agent) MediaUpload(mediatype string, filename string, body io.Reader) (string, error) {
	params := url.Values{}
	params.Set("type", mediatype)
	result := &resultMediaUpload{}
	err := a.UploadApiWithAccessToken(apiMediaUpload, params, filename, body, result)
	if err != nil {
		return "", err
	}
	return result.MediaID, nil
}
