package crypt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"log"
)

func getEncryptionKey(kmsKeyArn string, encryptionContext string, awsProfile string) (*dataKey, error) {
	kmsService, err := getKmsClient(awsProfile)
	if err != nil {
		return nil, err
	}

	generateDataKeyInput := &kms.GenerateDataKeyInput{
		KeyId:             &kmsKeyArn,
		KeySpec:           aws.String(kms.DataKeySpecAes256),
		EncryptionContext: map[string]*string{
			"context": &encryptionContext,
		},
	}

	generateDataKeyOutput, err := kmsService.GenerateDataKey(generateDataKeyInput)
	if err != nil {
		return nil, err
	}

	newDataKey := &dataKey{
		EncryptedDataKey: generateDataKeyOutput.CiphertextBlob,
		PlaintextDataKey: generateDataKeyOutput.Plaintext,
	}

	return newDataKey, nil
}

func getDecryptionKey(encryptedDataKey []byte, encryptionContext string, awsProfile string) (*dataKey, error) {
	kmsService, err := getKmsClient(awsProfile)
	if err != nil {
		return nil, err
	}

	decryptInput := &kms.DecryptInput{
		CiphertextBlob:    encryptedDataKey,
		EncryptionContext: map[string]*string{
			"context": &encryptionContext,
		},
	}

	decryptOutput, err := kmsService.Decrypt(decryptInput)
	if err != nil {
		return nil, err
	}

	newDataKey := dataKey{
		EncryptedDataKey: encryptedDataKey,
		PlaintextDataKey: decryptOutput.Plaintext,
	}

	return &newDataKey, nil
}

func getKmsClient(awsProfile string) (*kms.KMS, error) {
	awsSession, err := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
	})
	if err != nil {
		log.Println("Initial AWS session failed")
		return nil, err
	}

	kmsService := kms.New(awsSession)

	return kmsService, nil
}
