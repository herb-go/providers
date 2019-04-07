package menu

// Button wechat mp menu button struct
type Button struct {
	Name      string       `json:"name"`
	SubButton []*SubButton `json:"sub_button"`
	Type      *string      `json:"type"`
	Key       *string      `json:"key"`
	URL       *string      `json:"url"`
	MediaID   *string      `json:"media_id"`
	AppID     *string      `json:"appid"`
	Pagepath  *string      `json:"pagepath"`
}

func (b *Button) Validate() (string, error) {
	if b.Name == "" {
		return "按钮名不能为空", nil
	}
	if len(b.SubButton) > 5 {
		return "子菜单不能超过5个", nil
	}
	for _, v := range b.SubButton {
		r, err := v.Validate()
		if err != nil {
			return "", err
		}
		if r != "" {
			return r, nil
		}
	}
	return "", nil
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

func (b *SubButton) Validate() (string, error) {
	if b.Name == "" {
		return "子菜单名不能为空", nil
	}
	return "", nil
}

// Menu wechat mp menu struct
type Menu struct {
	Button []*Button `json:"button"`
}

func (m *Menu) Validate() (string, error) {
	if len(m.Button) > 3 {
		return "菜单不能超过3个", nil
	}
	for _, v := range m.Button {
		r, err := v.Validate()
		if err != nil {
			return "", err
		}
		if r != "" {
			return r, nil
		}
	}
	return "", nil

}
func (m *Menu) NewButton() *Button {
	b := &Button{
		SubButton: []*SubButton{},
	}
	m.Button = append(m.Button, b)
	return b
}

type MenuResult struct {
	Menu *Menu `json:"menu"`
}

func New() *Menu {
	menu := Menu{
		Button: []*Button{},
	}
	return &menu
}

func NewMenuResult() *MenuResult {
	return &MenuResult{
		Menu: New(),
	}
}
