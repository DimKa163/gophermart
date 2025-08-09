package gzip

import (
	"compress/gzip"
	"io"
)

type Reader struct {
	reader io.ReadCloser
	gz     *gzip.Reader
}

func NewReader(r io.ReadCloser) (*Reader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &Reader{
		reader: r,
		gz:     gz,
	}, nil
}

func (g *Reader) Read(p []byte) (n int, err error) {
	v, err := g.gz.Read(p)
	return v, err
}

func (g *Reader) Close() error {
	if err := g.reader.Close(); err != nil {
		return err
	}
	return g.gz.Close()
}
