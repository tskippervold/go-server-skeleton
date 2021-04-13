package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/app/auth"
	"github.com/tskippervold/golang-base-server/internal/app/proff"
	"github.com/tskippervold/golang-base-server/internal/utils/handler"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"
)

func CompanyHandler(r *mux.Router, env *env.Env) {
	r.Use(auth.JWTMiddleware)

	r.Handle("/lookup/{orgnum}", lookupCompany(env)).Methods("GET")
	r.Handle("/lookup/{orgnum}/ownership", ownership(env)).Methods("GET")
}

func lookupCompany(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		orgnum := mux.Vars(r)["orgnum"]

		if r.Header.Get("X-Mocked") != "" {
			file, err := ioutil.ReadFile("./mocks/companyLookup.json")
			if err != nil {
				return respond.GenericServerError(err)
			}

			var result interface{}
			err = json.Unmarshal(file, &result)
			if err != nil {
				return respond.GenericServerError(err)
			}

			return respond.Success(http.StatusOK, result)
		}

		result, err := proff.LookupCompany(orgnum)
		if err != nil {
			return respond.GenericServerError(err)
		}

		return respond.Success(http.StatusOK, result)
	})
}

func ownership(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		orgnum := mux.Vars(r)["orgnum"]

		if r.Header.Get("X-Mocked") != "" {
			file, err := ioutil.ReadFile("./mocks/companyOwners.json")
			if err != nil {
				return respond.GenericServerError(err)
			}

			var result interface{}
			err = json.Unmarshal(file, &result)
			if err != nil {
				return respond.GenericServerError(err)
			}

			return respond.Success(http.StatusOK, result)
		}

		result, err := proff.GetCompanyOwners(orgnum)
		if err != nil {
			return respond.GenericServerError(err)
		}

		return respond.Success(http.StatusOK, result)
	})
}
