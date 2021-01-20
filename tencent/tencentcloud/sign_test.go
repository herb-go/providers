package tencentcloud

import (
	"net/url"
	"testing"
)

func TestSign(t *testing.T) {
	secretid := "AKIDz8krbsJ5yKBZQpn74WFkmLPx3*******"
	key := "Gu5t9xGARNpq86cd98joQYCN3*******"
	wanted := "TC3-HMAC-SHA256 Credential=AKIDz8krbsJ5yKBZQpn74WFkmLPx3*******/2019-02-25/cvm/tc3_request, SignedHeaders=content-type;host, Signature=2230eefd229f582d8b1b891af7107b91597240707d778ab3738f756258d7652c"
	wantedrequest := `POST
/

content-type:application/json; charset=utf-8
host:cvm.tencentcloudapi.com

content-type;host
35e9c5b0e3ae67532d3c9f17ead6c90222632e5b1ff7f6e89887f1398934f064`
	wantedToSign := `TC3-HMAC-SHA256
1551113065
2019-02-25/cvm/tc3_request
5ffe6a04c0664d6b969fab9a13bdab201d63ee709638e2749d62a09ca18d7031`
	request := NewRequest()
	request.URL, _ = url.Parse("https://cvm.tencentcloudapi.com")
	request.Header.Set("Host", "cvm.tencentcloudapi.com")
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Action = "DescribeInstances"
	request.Service = "cvm"
	request.Timestamp = 1551113065
	request.Version = "2017-03-12"
	request.Region = "ap-guangzhou"
	request.Body = []byte(`{"Limit": 1, "Filters": [{"Values": ["\u672a\u547d\u540d"], "Name": "instance-name"}]}`)
	request.Method = "POST"
	data := request.CreateSignData(nil)
	if data.String() != wantedrequest {
		t.Fatal(data.String())
	}
	if data.ToSign(key) != wantedToSign {
		t.Fatal(data.ToSign(key))
	}
	auth := data.Authorization(secretid, key)
	if auth != wanted {
		t.Fatal(auth)
	}
}
