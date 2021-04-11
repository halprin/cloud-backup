package crypt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"log"
	"os"
)

var awsSession, sessionErr = session.NewSession()
var kmsService = kms.New(awsSession)

type dataKey struct {
	encryptedDataKey []byte
	plaintextDataKey []byte
}

func getEncryptionKey() (*dataKey, error) {
	if sessionErr != nil {
		log.Println("Initial AWS session failed")
		return nil, sessionErr
	}

	kmsKeyArn := os.Args[3]

	generateDataKeyInput := &kms.GenerateDataKeyInput{
		KeyId:             &kmsKeyArn,
		KeySpec:           aws.String(kms.DataKeySpecAes256),
		EncryptionContext: map[string]*string{
			"context": aws.String("a test context string"),
		},
	}

	generateDataKeyOutput, err := kmsService.GenerateDataKey(generateDataKeyInput)
	if err != nil {
		return nil, err
	}

	newDataKey := &dataKey{
		encryptedDataKey: generateDataKeyOutput.CiphertextBlob,
		plaintextDataKey: generateDataKeyOutput.Plaintext,
	}

	return newDataKey, nil
}