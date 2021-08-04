package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type EnvelopeEncryptionWriter interface {
	WriteEncryptedDataKey(encryptedDataKey []byte, writer io.Writer) error
	WriteEncryptedChunk(cipherText []byte, nonce []byte, writer io.Writer) error
}

type preamble struct {
	Version string
}

const PreambleVersion = "1.0.0"

type v100Preamble struct {
	EncryptedDataKey []byte
}

type v100Envelope struct {
	Nonce      []byte
	CipherText []byte
}

type dataKey struct {
	EncryptedDataKey []byte
	PlaintextDataKey []byte
}

func createAuthenticatedEncryption(dataKey *dataKey) (cipher.AEAD, error) {
	aesCipher, err := aes.NewCipher(dataKey.PlaintextDataKey)
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
	plaintextDataKey := dataKey.PlaintextDataKey

	for index := 0; index < len(dataKey.PlaintextDataKey); index++ {
		plaintextDataKey[index] = 0
	}

	dataKey.PlaintextDataKey = nil //sets the key to nil which will designate the data to be released via the garbage collector
}
