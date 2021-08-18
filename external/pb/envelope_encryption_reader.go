package pb

import (
	"encoding/binary"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

type ProtoBufEnvelopeEncryptionReader struct {
	version string
}

func (p *ProtoBufEnvelopeEncryptionReader) ReadEncryptedDataKey(reader io.Reader) ([]byte, error) {
	err := p.readPreamble(reader)
	if err != nil {
		return nil, err
	}

	encryptedDataKey, err := p.readVersionPreamble(reader)
	if err != nil {
		return nil, err
	}

	return encryptedDataKey, nil
}

func (p *ProtoBufEnvelopeEncryptionReader) ReadEncryptedChunk(reader io.Reader) ([]byte, []byte, error) {
	rawMessage, err := readMessage(reader)
	if err != nil {
		return nil, nil, err
	}

	v100Envelope := &V100Envelope{}
	err = proto.Unmarshal(rawMessage, v100Envelope)
	if err != nil {
		return nil, nil, err
	}

	return v100Envelope.GetCipherText(), v100Envelope.GetNonce(), nil
}

func (p *ProtoBufEnvelopeEncryptionReader) readPreamble(reader io.Reader) error {
	rawMessage, err := readMessage(reader)
	if err != nil {
		return err
	}

	preamble := &Preamble{}
	err = proto.Unmarshal(rawMessage, preamble)
	if err != nil {
		return err
	}

	p.version = preamble.GetVersion()

	return nil
}

func (p *ProtoBufEnvelopeEncryptionReader) readVersionPreamble(reader io.Reader) ([]byte, error) {
	if p.version == PreambleVersionV100 {
		return readV100Preamble(reader)
	}

	return nil, fmt.Errorf("%s is not a supported version", p.version)
}

func readV100Preamble(reader io.Reader) ([]byte, error) {
	rawMessage, err := readMessage(reader)
	if err != nil {
		return nil, err
	}

	v100Preamble := &V100Preamble{}
	err = proto.Unmarshal(rawMessage, v100Preamble)
	if err != nil {
		return nil, err
	}

	return v100Preamble.GetEncryptedDataKey(), nil
}

//https://github.com/golang/protobuf/issues/507#issuecomment-391144637
func readMessage(reader io.Reader) ([]byte, error) {
	field := make([]byte, 1)
	_, err := reader.Read(field)  //read a tag of field 1 and throw it away since we don't use this actually
	if err != nil {
		return nil, err
	}

	//read the length of the upcoming message, encoded in a UVarInt
	byteReader := newUnbufferedByteReader(reader)
	messageLength, err := binary.ReadUvarint(byteReader)
	if err != nil {
		return nil, err
	}

	rawMessage := make([]byte, messageLength)
	_, err = io.ReadFull(reader, rawMessage)
	if err != nil {
		return nil, err
	}

	return rawMessage, nil
}

type unbufferedByteReader struct {
	reader io.Reader
}

func newUnbufferedByteReader(reader io.Reader) *unbufferedByteReader {
	return &unbufferedByteReader{
		reader: reader,
	}
}

func (u *unbufferedByteReader) ReadByte() (byte, error) {
	theByte := make([]byte, 1)
	lengthRead, err := u.reader.Read(theByte)
	if err != nil {
		return 0, err
	} else if lengthRead == 0 {
		return 0, errors.New("unable to read a single byte")
	}

	return theByte[0], nil
}

