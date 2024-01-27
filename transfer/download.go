package transfer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	myAws "github.com/halprin/cloud-backup/aws"
	"github.com/halprin/cloud-backup/config"
	"io"
)

type downloader struct {
	s3Client    *s3.Client
	s3Bucket    string
	s3ObjectKey string
}

func NewDownloader(overallConfig config.BackupConfiguration, timestamp string, backupFile string) (*downloader, error) {
	awsConfig, err := myAws.GetConfig(overallConfig.AwsCredentialConfigPath, overallConfig.AwsProfile)
	if err != nil {
		return nil, err
	}

	s3ObjectKey := fmt.Sprintf("%s/%s", timestamp, backupFile)

	return &downloader{
		s3Client:    s3.NewFromConfig(awsConfig),
		s3Bucket:    overallConfig.S3Bucket,
		s3ObjectKey: s3ObjectKey,
	}, nil
}

func (receiver *downloader) Download() (io.Reader, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: &receiver.s3Bucket,
		Key:    &receiver.s3ObjectKey,
	}

	getObjectOutput, err := receiver.s3Client.GetObject(context.Background(), getObjectInput)
	if err != nil {
		return nil, err
	}

	return getObjectOutput.Body, nil
}
