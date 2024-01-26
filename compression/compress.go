package compression

import (
	"github.com/klauspost/pgzip"
	"io"
)

type compressor struct {
	gzipWriter *pgzip.Writer
}

func NewCompressor(outputWriter io.Writer) (*compressor, error) {
	gzipWriter, err := pgzip.NewWriterLevel(outputWriter, pgzip.BestCompression)
	if err != nil {
		return nil, err
	}

	return &compressor{
		gzipWriter: gzipWriter,
	}, nil
}

func (receiver *compressor) Writer() io.Writer {
	return receiver.gzipWriter
}

func (receiver *compressor) Close() error {
	err := receiver.gzipWriter.Close()
	return err
}
