package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

type GzipWriter struct {
	Response http.ResponseWriter
	Gz       *gzip.Writer
}

type GzipReader struct {
	IoReader io.ReadCloser
	Gz       *gzip.Reader
}

func NewGzipWriter(response http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		Response: response,
		Gz:       gzip.NewWriter(response),
	}
}

func NewGzipReader(reader io.ReadCloser) (*GzipReader, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return &GzipReader{
		IoReader: reader,
		Gz:       gzipReader,
	}, nil
}

func (g *GzipWriter) Header() http.Header {
	return g.Response.Header()
}

// Упакованные данные отправляемые клиенту
func (g *GzipWriter) Write(buffer []byte) (int, error) {
	return g.Gz.Write(buffer)
}

func (g *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		g.Response.Header().Set("Content-Encoding", "gzip")
	}
	g.Response.WriteHeader(statusCode)
}

func (g *GzipWriter) Close() error {
	return g.Gz.Close()
}

func (g GzipReader) Read(buffer []byte) (n int, err error) {
	return g.Gz.Read(buffer)
}

func (g GzipReader) Close() error {
	if err := g.IoReader.Close(); err != nil {
		return err
	}
	return g.Gz.Close()
}
