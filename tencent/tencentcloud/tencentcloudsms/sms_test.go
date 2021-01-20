package tencentcloudsms

import (
	"testing"
)

func TestSms(t *testing.T) {
	result, err := TestSMS.Send(NewTestMessage())
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal(result)
	}
}
