package aliyun

import (
	"testing"
)

func TestSign(t *testing.T) {
	accessKey := &AccessKey{
		AccessKeyID:     "testid",
		AccessKeySecret: "testsecret",
	}
	var wantedToSign = `GET&%2F&AccessKeyId%3Dtestid%26Action%3DDescribeRegions%26Format%3DXML%26SignatureMethod%3DHMAC-SHA1%26SignatureNonce%3D3ee8c1b8-83d3-44af-a94f-4e0ad82fd6cf%26SignatureVersion%3D1.0%26Timestamp%3D2016-02-23T12%253A46%253A24Z%26Version%3D2014-05-26`
	var wantedSigned = `OLeaidS1JvxuMvnyHOwuJ+uX5qY=`
	p := NewParams(accessKey)
	p.
		With("Timestamp", "2016-02-23T12:46:24Z").
		With("SignatureNonce", "3ee8c1b8-83d3-44af-a94f-4e0ad82fd6cf").
		With("Format", "XML").
		WithVersion("2014-05-26").
		WithAction("DescribeRegions")
	sp := p.ToSignParams()
	tosign := sp.ToSign("GET")
	if tosign != wantedToSign {
		t.Fatal(tosign)
	}
	signed := p.Sign("GET", accessKey.AccessKeySecret)
	if signed != wantedSigned {
		t.Fatal(signed)
	}
}
