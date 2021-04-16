package crypt

import (
	"bytes"
	"encoding/gob"
)

func Decrypt(ciphertext []byte) ([]byte, error) {
	readerForCipherText := bytes.NewReader(ciphertext)
	var totalPlaintext []byte

	for readerForCipherText.Len() > 0 {
		var messageEnvelope envelope

		gobDecoder := gob.NewDecoder(readerForCipherText)
		err := gobDecoder.Decode(&messageEnvelope)
		if err != nil {
			return nil, err
		}

		decryptionKey, err := getDecryptionKey(messageEnvelope.Key)
		if err != nil {
			return nil, err
		}

		decryptor, err := createAuthenticatedEncryption(decryptionKey)
		if err != nil {
			return nil, err
		}

		clearPlaintextDataKey(decryptionKey)

		plaintext, err := decryptor.Open(nil, messageEnvelope.Nonce, messageEnvelope.Message, []byte("a test context string"))
		if err != nil {
			return nil, err
		}

		totalPlaintext = append(totalPlaintext, plaintext...)
	}

	return totalPlaintext, nil
}
