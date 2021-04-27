package s3manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"time"
)

type Uploader struct {
	s3Client            *s3.S3
	body                io.Reader
	multipartUploadData *s3.CreateMultipartUploadOutput
}

func NewUploader(session client.ConfigProvider) *Uploader {
	s3Client := s3.New(session)
	return &Uploader{
		s3Client: s3Client,
	}
}

func (receiver *Uploader) Upload(uploadInput *s3manager.UploadInput) error {
	dayLater := time.Now().AddDate(0, 0, 1)
	createUploadInput := &s3.CreateMultipartUploadInput{
		Bucket: uploadInput.Bucket,
		Key:    uploadInput.Key,
		Expires: &dayLater,
	}

	receiver.body = uploadInput.Body

	createUploadOutput, err := receiver.s3Client.CreateMultipartUpload(createUploadInput)
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

func (receiver *Uploader) uploadAllTheParts() ([]*s3.CompletedPart, error) {
	partSize := int64(5 * 1024 * 1024)  //start at 5 MB
	numberOfIterationsPerPartSize := 909
	partNumber := int64(1)
	var completedParts []*s3.CompletedPart

	for ; ; partSize*=2 {
		for partIndex := 0; partIndex < numberOfIterationsPerPartSize; partIndex++ {
			completedPart, err := receiver.uploadPart(partNumber, partSize)
			if err != nil {
				if err == io.EOF {
					completedParts = append(completedParts, completedPart)
					return completedParts, nil
				}
				return nil, err
			}

			completedParts = append(completedParts, completedPart)
		}
	}
}

func (receiver *Uploader) uploadPart(partNumber int64, partSize int64) (*s3.CompletedPart, error) {
	partInput := &s3.UploadPartInput{
		Body:          aws.ReadSeekCloser(receiver.body),
		Bucket:        receiver.multipartUploadData.Bucket,
		ContentLength: &partSize,
		Key:           receiver.multipartUploadData.Key,
		PartNumber:    &partNumber,
		UploadId:      receiver.multipartUploadData.UploadId,
	}

	partOutput, err := receiver.s3Client.UploadPart(partInput)
	if err != nil {
		//return the completed part anyway just in case the error is io.EOF
		return &s3.CompletedPart{
			ETag:       partOutput.ETag,
			PartNumber: &partNumber,
		}, err
	}

	return &s3.CompletedPart{
		ETag:       partOutput.ETag,
		PartNumber: &partNumber,
	}, nil
}

func (receiver *Uploader) finishMultipartUpload(completedParts []*s3.CompletedPart) error {
	completeUploadInput := &s3.CompleteMultipartUploadInput{
		Bucket:          receiver.multipartUploadData.Bucket,
		Key:             receiver.multipartUploadData.Key,
		UploadId:        receiver.multipartUploadData.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}

	completeUploadOutput, err := receiver.s3Client.CompleteMultipartUpload(completeUploadInput)
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

	_, _ = receiver.s3Client.AbortMultipartUpload(abortUploadInput)
	//swallow error because we are stopping because there was already an error and that is more important to report
}
