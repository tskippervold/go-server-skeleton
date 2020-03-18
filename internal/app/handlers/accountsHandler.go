package handlers

import (
	"database/sql"
	"net/http"

	"github.com/tskippervold/golang-base-server/internal/app/model"

	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/utils/handler"
	"github.com/tskippervold/golang-base-server/internal/utils/log"
)

func AccountsHandlers(r *mux.Router, env *env.Env) {
	r.Handle("/me", getMyAccount(env)).Methods("GET")
	r.Handle("/update", updateAccount(env)).Methods("PUT")
}

func getMyAccount(env *env.Env) handler.Handler {
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

func updateAccount(env *env.Env) handler.Handler {

	type Update struct {
		Summary         string   `json:"summary"`
		AreaOfExpertise []string `json:"area_of_expertise"`
		Certifications  []string `json:"certifications"`
	}

	type Request struct {
		AccountEmail string `json:"accountEmail"`
		Update       Update `json:"update"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		var body Request
		if err := request.Decode(r.Body, &body); err != nil {
			return respond.GenericServerError(err)
		}

		account, err := model.GetAccount(env.DB, body.AccountEmail)
		if err != nil {
			return respond.Error(err, http.StatusNotFound, "Account not found", "no_account")
		}

		if summary := body.Update.Summary; summary != "" {
			account.Summary = sql.NullString{String: summary, Valid: true}
		}

		if aoe := body.Update.AreaOfExpertise; aoe != nil {
			account.AreaOfExpertise = aoe
		}

		if certs := body.Update.Certifications; certs != nil {
			account.Certifications = certs
		}

		if err := account.Update(env.DB); err != nil {
			return respond.GenericServerError(err)
		}

		return respond.Success(200, map[string]string{
			"status": "ok",
		})
	})
}
