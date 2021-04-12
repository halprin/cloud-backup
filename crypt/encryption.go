package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
)

func Encrypt(plaintext []byte) ([]byte, error) {

	dataKey, err := getEncryptionKey()
	if err != nil {
		return nil, err
	}

	encryptor, err := createAuthenticatedEncryption(dataKey)
	if err != nil {
		return nil, err
	}

	clearPlaintextDataKey(dataKey)

	nonce, err := generateNonce(encryptor)
	if err != nil {
		return nil, err
	}

	ciphertext := encryptor.Seal(nil, nonce, plaintext, []byte("a test context string"))

	messageEnvelope := &envelope{
		Key:     dataKey.encryptedDataKey,
		Nonce:   nonce,
		Message: ciphertext,
	}

	envelopeCipherText := &bytes.Buffer{}

	err = gob.NewEncoder(envelopeCipherText).Encode(messageEnvelope)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func generateNonce(aead cipher.AEAD) ([]byte, error) {
	nonce := make([]byte, aead.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
