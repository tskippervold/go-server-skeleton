package handler

import (
	"net/http"

	"github.com/tskippervold/golang-base-server/internal/utils/respond"
)

// Handler is equivalent to http.Handler but returns an error when the request
// should no longer be handled.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) *respond.Response

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := f(w, r)
	res.Write(w)
}
