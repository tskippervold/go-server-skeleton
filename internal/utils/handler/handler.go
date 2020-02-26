package handler

import (
	log "github.com/sirupsen/logrus"
	"github.com/tskippervold/go-server-skeleton/internal/utils/respond"
	"net/http"
)

// Handler is equivalent to http.Handler but returns an error when the request
// should no longer be handled.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HandleFunc func(http.ResponseWriter, *http.Request) *respond.Response

func (f HandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if res := f(w, r); res != nil {
		res.Write(w)
		return
	}

	log.Info("No response returned to handler.")
}
