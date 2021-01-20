package tencentcloud

import "fmt"

type Error struct {
	Code    string
	Message string
}

type BaseResponse struct {
	Err       *Error `json:"Error"`
	RequestId string
}

func (r *BaseResponse) CodeError() error {
	if r.Err == nil {
		return nil
	}
	return r
}
func (r *BaseResponse) Error() string {
	return fmt.Sprint("tencentcloud apierror:%s - %s", r.Err.Code, r.Err.Message)
}
