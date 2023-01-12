package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Options struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `env:"S3_ACCESS_KEY"`
	SecretKey string `env:"S3_SECRET_KEY"`
	Region    string `yaml:"region"`
}

type Client interface {
	Put(bucket string, key string, body io.ReadSeeker) error
	Delete(bucket, key string) error
}

type client struct {
	s3Service *s3.S3
}

func (c *client) Put(bucket string, key string, body io.ReadSeeker) error {
	_, err := c.s3Service.PutObject(&s3.PutObjectInput{
		Body:   body,
		Key:    aws.String(key),
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Delete(bucket, key string) error {
	_, err := c.s3Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

func New(opts Options) (Client, error) {
	region := aws.String(opts.Region)
	url := aws.String(opts.Endpoint)
	creds := credentials.NewStaticCredentials(opts.AccessKey, opts.SecretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Endpoint:    url,
		Region:      region,
	})
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	return &client{
		s3Service: svc,
	}, nil
}
