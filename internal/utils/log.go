package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	requestContextTraceId = "_rTraceId"
	kLogTraceId           = "traceId"
)

type Log struct {
	logger *logrus.Logger
	trace  *string
}

// Creates a new Log instance.
func NewLogger() *Log {
	logger := Log{
		logger: logrus.New(),
		trace:  nil,
	}

	logger.logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
	}

	return &logger
}

// Creates a Log instance containing a traceId from `http.Request`.
func (l *Log) ForRequest(r *http.Request) *Log {
	tid := r.Context().Value(requestContextTraceId).(string)

	c := NewLogger()
	c.trace = &tid

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
		ctx := context.WithValue(r.Context(), requestContextTraceId, tid)
		requestWithTrace := r.WithContext(ctx)

		l.logger.WithField(kLogTraceId, tid).Info(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

		next.ServeHTTP(w, requestWithTrace)
	})
}

func (l *Log) log(level logrus.Level, a ...interface{}) {
	msg := fmtLogMsg(a)

	if l.trace != nil {
		l.logger.WithField(kLogTraceId, l.trace).Log(level, msg)
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
