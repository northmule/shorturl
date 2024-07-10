package middlewarehandler

import (
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/compressor"
	"net/http"
	"strings"
)

func MiddlewareGzipCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		expectedContentTypes := map[string]bool{
			"application/json":   true,
			"text/html":          true,
			"application/x-gzip": true,
		}

		requestContentType := request.Header.Values("Content-type")
		contentTypeIsSupported := false
		for _, contentType := range requestContentType {
			if _, ok := expectedContentTypes[contentType]; ok {
				contentTypeIsSupported = true
				break
			}
		}
		if !contentTypeIsSupported {
			logger.Log.Sugar().Infof("Ответ сервера не сжимается, для Content-type не поддерживается сжатие")
			next.ServeHTTP(response, request)
			return
		}
		acceptIsEncodingSupported := false
		for _, acceptEncoding := range request.Header.Values("Accept-Encoding") {
			if strings.Contains(acceptEncoding, "gzip") {
				acceptIsEncodingSupported = true
				break
			}
		}
		acceptIsContentEncodingSupported := false
		for _, contentEncoding := range request.Header.Values("Content-Encoding") {
			if strings.Contains(contentEncoding, "gzip") {
				acceptIsContentEncodingSupported = true
				break
			}
		}
		modifiedResponse := response
		// Клиент поддерживает gzip, сжимаем данные для него
		if acceptIsEncodingSupported {
			gzipWriter := compressor.NewGzipWriter(response)
			modifiedResponse = gzipWriter
			defer gzipWriter.Close()
		}

		// Пришли сжатые данные, надо распаковать
		if acceptIsContentEncodingSupported {
			gzipReader, err := compressor.NewGzipReader(request.Body)

			if err != nil {
				logger.Log.Sugar().Error("Ошибка ", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			request.Body = gzipReader
			defer gzipReader.Close()
		}

		next.ServeHTTP(modifiedResponse, request)
	})
}
