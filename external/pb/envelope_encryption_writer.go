package pb

import (
	"google.golang.org/protobuf/proto"
	"io"
)

const PreambleVersion = "1.0.0"

type ProtoBufEnvelopeEncryptionWriter struct {
}

func (p ProtoBufEnvelopeEncryptionWriter) WriteEncryptedDataKey(encryptedDataKey []byte, writer io.Writer) error {
	//write the document preamble first
	err := writePreamble(writer)
	if err != nil {
		return err
	}

	//then write the v100 preamble because that's how we write out the encrypted data key
	err = writeV100Preamble(encryptedDataKey, writer)
	return err
}

func (p ProtoBufEnvelopeEncryptionWriter) WriteEncryptedChunk(cipherText []byte, nonce []byte, writer io.Writer) error {
	v100Envelope := &V100Envelope{
		Nonce: nonce,
		CipherText: cipherText,
	}

	v100EnvelopeBytes, err := proto.Marshal(v100Envelope)
	if err != nil {
		return err
	}

	_, err = writer.Write(v100EnvelopeBytes)
	return err
}

func writePreamble(writer io.Writer) error {
	preamble := &Preamble{
		Version: PreambleVersion,
	}

	preambleBytes, err := proto.Marshal(preamble)
	if err != nil {
		return err
	}

	_, err = writer.Write(preambleBytes)
	return err
}

func writeV100Preamble(encryptedDataKey []byte, writer io.Writer) error {
	v100Preamble := &V100Preamble{
		EncryptedDataKey: encryptedDataKey,
	}

	v100PreambleBytes, err := proto.Marshal(v100Preamble)
	if err != nil {
		return err
	}

	_, err = writer.Write(v100PreambleBytes)
	return err
}

