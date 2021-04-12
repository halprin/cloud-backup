package crypt

import (
	"bytes"
	"encoding/gob"
	"log"
)

func Decrypt(ciphertext []byte) ([]byte, error) {
	var messageEnvelope envelope

	err := gob.NewDecoder(bytes.NewReader(ciphertext)).Decode(&messageEnvelope)
	if err != nil {
		log.Println("this is here")
		return nil, err
	}

	dataKey, err := getDecryptionKey(messageEnvelope.Key)
	if err != nil {
		return nil, err
	}

	decryptor, err := createAuthenticatedEncryption(dataKey)
	if err != nil {
		return nil, err
	}

	clearPlaintextDataKey(dataKey)

	plaintext, err := decryptor.Open(nil, messageEnvelope.Nonce, messageEnvelope.Message, []byte("a test context string"))
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
