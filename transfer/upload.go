package transfer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/halprin/cloud-backup-go/config"
	"io"
	"log"
	"path"
)

type uploader struct {
	overallFolderName string
	fileName   string
	pipeWriter *io.PipeWriter
	pipeReader *io.PipeReader
	s3Uploader *s3manager.Uploader
}

func NewUploader(fileConfig config.BackupFileConfiguration, overallConfig config.BackupConfiguration, overallFolderName string) (*uploader, error) {
	pipeReader, pipeWriter := io.Pipe()

	s3Uploader, err := getUploader(overallConfig.AwsProfile)
	if err != nil {
		return nil, err
	}

	upParams := &s3manager.UploadInput{
		Bucket: &overallConfig.S3Bucket,
		Key:    aws.String(path.Join(overallFolderName, fileConfig.Title + ".cipher")),
		Body:   pipeReader,
	}

	_, err = s3Uploader.Upload(upParams)
	if err != nil {
		return nil, err
	}

	return &uploader{
		overallFolderName: overallFolderName,
		fileName: fileConfig.Title,
		pipeWriter: pipeWriter,
		pipeReader: pipeReader,
		s3Uploader: s3Uploader,
	}, nil
}

func (receiver *uploader) Write(inputBytes []byte) (int, error) {
	return receiver.pipeWriter.Write(inputBytes)
}

func getUploader(awsProfile string) (*s3manager.Uploader, error) {
	awsSession, err := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
	})
	if err != nil {
		log.Println("Initial AWS session failed")
		return nil, err
	}

	newUploader := s3manager.NewUploader(awsSession)
	return newUploader, nil
}
