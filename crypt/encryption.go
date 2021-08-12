package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"github.com/halprin/cloud-backup/config"
	"github.com/halprin/cloud-backup/external/pb"
	"io"
)

type encryptor struct {
	outputWriter     io.Writer
	config           config.BackupConfiguration
	encoderInterface EnvelopeEncryptionWriter

	authenticatedEncryption cipher.AEAD
	encryptedDataKey        []byte

	encryptedDataKeyWritten bool
}

func NewEncryptor(outputWriter io.Writer, config config.BackupConfiguration) *encryptor {
	return &encryptor{
		outputWriter:     outputWriter,
		config:           config,
		encoderInterface: &pb.ProtoBufEnvelopeEncryptionWriter{},
	}
}

func (receiver *encryptor) Write(plaintext []byte) (int, error) {
	err := receiver.setUpEncryptionOnce()
	if err != nil {
		return 0, err
	}

	err = receiver.writeEncryptedDataKeyOnce()
	if err != nil {
		return 0, err
	}

	err = receiver.encrypt(plaintext)
	if err != nil {
		return 0, err
	}

	return len(plaintext), nil
}

func (receiver *encryptor) encrypt(plaintext []byte) error {
	nonce, err := receiver.generateNonce()
	if err != nil {
		return err
	}

	ciphertext := receiver.authenticatedEncryption.Seal(nil, nonce, plaintext, []byte(receiver.config.EncryptionContext))

	err = receiver.encoderInterface.WriteEncryptedChunk(ciphertext, nonce, receiver.outputWriter)
	return err
}

func (receiver *encryptor) generateNonce() ([]byte, error) {
	aead := receiver.authenticatedEncryption
	nonce := make([]byte, aead.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func (receiver *encryptor) setUpEncryptionOnce() error {
	if receiver.authenticatedEncryption != nil && len(receiver.encryptedDataKey) != 0 {
		return nil
	}

	encryptionKey, err := getEncryptionKey(receiver.config.KmsKey, receiver.config.EncryptionContext, receiver.config.AwsCredentialConfigPath, receiver.config.AwsProfile)
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

	return nil
}

func (receiver *encryptor) writeEncryptedDataKeyOnce() error {
	if receiver.encryptedDataKeyWritten {
		return nil
	}

	err := receiver.encoderInterface.WriteEncryptedDataKey(receiver.encryptedDataKey, receiver.outputWriter)
	if err != nil {
		return err
	}

	receiver.encryptedDataKeyWritten = true
	return nil
}
