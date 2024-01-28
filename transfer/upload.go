package transfer

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	myAws "github.com/halprin/cloud-backup/aws"
	"github.com/halprin/cloud-backup/aws/myS3Manager"
	"github.com/halprin/cloud-backup/config"
	"io"
	"path"
	"sync"
)

type uploader struct {
	pipeWriter  *io.PipeWriter
	pipeReader  *io.PipeReader
	s3Uploader  *myS3Manager.Uploader
	uploadInput *s3.PutObjectInput
	waitGroup   sync.WaitGroup
}

func NewUploader(fileConfig config.BackupFileConfiguration, overallConfig config.BackupConfiguration, overallFolderName string) (*uploader, error) {
	pipeReader, pipeWriter := io.Pipe()

	s3Uploader, err := getUploader(overallConfig)
	if err != nil {
		return nil, err
	}

	uploadInput := &s3.PutObjectInput{
		Bucket: &overallConfig.S3Bucket,
		Key:    aws.String(path.Join(overallFolderName, fileConfig.Title+".tar.gz")),
		Body:   pipeReader,
	}

	if overallConfig.KmsKey != "" {
		uploadInput.ServerSideEncryption = types.ServerSideEncryptionAwsKms
		uploadInput.SSEKMSKeyId = &overallConfig.KmsKey
		uploadInput.SSEKMSEncryptionContext = aws.String(base64.StdEncoding.EncodeToString([]byte(overallConfig.EncryptionContext)))
	}

	newUploader := &uploader{
		pipeWriter:  pipeWriter,
		pipeReader:  pipeReader,
		s3Uploader:  s3Uploader,
		uploadInput: uploadInput,
		waitGroup:   sync.WaitGroup{},
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
	err := receiver.s3Uploader.Upload(receiver.uploadInput)
	if err != nil {
		_ = receiver.pipeReader.CloseWithError(err)
	}

	receiver.waitGroup.Done()
}

func getUploader(overallConfig config.BackupConfiguration) (*myS3Manager.Uploader, error) {
	awsConfig, err := myAws.GetConfig(overallConfig.AwsCredentialConfigPath, overallConfig.AwsProfile)
	if err != nil {
		return nil, err
	}

	newUploader := myS3Manager.NewUploader(awsConfig)

	return newUploader, nil
}
