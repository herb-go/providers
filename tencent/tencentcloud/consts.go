package tencentcloud

type ContentType string

const (
	ContentTypeURLEncoded = ContentType("application/x-www-form-urlencoded")
	ContentTypeJSON       = ContentType("application/json; charset=utf-8")
	ContentTypeFormData   = ContentType("multipart/form-data")
)
