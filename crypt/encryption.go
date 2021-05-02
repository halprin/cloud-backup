package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"github.com/halprin/cloud-backup-go/config"
	"io"
)

type encryptor struct {
	outputWriter io.Writer
	config       config.BackupConfiguration
	gobEncoder   *gob.Encoder

	authenticatedEncryption cipher.AEAD
	encryptedDataKey        []byte

	preambleWritten   bool
	v1PreambleWritten bool
}

func NewEncryptor(outputWriter io.Writer, config config.BackupConfiguration) *encryptor {
	return &encryptor{
		outputWriter: outputWriter,
		config: config,
		gobEncoder: gob.NewEncoder(outputWriter),
	}
}

func (receiver *encryptor) Write(plaintext []byte) (int, error) {
	err := receiver.setUpEncryptionOnce()
	if err != nil {
		return 0, err
	}

	err = receiver.writePreambleOnce()
	if err != nil {
		return 0, err
	}

	err = receiver.writeV1PreambleOnce()
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

	messageEnvelope := &v100Envelope{
		Nonce:      nonce,
		CipherText: ciphertext,
	}

	err = receiver.gobEncoder.Encode(messageEnvelope)
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

	return nil
}

func (receiver *encryptor) writePreambleOnce() error {
	if receiver.preambleWritten {
		return nil
	}

	preambleStruct := &preamble{
		Version: PreambleVersion,
	}

	err := receiver.gobEncoder.Encode(preambleStruct)
	if err != nil {
		return err
	}

	receiver.preambleWritten = true
	return nil
}

func (receiver *encryptor) writeV1PreambleOnce() error {
	if receiver.v1PreambleWritten {
		return nil
	}

	v1preambleStruct := &v100Preamble{
		EncryptedDataKey: receiver.encryptedDataKey,
	}

	err := receiver.gobEncoder.Encode(v1preambleStruct)
	if err != nil {
		return err
	}

	receiver.v1PreambleWritten = true
	return nil
}
