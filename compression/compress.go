package compression

import (
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
