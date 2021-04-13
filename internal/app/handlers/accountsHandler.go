package handlers

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/tskippervold/golang-base-server/internal/app/auth"

	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/utils/handler"
)

func AccountsHandlers(r *mux.Router, env *env.Env) {
	r.Use(auth.JWTMiddleware)

	r.Handle("/me", getMyAccount(env)).Methods("GET")
	r.Handle("/me", updateMyAccount(env)).Methods("PUT")

	r.Handle("/tokens", getAccessTokens(env)).Methods("GET")
}

func getMyAccount(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		account, err := auth.AuthenticatedAccount(r, env.DB)
		if err != nil {
			return respond.GenericServerError(err)
		}

		return respond.Success(http.StatusOK, account.Profile)
	})
}

func updateMyAccount(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		/*account, err := auth.AuthenticatedAccount(r, env.DB)
		if err != nil {
			return respond.GenericServerError(err)
		}*/

		var update map[string]interface{}
		request.Decode(r.Body, &update)

		return respond.Success(http.StatusOK, update)
	})
}

func getAccessTokens(env *env.Env) handler.Handler {
	type Response struct {
		DevToken  string  `json:"dev_token"`
		ProdToken *string `json:"prod_token"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {

		devToken := uuid.NewV4().String()

		return respond.Success(http.StatusOK, Response{
			DevToken:  devToken,
			ProdToken: nil,
		})
	})
}
