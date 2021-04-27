package transfer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/halprin/cloud-backup-go/config"
	"io"
	"log"
	"path"
	"sync"
)

type uploader struct {
	pipeWriter  *io.PipeWriter
	pipeReader  *io.PipeReader
	s3Uploader  *s3manager.Uploader
	uploadInput *s3manager.UploadInput
	waitGroup   sync.WaitGroup
}

func NewUploader(fileConfig config.BackupFileConfiguration, overallConfig config.BackupConfiguration, overallFolderName string) (*uploader, error) {
	pipeReader, pipeWriter := io.Pipe()

	s3Uploader, err := getUploader(overallConfig.AwsProfile)
	if err != nil {
		return nil, err
	}

	uploadInput := &s3manager.UploadInput{
		Bucket: &overallConfig.S3Bucket,
		Key:    aws.String(path.Join(overallFolderName, fileConfig.Title + ".cipher")),
		Body:   pipeReader,
	}

	newUploader := &uploader{
		pipeWriter: pipeWriter,
		pipeReader: pipeReader,
		s3Uploader: s3Uploader,
		uploadInput: uploadInput,
		waitGroup: sync.WaitGroup{},
	}

	newUploader.waitGroup.Add(1)
	go newUploader.initiateUpload()

	return newUploader, nil
}

func (receiver *uploader) Write(inputBytes []byte) (int, error) {
	return receiver.pipeWriter.Write(inputBytes)
}

func (receiver *uploader) Close() error {
	err := receiver.pipeWriter.Close()
	receiver.waitGroup.Wait()

	return err
}

func (receiver *uploader) initiateUpload() {
	_, err := receiver.s3Uploader.Upload(receiver.uploadInput)
	if err != nil {
		_ = receiver.pipeReader.CloseWithError(err)
	}

	receiver.waitGroup.Done()
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
	//allows up to 195.3125 GB
	newUploader.PartSize = 1024 * 1024 * 20

	return newUploader, nil
}
