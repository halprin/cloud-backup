package compression

import (
	"bytes"
	"compress/gzip"
	"io"
)

type compressor struct {
	gzipWriter *gzip.Writer
}

func NewCompressor(outputWriter io.Writer) *compressor {
	gzipWriter := gzip.NewWriter(outputWriter)
	return &compressor{
		gzipWriter: gzipWriter,
	}
}

func (receiver *compressor) Writer() io.Writer {
	return receiver.gzipWriter
}

func (receiver *compressor) Close() error {
	err := receiver.gzipWriter.Close()
	return err
}

func Compress(contentToCompress []byte) ([]byte, error) {
	byteBuffer := bytes.Buffer{}
	gzipWriter := gzip.NewWriter(&byteBuffer)
	defer gzipWriter.Close()

	_, err := gzipWriter.Write(contentToCompress)
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return byteBuffer.Bytes(), nil
}
