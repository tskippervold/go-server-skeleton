package handlers

import (
	"net/http"

	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/utils/handler"
	"github.com/tskippervold/golang-base-server/internal/utils/log"
)

func AccountsHandlers(r *mux.Router, env *env.Env) {
	r.Handle("/", getAccounts(env)).Methods("GET")
	r.Handle("/", createAccount(env)).Methods("POST")
}

func getAccounts(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		/*
			Get the logger for this request.
			Using this logger makes it possible to trace this specific request.
			Notice the `traceId` field in the log for this request, tying it all together.
		*/
		logger := log.ForRequest(r)
		logger.Info("Logging something..")
		return respond.Success(http.StatusOK, "🙌")
	})
}

func createAccount(env *env.Env) handler.Handler {

	/*
		Keep your request structs local to the handler.
		This way you dont need excessive files just containing requests.
	*/
	type Request struct {
		Foo string
		Bar int16
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		var body Request
		if err := request.Decode(r, &body); err != nil {
			return respond.GenericServerError(err)
		}

		/*
			You can either return a `struct` directly to JSON using:
			`respond.OK(w, yourStructHere)`
			Or you can compose your response as shown:
		*/
		return respond.Success(http.StatusOK, map[string]interface{}{
			"requestBody": body,
		})
	})
}
