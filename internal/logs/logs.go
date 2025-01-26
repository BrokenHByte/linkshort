package logs

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

var sugar *zap.SugaredLogger = nil

func Logs() *zap.SugaredLogger {
	if sugar == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			// вызываем панику, если ошибка
			panic("cannot initialize zap")
		}
		defer logger.Sync()
		sugar = logger.Sugar()
	}
	return sugar
}

func LoggingRequest(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Logs().Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration.Microseconds(),
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
