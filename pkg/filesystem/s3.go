package filesystem

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

var s3sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String(viper.GetString("AWS_REGION")),
}))

func DeleteS3File(bucket string, filename string) error {
	svc := s3.New(s3sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}
	_, err := svc.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}

func DownloadS3File(bucket string, filename string) ([]byte, error) {
	downloader := s3manager.NewDownloader(s3sess)
	buff := &aws.WriteAtBuffer{}

	_, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func UploadS3File(bucket string, filename string, content string) error {
	uploader := s3manager.NewUploader(s3sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   strings.NewReader(content),
	})
	if err != nil {
		return err
	}

	return nil
}
