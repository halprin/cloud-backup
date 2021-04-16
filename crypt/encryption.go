package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"io"
)

type encryptor struct {
	outputWriter io.Writer
}

func NewEncryptor(outputWriter io.Writer) *encryptor {
	return &encryptor{
		outputWriter: outputWriter,
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

	encryptionKey, err := getEncryptionKey()
	if err != nil {
		return err
	}

	authenticatedEncryption, err := createAuthenticatedEncryption(encryptionKey)
	if err != nil {
		return err
	}

	clearPlaintextDataKey(encryptionKey)

	nonce, err := receiver.generateNonce(authenticatedEncryption)
	if err != nil {
		return err
	}

	ciphertext := authenticatedEncryption.Seal(nil, nonce, plaintext, []byte("a test context string"))

	messageEnvelope := &envelope{
		Key:     encryptionKey.EncryptedDataKey,
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
