package crypt

import (
	"crypto/aes"
	"crypto/cipher"
)

type envelope struct {
	Key     []byte
	Nonce   []byte
	Message []byte
}

type dataKey struct {
	encryptedDataKey []byte
	plaintextDataKey []byte
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
