package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"github.com/halprin/cloud-backup-go/config"
	"io"
)

type encryptor struct {
	outputWriter            io.Writer
	config                  config.BackupConfiguration

	authenticatedEncryption cipher.AEAD
	encryptedDataKey        []byte
}

func NewEncryptor(outputWriter io.Writer, config config.BackupConfiguration) *encryptor {
	return &encryptor{
		outputWriter: outputWriter,
		config: config,
	}
}

func (receiver *encryptor) Write(plaintext []byte) (int, error) {
	err := receiver.encrypt(plaintext)
	if err != nil {
		return 0, err
	}

	return len(plaintext), nil
}

func (receiver *encryptor) encrypt(plaintext []byte) error {

	if receiver.authenticatedEncryption == nil || len(receiver.encryptedDataKey) == 0 {
		encryptionKey, err := getEncryptionKey(receiver.config.KmsKey, receiver.config.EncryptionContext, receiver.config.AwsProfile)
		if err != nil {
			return err
		}

		receiver.encryptedDataKey = encryptionKey.EncryptedDataKey

		authenticatedEncryption, err := createAuthenticatedEncryption(encryptionKey)
		if err != nil {
			return err
		}
		receiver.authenticatedEncryption = authenticatedEncryption

		clearPlaintextDataKey(encryptionKey)
	}

	nonce, err := receiver.generateNonce(receiver.authenticatedEncryption)
	if err != nil {
		return err
	}

	ciphertext := receiver.authenticatedEncryption.Seal(nil, nonce, plaintext, []byte(receiver.config.EncryptionContext))

	messageEnvelope := &envelope{
		Key:     receiver.encryptedDataKey,
		Nonce:   nonce,
		Message: ciphertext,
	}

	err = gob.NewEncoder(receiver.outputWriter).Encode(messageEnvelope)
	return err
}

func (receiver *encryptor) generateNonce(aead cipher.AEAD) ([]byte, error) {
	nonce := make([]byte, aead.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
