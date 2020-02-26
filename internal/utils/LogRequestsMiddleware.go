package utils

import (
	"context"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type contextKey string

const (
	requestLoggerFields = contextKey("_rRequestLoggerFields")
	logFilename         = contextKey("filename")
	logTraceID          = contextKey("traceId")
	logRequestMethod    = contextKey("method")
	logRequestPath      = contextKey("path")
	logRequestUserAgent = contextKey("userAgent")
)

// Logs every HTTP request and sets a traceId in the request's context.
func LogRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid := uuid.NewV4().String()
		fields := log.Fields{
			string(logTraceID):          tid,
			string(logRequestMethod):    r.Method,
			string(logRequestPath):      r.URL.Path,
			string(logRequestUserAgent): r.UserAgent(),
		}

		// Logging the request before handling it
		logger := log.WithFields(fields)
		logger.Info()

		ctx := context.WithValue(r.Context(), requestLoggerFields, fields)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
