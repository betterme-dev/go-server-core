package filesystem

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/fclairamb/afero-s3"
	"github.com/spf13/afero"
)

type (
	S3FsFactory struct {
		Bucket string
		Region string
	}
)

func (sfs S3FsFactory) New() afero.Fs {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(sfs.Region),
	})

	return s3.NewFs(sfs.Bucket, sess)
}
