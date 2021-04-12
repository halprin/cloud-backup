package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
)

type envelope struct {
	Key     []byte
	Nonce   []byte
	Message []byte
}

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

func createAuthenticatedEncryption(dataKey *dataKey) (cipher.AEAD, error) {
	aesCipher, err := aes.NewCipher(dataKey.plaintextDataKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

func clearPlaintextDataKey(dataKey *dataKey)  {
	plaintextDataKey := dataKey.plaintextDataKey

	for index := 0; index < len(dataKey.plaintextDataKey); index++ {
		plaintextDataKey[index] = 0
	}

	dataKey.plaintextDataKey = nil  //sets the key to nil which will designate the data to be released via the garbage collector
}

func generateNonce(aead cipher.AEAD) ([]byte, error) {
	nonce := make([]byte, aead.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
