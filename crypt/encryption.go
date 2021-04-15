package crypt

import (
	"bytes"
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
	cipherText, err := Encrypt(plaintext)
	if err != nil {
		return 0, err
	}

	writtenSize, err := receiver.outputWriter.Write(cipherText)
	return writtenSize, err
}

func Encrypt(plaintext []byte) ([]byte, error) {

	encryptionKey, err := getEncryptionKey()
	if err != nil {
		return nil, err
	}

	authenticatedEncryption, err := createAuthenticatedEncryption(encryptionKey)
	if err != nil {
		return nil, err
	}

	clearPlaintextDataKey(encryptionKey)

	nonce, err := generateNonce(authenticatedEncryption)
	if err != nil {
		return nil, err
	}

	ciphertext := authenticatedEncryption.Seal(nil, nonce, plaintext, []byte("a test context string"))

	messageEnvelope := &envelope{
		Key:     encryptionKey.encryptedDataKey,
		Nonce:   nonce,
		Message: ciphertext,
	}

	envelopeCipherText := &bytes.Buffer{}

	err = gob.NewEncoder(envelopeCipherText).Encode(messageEnvelope)
	if err != nil {
		return nil, err
	}

	return envelopeCipherText.Bytes(), nil
}

func generateNonce(aead cipher.AEAD) ([]byte, error) {
	nonce := make([]byte, aead.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
