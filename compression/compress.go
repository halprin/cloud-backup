package compression

import (
	"bytes"
	"compress/gzip"
)

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
