package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"github.com/tskippervold/golang-base-server/internal/utils/response"
)

func AccountsHandlers(r *mux.Router, env *env.Env) {
	r.Handle("/", getAccounts(env)).Methods("GET")
	r.Handle("/", createAccount(env)).Methods("POST")
}

func getAccounts(env *env.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			Get the logger for this request.
			Using this logger makes it possible to trace this specific request.
			Notice the `traceId` field in the log for this request, tying it all together.
		*/
		logger := env.Log.ForRequest(r)
		logger.Info("Logging something..")
		response.Ok(w, "🙌")
	})
}

func createAccount(env *env.Env) http.Handler {

	/*
		Keep your request structs local to the handler.
		This way you dont need excessive files just containing requests.
	*/
	type Request struct {
		Foo string
		Bar int16
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body Request
		if err := request.Decode(r, &body); err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		/*
			You can either return a `struct` directly to JSON using:
			`response.OK(w, yourStructHere)`
			Or you can compose your response as shown:
		*/
		response.Ok(w, map[string]interface{}{
			"requestBody": body,
		})
	})
}
