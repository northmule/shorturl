package middlewarehandler

import (
	"github.com/northmule/shorturl/internal/app/logger"
	"net/http"
	"time"
)

type loggingData struct {
	url         string
	method      string
	executeTime string
	size        int
	statusCode  int
}

// ResponseWriterWrapper структура для захвата ответа
type ResponseWriterWrapper struct {
	originResponse *http.ResponseWriter
	originRequest  *http.Request
	loggingData    *loggingData
}

func NewResponseWriterWrapper(rw http.ResponseWriter, request http.Request) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		originResponse: &rw,
		originRequest:  &request,
		loggingData:    &loggingData{},
	}
}

// Write при записи ответа
func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	size, err := (*rww.originResponse).Write(buf)
	if err != nil {
		return 0, err
	}
	rww.loggingData.size += size
	return size, nil
}

// Header срабатывает перед сеттом заголовка
func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.originResponse).Header()
}

// WriteHeader срабоатет при записи заголовков в основном запросе
func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.originResponse).WriteHeader(statusCode)

	rww.loggingData.statusCode = statusCode
	rww.loggingData.method = rww.originRequest.Method
	rww.loggingData.url = rww.originRequest.URL.String()
}

// MiddlewareLogger логгер запросов/ответов
func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		listenerResponse := NewResponseWriterWrapper(response, *request)
		next.ServeHTTP(listenerResponse, request)

		listenerResponse.loggingData.executeTime = time.Since(startTime).String()
		sugarLogger := logger.Log.Sugar()

		sugarLogger.Infof(
			"Request: URL %s, method: %s, executeTime: %s",
			listenerResponse.loggingData.url,
			listenerResponse.loggingData.method,
			listenerResponse.loggingData.executeTime,
		)

		sugarLogger.Infof(
			"Response: statusCode: %d, size: %d bytes",
			listenerResponse.loggingData.statusCode,
			listenerResponse.loggingData.size,
		)
	})
}
