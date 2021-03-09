package receiver

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/herb-go/providers/tencent/wechatmp"
)

type Receiver struct {
	App   *wechatmp.App
	Token string
}

type Message struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Event        *string
	EventKey     *string
	Latitude     *float64
	Longitude    *float64
	Precision    *float64
	Ticket       *string
	Content      *string
	MsgID        *int64
	PicURL       *string
	MediaID      *string
	Format       *string
	Recognition  *string
	ThumbMediaID *string
	LocationX    *float64 `xml:"Location_X"`
	LocationY    *float64 `xml:"Location_Y"`
	Scale        *int64
	Title        *string
	Description  *string
	URL          *string
}

type Handler func(r *http.Request, content []byte, msg *Message)

func (r *Receiver) Auth(q url.Values) (bool, error) {
	signature := q.Get("signature")
	timestamp := q.Get("timestamp")
	nonce := q.Get("nonce")
	if signature == "" || timestamp == "" || nonce == "" {
		return false, nil
	}
	l := sort.StringSlice([]string{r.Token, timestamp, nonce})
	l.Sort()
	data := strings.Join(l, "")
	h := sha1.New()
	h.Write([]byte(data))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash == signature, nil
}
func (r *Receiver) Middleware(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ok, err := r.Auth(req.URL.Query())
	if err != nil {
		panic(err)
	}
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	next(w, req)
}

func (r *Receiver) HandlerAction(handler Handler) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r.Handle(w, req, handler)
	}
}

func (r *Receiver) Handle(w http.ResponseWriter, req *http.Request, handler Handler) {
	q := req.URL.Query()
	ok, err := r.Auth(req.URL.Query())
	if err != nil {
		panic(err)
	}
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if req.Method == "GET" {
		if q.Get("echostr") == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		_, err = w.Write([]byte(q.Get("echostr")))
		if err != nil {
			panic(err)
		}
		return
	} else if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		msg := &Message{}
		err = xml.Unmarshal(body, msg)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		handler(req, body, msg)
		_, err = w.Write([]byte{})
		if err != nil {
			panic(err)
		}
		return
	} else {
		http.Error(w, http.StatusText(405), 405)
		return
	}
}
