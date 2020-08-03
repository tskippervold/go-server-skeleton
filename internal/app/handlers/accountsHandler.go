package handlers

import (
	"net/http"

	"github.com/tskippervold/golang-base-server/internal/app/auth"

	"github.com/tskippervold/golang-base-server/internal/utils/respond"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/utils/handler"
)

func AccountsHandlers(r *mux.Router, env *env.Env) {
	r.Handle("/me", getMyAccount(env)).Methods("GET")
}

func getMyAccount(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		account, err := auth.AuthenticatedAccount(r, env.DB)
		if err != nil {
			return respond.GenericServerError(err)
		}

		return respond.Success(http.StatusOK, account)
	})
}
