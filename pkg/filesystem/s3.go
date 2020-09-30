package filesystem

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/fclairamb/afero-s3"
	"github.com/spf13/afero"
)

func NewS3Fs(region string, bucket string) afero.Fs {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	return s3.NewFs(bucket, sess)
}
