package tencentminiprogramum

import (
	"testing"
)

func TestSend(t *testing.T) {
	app := TestApp
	msg := NewTestMessage()
	t.Fatal(Send(app, msg))
}
