package aliyunsms

import "testing"

func TestSend(t *testing.T) {
	msg := NewTestMessage()
	result, err := Send(TestKey, msg)
	if result == nil || err != nil {
		t.Fatal(result, err)
	}
}
