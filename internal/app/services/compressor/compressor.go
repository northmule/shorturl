package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
)

// GzipWriter структура записи
type GzipWriter struct {
	Response http.ResponseWriter
	Writer   *gzip.Writer
}

// GzipReader структура чтения
type GzipReader struct {
	IoReader io.ReadCloser
	Reader   *gzip.Reader
}

// NewGzipWriter Конструктор записи.
func NewGzipWriter(response http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		Response: response,
		Writer:   gzip.NewWriter(response),
	}
}

// NewGzipReader Конструктор чтение.
func NewGzipReader(reader io.ReadCloser) (*GzipReader, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return &GzipReader{
		IoReader: reader,
		Reader:   gzipReader,
	}, nil
}

// Header загологвки.
func (g *GzipWriter) Header() http.Header {
	return g.Response.Header()
}

// Write Упакованные данные отправляемые клиенту.
func (g *GzipWriter) Write(buffer []byte) (int, error) {
	g.Response.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
	return g.Writer.Write(buffer)
}

// WriteHeader Переопределение записи заголовков
func (g *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		g.Response.Header().Set("Content-Encoding", "gzip")
	}
	g.Response.WriteHeader(statusCode)
}

// Close Переопределение закрытия
func (g *GzipWriter) Close() error {
	return g.Writer.Close()
}

// Read Переопределение чтения
func (g GzipReader) Read(buffer []byte) (n int, err error) {
	return g.Reader.Read(buffer)
}

// Close Переопределение закрытия
func (g GzipReader) Close() error {
	if err := g.IoReader.Close(); err != nil {
		return err
	}
	return g.Reader.Close()
}
