package gzip

import (
	"bufio"
	"compress/gzip"
	"errors"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type Writer struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func NewWriter(writer gin.ResponseWriter) *Writer {
	gz := gzip.NewWriter(writer)
	return &Writer{
		ResponseWriter: writer,
		writer:         gz,
	}
}

func (g *Writer) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}
func (g *Writer) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

func (g *Writer) WriteHeader(statusCode int) {
	g.ResponseWriter.WriteHeader(statusCode)
}

func (g *Writer) Flush() {
	_ = g.writer.Flush()
	g.ResponseWriter.Flush()
}

func (g *Writer) Close() error {
	return g.writer.Close()
}

var _ http.Hijacker = (*Writer)(nil)

func (g *Writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := g.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
