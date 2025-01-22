package middleware

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type (
	responseData struct {
		statusCode int
		bodySize   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.responseData.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseData.bodySize += size
	return size, err
}

func Logger(next http.Handler) http.Handler {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			statusCode: http.StatusOK,
			bodySize:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lrw, r)

		duration := time.Since(start)

		log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("duration", duration.String()).
			Msg("[REQUEST]")
		log.Info().
			Str("status", strconv.Itoa(responseData.statusCode)).
			Str("bodySize", strconv.Itoa(lrw.responseData.bodySize)).
			Msg("[RESPONSE]")
	}
	return http.HandlerFunc(logFunc)
}
