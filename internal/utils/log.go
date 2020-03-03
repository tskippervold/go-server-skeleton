package utils

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
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

type Log struct {
	logger *logrus.Logger
	fields *logrus.Fields
}

// Creates a new Log instance.
func NewLogger() *Log {
	_, file, no, _ := runtime.Caller(1)

	logger := Log{
		logger: logrus.New(),
		fields: &logrus.Fields{
			string(logFilename): fmt.Sprintf("%s#%d", file, no),
		},
	}

	logger.logger.Formatter = &logrus.TextFormatter{
		DisableTimestamp: false,
	}

	return &logger
}

// Creates a Log instance containing a traceId from `http.Request`.
func (l *Log) ForRequest(r *http.Request) *Log {
	c := NewLogger()

	if f, ok := r.Context().Value(requestLoggerFields).(logrus.Fields); ok {
		c.fields = &f
	}

	return c
}

func (l *Log) Debug(a ...interface{}) {
	l.log(logrus.DebugLevel, a...)
}

func (l *Log) Info(a ...interface{}) {
	l.log(logrus.InfoLevel, a...)
}

func (l *Log) Error(a ...interface{}) {
	l.log(logrus.ErrorLevel, a...)
}

// Logs every HTTP request and sets a traceId in the request's context.
func (l *Log) HTTPRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid := uuid.NewV4().String()
		fields := logrus.Fields{
			string(logTraceID):          tid,
			string(logRequestMethod):    r.Method,
			string(logRequestPath):      r.URL.Path,
			string(logRequestUserAgent): r.UserAgent(),
		}

		ctx := context.WithValue(r.Context(), requestLoggerFields, fields)
		r = r.WithContext(ctx)

		// Logging the request before handling it
		logger := l.logger.WithFields(fields)
		logger.Info()

		next.ServeHTTP(w, r)
	})
}

func (l *Log) log(level logrus.Level, a ...interface{}) {
	msg := fmtLogMsg(a)

	if l.fields != nil {
		l.logger.WithFields(*l.fields).Log(level, msg)
	} else {
		l.logger.Log(level, msg)
	}
}

func fmtLogMsg(args ...interface{}) string {
	var argsStr = []string{}

	for _, v := range args {
		// error needs special treatment
		if _, ok := v.(error); ok {
			argsStr = append(argsStr, fmt.Sprint(v))
			continue
		}

		argsStr = append(argsStr, fmt.Sprint(v))
	}

	return strings.Trim(strings.Join(argsStr, ", "), " []")
}
