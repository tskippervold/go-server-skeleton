package handlers

import (
	"github.com/gorilla/mux"
	"github.com/tskippervold/go-server-skeleton/internal/app"
	"github.com/tskippervold/go-server-skeleton/internal/utils/handler"
	"github.com/tskippervold/go-server-skeleton/internal/utils/respond"
	"net/http"
)

// Setup configures all the routes for this mux.Router `r`.
// You should pass the Env `env` to all the registered handlers.
func Setup(r *mux.Router, env *app.Environment) {
	testHandler(r.PathPrefix("/test").Subrouter(), env)
}

func testHandler(r *mux.Router, env *app.Environment) {
	r.Handle("/", test(env)).Methods("GET")
}

func test(env *app.Environment) handler.Handler {
	return handler.HandleFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		return respond.Success(http.StatusOK, "Hello world")
	})
}
