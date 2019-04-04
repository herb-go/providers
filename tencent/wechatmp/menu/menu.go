package menu

// Button wechat mp menu button struct
type Button struct {
	Name      string       `json:"name"`
	SubButton []*SubButton `json:"subbutton"`
	Type      *string      `json:"type"`
	Key       *string      `json:"key"`
	URL       *string      `json:"url"`
	MediaID   *string      `json:"media_id"`
	AppID     *string      `json:"appid"`
	Pagepath  *string      `json:"pagepath"`
}

// SubButton wechat mp menu subbutton struct
type SubButton struct {
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Key      *string `json:"key"`
	URL      *string `json:"url"`
	MediaID  *string `json:"media_id"`
	AppID    *string `json:"appid"`
	Pagepath *string `json:"pagepath"`
}

// Menu wechat mp menu struct
type Menu struct {
	Button []*Button `json:"button"`
}

func (m *Menu) NewButton() *Button {
	b := &Button{
		SubButton: []*SubButton{},
	}
	m.Button = append(m.Button, b)
	return b
}
func New() *Menu {
	menu := Menu{
		Button: []*Button{},
	}
	return &menu
}
