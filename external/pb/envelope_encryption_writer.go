package pb

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

const PreambleVersionV100 = "1.0.0"

type ProtoBufEnvelopeEncryptionWriter struct {
}

func (p *ProtoBufEnvelopeEncryptionWriter) WriteEncryptedDataKey(encryptedDataKey []byte, writer io.Writer) error {
	//write the document preamble first
	err := writePreamble(writer)
	if err != nil {
		return err
	}

	//then write the v100 preamble because that's how we write out the encrypted data key
	err = writeV100Preamble(encryptedDataKey, writer)
	return err
}

func (p *ProtoBufEnvelopeEncryptionWriter) WriteEncryptedChunk(cipherText []byte, nonce []byte, writer io.Writer) error {
	v100Envelope := &V100Envelope{
		Nonce: nonce,
		CipherText: cipherText,
	}

	v100EnvelopeBytes, err := proto.Marshal(v100Envelope)
	if err != nil {
		return err
	}

	err = writeMessage(v100EnvelopeBytes, writer)
	return err
}

func writePreamble(writer io.Writer) error {
	preamble := &Preamble{
		Version: PreambleVersionV100,
	}

	preambleBytes, err := proto.Marshal(preamble)
	if err != nil {
		return err
	}

	err = writeMessage(preambleBytes, writer)
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

	err = writeMessage(v100PreambleBytes, writer)
	return err
}

//https://github.com/golang/protobuf/issues/507#issuecomment-391144637
func writeMessage(rawMessage []byte, writer io.Writer) error {
	_, err := writer.Write([]byte{1<<3 | 2})  //write a tag of field 1 of type bytes
	if err != nil {
		return err
	}

	//encode the length of the message in as a UVarInt
	var messageLengthArray [binary.MaxVarintLen64]byte
	messageLengthSlice := messageLengthArray[:binary.PutUvarint(messageLengthArray[:], uint64(len(rawMessage)))]

	_, err = writer.Write(messageLengthSlice)  //write the length
	if err != nil {
		return err
	}

	_, err = writer.Write(rawMessage)  //finally, write the actual message
	return err
}
