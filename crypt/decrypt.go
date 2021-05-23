package crypt

import (
	"crypto/cipher"
	"encoding/gob"
	"fmt"
	"github.com/halprin/cloud-backup-go/config"
	"io"
)

type decryptor struct {
	inputReader  io.Reader
	outputWriter io.Writer
	config config.BackupConfiguration
	gobDecoder *gob.Decoder

	authenticatedEncryption cipher.AEAD
	encryptedDataKey []byte
}

func NewDecryptor(inputReader io.Reader, outputWriter io.Writer, config config.BackupConfiguration) *decryptor {
	return &decryptor{
		inputReader: inputReader,
		outputWriter: outputWriter,
		config: config,
		gobDecoder: gob.NewDecoder(inputReader),
	}
}

func (receiver *decryptor) Decrypt() error {

	err := receiver.readPreamble()
	if err != nil {
		return err
	}

	err = receiver.readV1Preamble()
	if err != nil {
		return err
	}

	err = receiver.setUpEncryption()
	if err != nil {
		return err
	}

	err = receiver.decrypt()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *decryptor) readPreamble() error {
	var preambleStruct preamble

	err := receiver.gobDecoder.Decode(&preambleStruct)
	if err != nil {
		return err
	}

	if preambleStruct.Version != PreambleVersion {
		return fmt.Errorf("unsupported cipher format. cipher format was %s", preambleStruct.Version)
	}

	return nil
}

func (receiver *decryptor) readV1Preamble() error {
	var v1PreambleStruct v100Preamble

	err := receiver.gobDecoder.Decode(&v1PreambleStruct)
	if err != nil {
		return err
	}

	if len(v1PreambleStruct.EncryptedDataKey) == 0 {
		return fmt.Errorf("encrypted data key is empty")
	}

	receiver.encryptedDataKey = v1PreambleStruct.EncryptedDataKey
	return nil
}

func (receiver *decryptor) setUpEncryption() error {
	decryptionKey, err := getDecryptionKey(receiver.encryptedDataKey, receiver.config.EncryptionContext, receiver.config.AwsCredentialConfigPath, receiver.config.AwsProfile)
	if err != nil {
		return err
	}

	authenticatedEncryption, err := createAuthenticatedEncryption(decryptionKey)
	if err != nil {
		return err
	}

	clearPlaintextDataKey(decryptionKey)

	receiver.authenticatedEncryption = authenticatedEncryption
	return nil
}

func (receiver *decryptor) decrypt() error {
	for {
		var messageEnvelope v100Envelope

		err := receiver.gobDecoder.Decode(&messageEnvelope)
		if err == io.EOF {
			//we've exhausted the reader; end decrypting successfully
			return nil
		} else if err != nil {
			return err
		}

		plaintext, err := receiver.authenticatedEncryption.Open(nil, messageEnvelope.Nonce, messageEnvelope.CipherText, []byte(receiver.config.EncryptionContext))
		if err != nil {
			return err
		}

		_, err = receiver.outputWriter.Write(plaintext)
		if err != nil {
			return err
		}
	}
}
