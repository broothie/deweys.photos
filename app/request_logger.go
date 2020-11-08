package app

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type RequestLoggerRecorder struct {
	http.ResponseWriter
	Status   int
	BodySize int
}

func (r *RequestLoggerRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *RequestLoggerRecorder) Write(bytes []byte) (int, error) {
	r.BodySize = len(bytes)
	return r.ResponseWriter.Write(bytes)
}

func RequestLogger(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := RequestLoggerRecorder{ResponseWriter: w, Status: http.StatusOK}

		query := ""
		if r.URL.RawQuery != "" {
			query = fmt.Sprintf("?%s", r.URL.RawQuery)
		}

		requestBodySize := ""
		if r.ContentLength != 0 {
			requestBodySize = fmt.Sprintf("%dB ", r.ContentLength)
		}

		before := time.Now()
		next.ServeHTTP(&response, r)
		elapsed := time.Since(before)

		responseBodySize := ""
		if response.BodySize > 0 {
			responseBodySize = fmt.Sprintf("%dB ", response.BodySize)
		}

		logger.Printf("%s %s%s %s| %d %s %s| %v",
			// Request
			r.Method,
			r.URL.Path,
			query,
			requestBodySize,
			// Response
			response.Status,
			http.StatusText(response.Status),
			responseBodySize,
			// Elapsed
			elapsed,
		)
	})
}
