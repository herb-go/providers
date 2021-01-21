package aliyun

import "github.com/herb-go/fetcher"

type AccessKey struct {
	AccessKeyID     string
	AccessKeySecret string
	Client          fetcher.Client
}
