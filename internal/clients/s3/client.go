package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Options struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `yaml:"use_ssl"`
}

type Client interface {
	//TODO methods for add and delete files
}

type client struct {
	minioClient *minio.Client
}

func New(opts Options) (Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(opts.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.AccessKeyID, opts.SecretAccessKey, ""),
		Secure: opts.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &client{
		minioClient: minioClient,
	}, nil
}
