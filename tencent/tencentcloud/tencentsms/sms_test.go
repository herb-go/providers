package tencentsms

import (
	"testing"
)

func TestSms(t *testing.T) {
	result, err := TestSMS.Send(NewTestMessage())
	if err != nil || result == nil {
		t.Fatal(result, err)
	}
}
