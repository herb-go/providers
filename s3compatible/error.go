package s3compatible

import (
	"errors"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
)

func IsHTTPError(err error, statuscode int) bool {
	if err != nil {
		httpErr := &awshttp.ResponseError{}
		if errors.As(err, &httpErr) {
			return httpErr.HTTPStatusCode() == statuscode
		}
	}
	return false
}
