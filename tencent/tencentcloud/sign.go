package tencentcloud

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}
func HMAC_Sha256(key []byte, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

type SignData struct {
	HTTPRequestMethod    string
	CanonicalURI         string
	CanonicalQueryString string
	CanonicalHeaders     string
	SignedHeaders        string
	HashedRequestPayload string
	Service              string
	Timestamp            int64
}

func (d *SignData) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", d.HTTPRequestMethod, d.CanonicalURI, d.CanonicalQueryString, d.CanonicalHeaders, d.SignedHeaders, d.HashedRequestPayload)
}
func (d *SignData) ToSign(key string) string {
	date := d.Date()
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, d.Service)
	return fmt.Sprintf("%s\n%d\n%s\n%s", "TC3-HMAC-SHA256", d.Timestamp, credentialScope, sha256hex(d.String()))
}
func (d *SignData) Sign(key string) string {
	date := d.Date()
	tosign := d.ToSign(key)
	secretDate := HMAC_Sha256([]byte("TC3"+key), []byte(date))
	secretService := HMAC_Sha256(secretDate, []byte(d.Service))
	SecretSigning := HMAC_Sha256(secretService, []byte("tc3_request"))
	sign := strings.ToLower(hex.EncodeToString(HMAC_Sha256([]byte(SecretSigning), []byte(tosign))))
	return sign
}
func (d *SignData) Date() string {
	ts := time.Unix(d.Timestamp, 0)
	return ts.UTC().Format("2006-01-02")
}
func (d *SignData) Authorization(secretid string, key string) string {

	credentialScope := fmt.Sprintf("%s/%s/tc3_request", d.Date(), d.Service)
	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		"TC3-HMAC-SHA256",
		secretid,
		credentialScope,
		d.SignedHeaders,
		d.Sign(key),
	)
}
