package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/tskippervold/golang-base-server/internal/utils/slice"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/app/model"
	"github.com/tskippervold/golang-base-server/internal/utils/handler"
	"github.com/tskippervold/golang-base-server/internal/utils/log"
	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"
)

func authHandlers(r *mux.Router, env *env.Env) {
	r.Handle("/signup", signup(env)).Methods("POST")
	r.Handle("/login", login(env)).Methods("POST")

	r.Handle("/oauth/{provider}", oauth(env)).Queries("account_type", "{accountType}").Methods("GET")
	r.Handle("/oauth/{provider}/callback", oauthCallback(env)).Methods("GET")
}

func signup(env *env.Env) handler.Handler {

	type Request struct {
		Email    string            `json:"email"`
		Fullname string            `json:"fullname"`
		Password string            `json:"password"`
		Type     model.AccountType `json:"type"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		var body Request
		if err := request.Decode(r.Body, &body); err != nil {
			return respond.GenericServerError(err)
		}

		accountExists, err := model.AccountExists(env.DB, body.Email)
		if err != nil || accountExists {
			if accountExists {
				err = errors.New("Account exists")
				return respond.Error(err, http.StatusConflict, "Account with email already exists", "account_exists")
			}

			return respond.GenericServerError(err)
		}

		tx := env.DB.MustBegin()

		account := model.NewAccount(body.Email, body.Type)
		if err := account.Validate(); err != nil {
			return respond.Error(err, http.StatusBadRequest, "Invalid request", "invalid_request")
		}

		accountIID, err := account.Insert(tx)
		if err != nil {
			tx.Rollback()
			return respond.GenericServerError(err)
		}

		hash, _ := hashPassword(body.Password)
		ident := model.NewIdentity(model.IdentityTypeEmail, body.Email, accountIID)
		ident.PWHash = hash

		if err := ident.Insert(tx); err != nil {
			tx.Rollback()
			return respond.GenericServerError(err)
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return respond.GenericServerError(err)
		}

		return respond.Success(http.StatusCreated, nil)
	})
}

func login(env *env.Env) handler.Handler {

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		logger := log.ForRequest(r)

		var body Request
		if err := request.Decode(r.Body, &body); err != nil {
			return respond.GenericServerError(err)
		}

		account, err := model.GetAccount(env.DB, body.Email)
		if err != nil {
			logger.Error(err)
			return respond.Error(err, http.StatusNotFound, "Account not found", "no_account")
		}

		ident, err := model.GetIdentity(env.DB, model.IdentityTypeEmail, account.Email)
		if err != nil {
			logger.Error(err)
			return respond.Error(err, http.StatusNotFound, "Account not found", "no_account")
		}

		if err := compareHashAndPassword(ident.PWHash, body.Password); err != nil {
			logger.Error(err)
			return respond.Error(err, http.StatusUnauthorized, "Wrong credentials", "invalid_credentials")
		}

		return loginResponse(strconv.Itoa(account.IID))
	})
}

func oauth(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		accountType := mux.Vars(r)["accountType"]
		passthrough := struct {
			AccountType string `json:"account_type"`
		}{
			AccountType: accountType,
		}

		provider := OAuthProvider(mux.Vars(r)["provider"])
		if err := OAuthAuthenticate(provider, w, r, passthrough); err != nil {
			return respond.Error(err, http.StatusNotFound, "Not found", "not_found")
		}

		return nil
	})
}

func oauthCallback(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		logger := log.ForRequest(r)

		provider := OAuthProvider(mux.Vars(r)["provider"])
		oauth, err := OAuthCallback(provider, r)
		if err != nil {
			logger.Error(err)
			return respond.GenericServerError(err)
		}

		// Get the passthrough value, which in this case is `account_type`.
		var accountType string
		if passthroughJSON := oauth.JSONPassthrough; passthroughJSON != "" {
			var v map[string]string
			if err := json.Unmarshal([]byte(passthroughJSON), &v); err != nil {
				return respond.GenericServerError(err)
			}

			accountType = v["account_type"]
		}

		account, err := model.GetAccount(env.DB, oauth.Email)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Error(err)
				return respond.GenericServerError(err)
			}

			// Create account with Google oauth identity

			account = model.NewAccount(oauth.Email, model.AccountType(accountType))
			if err := account.Validate(); err != nil {
				return respond.Error(err, http.StatusBadRequest, "Invalid request", "invalid_request")
			}

			tx := env.DB.MustBegin()

			accountIID, err := account.Insert(tx)
			if err != nil {
				tx.Rollback()
				return respond.GenericServerError(err)
			}

			ident := model.NewIdentity(model.IdentityType(provider), oauth.Email, accountIID)
			ident.UID = oauth.ID

			if err := ident.Insert(tx); err != nil {
				tx.Rollback()
				return respond.GenericServerError(err)
			}

			if err := tx.Commit(); err != nil {
				tx.Rollback()
				return respond.GenericServerError(err)
			}
		} else {
			// Update account with Google oauth identity and account type

			_, err := model.GetIdentity(env.DB, model.IdentityType(provider), oauth.ID)
			if err != nil {
				logger.Error(err)
				return respond.Error(err, http.StatusNotFound, "Account not found", "no_account")
			}

			if slice.ContainsString(account.Type, accountType) == false {
				account.Type = append(account.Type, accountType)
			}

			if err := account.Validate(); err != nil {
				return respond.Error(err, http.StatusBadRequest, "Invalid request", "invalid_request")
			}

			if err := account.Update(env.DB); err != nil {
				logger.Error(err)
				return respond.GenericServerError(err)
			}
		}

		return loginResponse(strconv.Itoa(account.IID))
	})
}
