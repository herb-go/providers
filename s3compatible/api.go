package s3compatible

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type API struct {
	S3 *s3.Client
}

func (a *API) Load(ctx context.Context, bucket string, key string, writer io.Writer) (int64, error) {
	result, err := a.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return 0, err
	}
	defer result.Body.Close()
	return io.Copy(writer, result.Body)
}
func (a *API) Save(ctx context.Context, bucket string, key string, reader io.Reader) error {
	_, err := a.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   reader,
	})
	return err
}

type Info struct {
	LastModified *time.Time
	Size         int64
}

func (a *API) Info(ctx context.Context, bucket string, key string) (*Info, error) {
	result, err := a.S3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	info := &Info{
		LastModified: result.LastModified,
		Size:         result.ContentLength,
	}
	return info, nil
}

func (a *API) Remove(ctx context.Context, bucket string, key string) error {
	_, err := a.S3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	return err
}

func (a *API) PresignGetObject(ctx context.Context, bucket string, key string, ttl time.Duration) (string, error) {
	psClient := s3.NewPresignClient(a.S3)
	result, err := psClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func (a *API) PresignPutObject(ctx context.Context, bucket string, key string, ttl time.Duration) (string, error) {
	psClient := s3.NewPresignClient(a.S3)
	result, err := psClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func New() *API {
	return &API{}
}
