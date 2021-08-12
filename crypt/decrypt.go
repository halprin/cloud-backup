package crypt

import (
	"crypto/cipher"
	"fmt"
	"github.com/halprin/cloud-backup/config"
	"github.com/halprin/cloud-backup/external/pb"
	"io"
)

type decryptor struct {
	inputReader      io.Reader
	outputWriter     io.Writer
	config           config.BackupConfiguration
	decoderInterface EnvelopeEncryptionReader

	authenticatedEncryption cipher.AEAD
	encryptedDataKey []byte
}

func NewDecryptor(inputReader io.Reader, outputWriter io.Writer, config config.BackupConfiguration) *decryptor {
	return &decryptor{
		inputReader: inputReader,
		outputWriter: outputWriter,
		config: config,
		decoderInterface: &pb.ProtoBufEnvelopeEncryptionReader{},
	}
}

func (receiver *decryptor) Decrypt() error {
	err := receiver.readEncryptedDataKey()
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

func (receiver *decryptor) readEncryptedDataKey() error {
	encryptedDataKey, err := receiver.decoderInterface.ReadEncryptedDataKey(receiver.inputReader)
	if err != nil {
		return err
	}

	if len(encryptedDataKey) == 0 {
		return fmt.Errorf("encrypted data key is empty")
	}

	receiver.encryptedDataKey = encryptedDataKey
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
		cipherText, nonce, err := receiver.decoderInterface.ReadEncryptedChunk(receiver.inputReader)
		if err == io.EOF {
			//we've exhausted the reader; end decrypting successfully
			return nil
		} else if err != nil {
			return err
		}

		plaintext, err := receiver.authenticatedEncryption.Open(nil, nonce, cipherText, []byte(receiver.config.EncryptionContext))
		if err != nil {
			return err
		}

		_, err = receiver.outputWriter.Write(plaintext)
		if err != nil {
			return err
		}
	}
}
