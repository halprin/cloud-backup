package crypt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	myAws "github.com/halprin/cloud-backup-go/aws"
)

func getEncryptionKey(kmsKeyArn string, encryptionContext string, awsCredentialConfigPath string, awsProfile string) (*dataKey, error) {
	kmsService, err := getKmsClient(awsCredentialConfigPath, awsProfile)
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

func getDecryptionKey(encryptedDataKey []byte, encryptionContext string, awsCredentialConfigPath string, awsProfile string) (*dataKey, error) {
	kmsService, err := getKmsClient(awsCredentialConfigPath, awsProfile)
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

func getKmsClient(awsCredentialConfigPath string, awsProfile string) (*kms.KMS, error) {
	awsSession, err := myAws.GetSession(awsCredentialConfigPath, awsProfile)
	if err != nil {
		return nil, err
	}

	kmsService := kms.New(awsSession)

	return kmsService, nil
}
