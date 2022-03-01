package s3compatible

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

var Debug = false

var DebugWriter = os.Stdout

type S3Config struct {
	Region       string
	Endpoint     string
	SecretID     string
	Secret       string
	SessionToken string
	UsePathStyle bool
}

func (c *S3Config) ApplyTo(api *API) error {
	cfg := aws.Config{
		Region:      c.Region,
		Credentials: credentials.NewStaticCredentialsProvider(c.SecretID, c.Secret, c.SessionToken),
	}
	if Debug {
		cfg.ClientLogMode = aws.LogRequest | aws.LogResponse
		cfg.Logger = logging.NewStandardLogger(DebugWriter)
	}
	if c.Endpoint != "" {
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               c.Endpoint, // or where ever you ran minio
				SigningRegion:     c.Region,
				HostnameImmutable: true,
			}, nil
		})
	}
	opts := []func(*s3.Options){
		func(opt *s3.Options) {
			opt.UsePathStyle = c.UsePathStyle
		},
	}
	api.S3 = s3.NewFromConfig(cfg, opts...)
	return nil
}
