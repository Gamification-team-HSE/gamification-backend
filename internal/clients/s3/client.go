package s3

import (
	"fmt"
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
	Bucket    string `yaml:"bucket"`
}

type Client interface {
	Put(folder string, filename string, body io.ReadSeeker) error
	Delete(folder, filename string) error
	BuildURL(folder, filename string) string
}

type client struct {
	s3Service *s3.S3

	opts Options
}

func (c *client) BuildURL(folder, filename string) string {
	return fmt.Sprintf("https://%s.%s/%s/%s", c.opts.Bucket, c.opts.Endpoint, folder, filename)
}

func key(folder, filename string) string {
	return fmt.Sprintf("%s/%s", folder, filename)
}

func (c *client) Put(folder string, filename string, body io.ReadSeeker) error {
	_, err := c.s3Service.PutObject(&s3.PutObjectInput{
		Body:   body,
		Key:    aws.String(key(folder, filename)),
		Bucket: aws.String(c.opts.Bucket),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Delete(folder, filename string) error {
	_, err := c.s3Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.opts.Bucket),
		Key:    aws.String(key(folder, filename)),
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
		LogLevel:    aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	return &client{
		s3Service: svc,
		opts:      opts,
	}, nil
}
