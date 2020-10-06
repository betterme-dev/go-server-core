package filesystem

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/fclairamb/afero-s3"
	"github.com/spf13/afero"
)

type (
	S3Config struct {
		Bucket string
		Region string
	}
)

func NewS3Fs(conf S3Config) afero.Fs {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(conf.Region),
	})

	return s3.NewFs(conf.Bucket, sess)
}
