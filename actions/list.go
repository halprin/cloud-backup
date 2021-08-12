package actions

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	myAws "github.com/halprin/cloud-backup/aws"
	"github.com/halprin/cloud-backup/config"
	"log"
	"path"
)

func List(configFilePath string, timestamp string) error {
	log.Println("Listing backups")

	overallConfig, err := config.New(configFilePath)
	if err != nil {
		return err
	}

	if timestamp == "" {
		err = listTimestamps(overallConfig)
	} else {
		err = listBackups(overallConfig, timestamp)
	}

	log.Println("Done listing backups")
	return err
}

func listTimestamps(overallConfig config.BackupConfiguration) error {
	log.Println("Listing timestamps")

	awsSession, err := myAws.GetSession(overallConfig.AwsCredentialConfigPath, overallConfig.AwsProfile)
	if err != nil {
		return err
	}

	s3Client := s3.New(awsSession)

	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket:    &overallConfig.S3Bucket,
		Delimiter: aws.String("/"),
	}
	listObjectsOutput, err := s3Client.ListObjectsV2(listObjectsInput)
	if err != nil {
		return err
	}

	prefixes := listObjectsOutput.CommonPrefixes

	if len(prefixes) == 0 {
		log.Println("No backup timestamps")
		return nil
	}

	for _, prefix := range prefixes {
		prefixSansSlash := (*prefix.Prefix)[0:len(*prefix.Prefix) - 1]
		log.Printf("- %s", prefixSansSlash)
	}

	return nil
}

func listBackups(overallConfig config.BackupConfiguration, timestamp string) error {
	log.Printf("Listing backups in %s", timestamp)

	awsSession, err := myAws.GetSession(overallConfig.AwsCredentialConfigPath, overallConfig.AwsProfile)
	if err != nil {
		return err
	}

	s3Client := s3.New(awsSession)

	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: &overallConfig.S3Bucket,
		Prefix: &timestamp,
	}
	listObjectsOutput, err := s3Client.ListObjectsV2(listObjectsInput)
	if err != nil {
		return err
	}

	objects := listObjectsOutput.Contents

	if len(objects) == 0 {
		log.Println("No backups")
		return nil
	}

	for _, object := range objects {
		baseFilename := path.Base(*object.Key)
		log.Printf("- %s", baseFilename)
	}

	return nil
}