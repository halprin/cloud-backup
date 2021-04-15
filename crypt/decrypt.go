package crypt

import (
	"bytes"
	"encoding/gob"
	"log"
)

func Decrypt(ciphertext []byte) ([]byte, error) {
	readerForCipherText := bytes.NewReader(ciphertext)
	gobDecoder := gob.NewDecoder(readerForCipherText)


	//TODO: the idea is to do a for loop that loops until the readerForCipherText.Len() is 0.  In the loop, I constantly call gobDecoder.Decode(&messageEnvelope).  I think that might work.
	var messageEnvelope envelope

	err := gobDecoder.Decode(&messageEnvelope)
	if err != nil {
		log.Println("this is here")
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

	return plaintext, nil
}
