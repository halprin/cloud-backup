package crypt

import (
	"encoding/gob"
	"io"
)

type decryptor struct {
	inputReader  io.Reader
	outputWriter io.Writer
}

func NewDecryptor(inputReader io.Reader, outputWriter io.Writer) *decryptor {
	return &decryptor{
		inputReader: inputReader,
		outputWriter: outputWriter,
	}
}

func (receiver *decryptor) Decrypt() error {
	for {
		var messageEnvelope envelope

		gobDecoder := gob.NewDecoder(receiver.inputReader)
		err := gobDecoder.Decode(&messageEnvelope)
		if err == io.EOF {
			//we've exhausted the reader; end decrypting successfully
			return nil
		} else if err != nil {
			return err
		}

		decryptionKey, err := getDecryptionKey(messageEnvelope.Key)
		if err != nil {
			return err
		}

		authenticatedEncryption, err := createAuthenticatedEncryption(decryptionKey)
		if err != nil {
			return err
		}

		clearPlaintextDataKey(decryptionKey)

		plaintext, err := authenticatedEncryption.Open(nil, messageEnvelope.Nonce, messageEnvelope.Message, []byte("a test context string"))
		if err != nil {
			return err
		}

		_, err = receiver.outputWriter.Write(plaintext)
		if err != nil {
			return err
		}
	}
}
