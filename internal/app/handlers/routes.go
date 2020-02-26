package handlers

import (
	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
)

// Setup configures all the routes for this mux.Router `r`.
// You should pass the Env `env` to all the registered handlers.
func Setup(r *mux.Router, env *env.Env) {
	AccountsHandlers(r.PathPrefix("/accounts").Subrouter(), env)
}
