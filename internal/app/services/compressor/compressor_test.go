package compressor

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

type mockResponseWriter struct {
	header http.Header
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) Write(buffer []byte) (int, error) {
	return len(buffer), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
}

type mockReadCloser struct {
	io.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}

func TestNewGzipWriter(t *testing.T) {
	writer := NewGzipWriter(&mockResponseWriter{header: make(http.Header)})
	if writer == nil {
		t.Error("Expected GzipWriter to be created, but got nil")
	}
}

func TestGzipWriter_Write(t *testing.T) {
	responseWriter := &mockResponseWriter{header: make(http.Header)}
	gzipWriter := NewGzipWriter(responseWriter)
	buffer := []byte("test")
	n, err := gzipWriter.Write(buffer)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if n != len(buffer) {
		t.Errorf("Expected %d bytes written, but got %d", len(buffer), n)
	}
	if responseWriter.Header().Get("Content-Length") != strconv.Itoa(len(buffer)) {
		t.Errorf("Expected Content-Length to be %d, but got %s", len(buffer), responseWriter.Header().Get("Content-Length"))
	}

}

func TestGzipWriter_WriteHeader(t *testing.T) {
	responseWriter := &mockResponseWriter{header: make(http.Header)}
	gzipWriter := NewGzipWriter(responseWriter)
	statusCode := 200
	gzipWriter.WriteHeader(statusCode)
	if responseWriter.Header().Get("Content-Encoding") != "gzip" {
		t.Error("Expected Content-Encoding to be gzip, but got nil")
	}
}
