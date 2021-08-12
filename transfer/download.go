package transfer

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	myAws "github.com/halprin/cloud-backup/aws"
	"github.com/halprin/cloud-backup/config"
	"io"
)

type downloader struct {
	s3Client    *s3.S3
	s3Bucket    string
	s3ObjectKey string
}

func NewDownloader(overallConfig config.BackupConfiguration, timestamp string, backupFile string) (*downloader, error) {
	awsSession, err := myAws.GetSession(overallConfig.AwsCredentialConfigPath, overallConfig.AwsProfile)
	if err != nil {
		return nil, err
	}

	s3ObjectKey := fmt.Sprintf("%s/%s", timestamp, backupFile)

	return &downloader{
		s3Client:    s3.New(awsSession),
		s3Bucket:    overallConfig.S3Bucket,
		s3ObjectKey: s3ObjectKey,
	}, nil
}

func (receiver *downloader) Download() (io.Reader, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: &receiver.s3Bucket,
		Key:    &receiver.s3ObjectKey,
	}

	getObjectOutput, err := receiver.s3Client.GetObject(getObjectInput)
	if err != nil {
		return nil, err
	}

	return getObjectOutput.Body, nil
}
