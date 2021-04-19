package crypt

import (
	"encoding/gob"
	"github.com/halprin/cloud-backup-go/config"
	"io"
)

type decryptor struct {
	inputReader  io.Reader
	outputWriter io.Writer
	config config.BackupConfiguration
}

func NewDecryptor(inputReader io.Reader, outputWriter io.Writer, config config.BackupConfiguration) *decryptor {
	return &decryptor{
		inputReader: inputReader,
		outputWriter: outputWriter,
		config: config,
	}
}

func (receiver *decryptor) Decrypt() error {
	for {
		var messageEnvelope envelope

		gobDecoder := gob.NewDecoder(receiver.inputReader)
		err := gobDecoder.Decode(&messageEnvelope)
		if err == io.EOF {
			//we've exhausted the reader; end decrypting successfully
			return nil
		} else if err != nil {
			return err
		}

		decryptionKey, err := getDecryptionKey(messageEnvelope.Key, receiver.config.EncryptionContext, receiver.config.AwsProfile)
		if err != nil {
			return err
		}

		authenticatedEncryption, err := createAuthenticatedEncryption(decryptionKey)
		if err != nil {
			return err
		}

		clearPlaintextDataKey(decryptionKey)

		plaintext, err := authenticatedEncryption.Open(nil, messageEnvelope.Nonce, messageEnvelope.Message, []byte(receiver.config.EncryptionContext))
		if err != nil {
			return err
		}

		_, err = receiver.outputWriter.Write(plaintext)
		if err != nil {
			return err
		}
	}
}
