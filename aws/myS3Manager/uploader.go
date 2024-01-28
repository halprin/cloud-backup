package myS3Manager

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/halprin/cloud-backup/parallel"
	"io"
	"time"
)

type Uploader struct {
	s3Client            *s3.Client
	body                io.Reader
	multipartUploadData *s3.CreateMultipartUploadOutput
}

func NewUploader(awsConfig aws.Config) *Uploader {
	s3Client := s3.NewFromConfig(awsConfig)
	return &Uploader{
		s3Client: s3Client,
	}
}

func (receiver *Uploader) Upload(uploadInput *s3.PutObjectInput) error {
	dayLater := time.Now().AddDate(0, 0, 1)
	createUploadInput := &s3.CreateMultipartUploadInput{
		Bucket:                  uploadInput.Bucket,
		Key:                     uploadInput.Key,
		Expires:                 &dayLater,
		ServerSideEncryption:    uploadInput.ServerSideEncryption,
		SSEKMSKeyId:             uploadInput.SSEKMSKeyId,
		SSEKMSEncryptionContext: uploadInput.SSEKMSEncryptionContext,
	}

	receiver.body = uploadInput.Body

	createUploadOutput, err := receiver.s3Client.CreateMultipartUpload(context.Background(), createUploadInput)
	if err != nil {
		return err
	}

	receiver.multipartUploadData = createUploadOutput

	completedParts, err := receiver.uploadAllTheParts()
	if err != nil {
		receiver.stopMultipartUpload()
		return err
	}

	err = receiver.finishMultipartUpload(completedParts)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Uploader) uploadAllTheParts() ([]types.CompletedPart, error) {
	partSize := int64(5 * 1024 * 1024) //start at 5 MB
	numberOfIterationsPerPartSize := 909
	partNumber := int32(1)
	var completedPartChannels []chan types.CompletedPart
	var errorChannels []chan error

	poolSize := 5
	taskQueueSize := poolSize * 2
	pool := parallel.NewPool(poolSize, taskQueueSize)
	defer pool.Release()

	for ; ; partSize *= 2 {
		for partIndex := 0; partIndex < numberOfIterationsPerPartSize; partIndex++ {

			partBytes, err := receiver.readPart(partSize)
			if err != nil {
				if err == io.EOF {
					//we're done reading, check the upload error channels first before we call this a success
					err := returnFirstErrorInSlice(parallel.ConvertChannelsOfErrorToErrorSlice(errorChannels))
					if err != nil {
						return nil, err
					}
					//no upload errors, so return the completed parts
					return parallel.ConvertChannelsOfCompletedPartsToSlice(completedPartChannels), nil
				}
				return nil, err
			}

			//check upload errors every taskQueueSize times
			if partNumber%int32(taskQueueSize) == 0 {
				err := returnFirstErrorInSlice(parallel.ConvertChannelsOfErrorToErrorSlice(errorChannels))
				if err != nil {
					return nil, err
				}
				errorChannels = make([]chan error, 0, taskQueueSize)
			}

			errorChannel := make(chan error, 1)
			errorChannels = append(errorChannels, errorChannel)
			completedPartChannel := make(chan types.CompletedPart, 1)
			completedPartChannels = append(completedPartChannels, completedPartChannel)

			func(partBytes []byte, partNumber int32, completedPartChannel chan types.CompletedPart, errorChannel chan error) {
				pool.Submit(func() {
					completedPart, err := receiver.uploadPart(partBytes, partNumber)
					if err != nil {
						errorChannel <- err
						close(errorChannel)
						close(completedPartChannel)
						return
					}
					completedPartChannel <- completedPart
					close(errorChannel)
					close(completedPartChannel)
				})
			}(partBytes, partNumber, completedPartChannel, errorChannel) //copy partBytes and partNumber so it is unique for the closure

			partNumber++
		}
	}
}

func (receiver *Uploader) readPart(partSize int64) ([]byte, error) {
	//read up to partSize amount of bytes
	fullPartBytes := make([]byte, partSize)

	partBytesRead, err := io.ReadFull(receiver.body, fullPartBytes) //ReadFull to try to fill the full size of partSize
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	//possible that the last read returns less than a full partSize, so we need to slice the bytes to the read size
	partBytes := fullPartBytes[:partBytesRead]

	return partBytes, nil
}

func (receiver *Uploader) uploadPart(partBytes []byte, partNumber int32) (types.CompletedPart, error) {
	md5Hash := calculateMd5Hash(partBytes)

	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(partBytes),
		Bucket:        receiver.multipartUploadData.Bucket,
		ContentLength: aws.Int64(int64(len(partBytes))),
		ContentMD5:    aws.String(md5Hash),
		Key:           receiver.multipartUploadData.Key,
		PartNumber:    &partNumber,
		UploadId:      receiver.multipartUploadData.UploadId,
	}

	partOutput, err := receiver.s3Client.UploadPart(context.Background(), partInput)
	if err != nil {
		return types.CompletedPart{}, err
	}

	return types.CompletedPart{
		ETag:       partOutput.ETag,
		PartNumber: &partNumber,
	}, nil
}

func (receiver *Uploader) finishMultipartUpload(completedParts []types.CompletedPart) error {
	completeUploadInput := &s3.CompleteMultipartUploadInput{
		Bucket:   receiver.multipartUploadData.Bucket,
		Key:      receiver.multipartUploadData.Key,
		UploadId: receiver.multipartUploadData.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}

	_, err := receiver.s3Client.CompleteMultipartUpload(context.Background(), completeUploadInput)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Uploader) stopMultipartUpload() {
	abortUploadInput := &s3.AbortMultipartUploadInput{
		Bucket:   receiver.multipartUploadData.Bucket,
		Key:      receiver.multipartUploadData.Key,
		UploadId: receiver.multipartUploadData.UploadId,
	}

	_, _ = receiver.s3Client.AbortMultipartUpload(context.Background(), abortUploadInput)
	//swallow error because we are stopping because there was already an error and that is more important to report
}

func calculateMd5Hash(partBytes []byte) string {
	md5Algorithm := md5.New()
	md5Algorithm.Write(partBytes)
	md5Hash := base64.StdEncoding.EncodeToString(md5Algorithm.Sum(nil))
	return md5Hash
}

func returnFirstErrorInSlice(errorSlice []error) error {
	for _, err := range errorSlice {
		if err != nil {
			return err
		}
	}

	return nil
}
