package auth

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

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

	r.Handle("/oauth/{provider}", oauth(env)).Queries("account_type", "{accountType}", "redirect_uri", "{redirectUri}").Methods("GET")
	r.Handle("/oauth/{provider}/callback", oauthCallback(env)).Methods("GET")
}

func signup(env *env.Env) handler.Handler {

	type Request struct {
		Email    string            `json:"email"`
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

		response := loginResponse(account.Email)
		return respond.Success(http.StatusCreated, response)
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
			return respond.Error(err, http.StatusUnprocessableEntity, "Wrong credentials", "invalid_credentials")
		}

		return loginResponse(account.Email)
	})
}

func oauth(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		v := mux.Vars(r)
		passthrough := struct {
			AccountType string `json:"account_type"`
			RedirectURI string `json:"redirect_uri"`
		}{
			AccountType: v["accountType"],
			RedirectURI: v["redirectUri"],
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
		var redirectURI string
		if passthroughJSON := oauth.JSONPassthrough; passthroughJSON != "" {
			var v map[string]string
			if err := json.Unmarshal([]byte(passthroughJSON), &v); err != nil {
				return respond.GenericServerError(err)
			}

			accountType = v["account_type"]
			redirectURI = v["redirect_uri"]
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

			ident := model.NewIdentity(model.IdentityType(provider), oauth.ID, accountIID)
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
			ident := model.NewIdentity(model.IdentityType(provider), oauth.ID, account.IID)
			if err := ident.InsertIfNew(env.DB); err != nil {
				return respond.GenericServerError(err)
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

		response := loginResponse(account.Email)
		responseBytes, _ := response.MarshalJSON()
		b64 := base64.URLEncoding.EncodeToString(responseBytes)

		u, _ := url.Parse(redirectURI)
		q, _ := url.ParseQuery(u.RawQuery)
		q.Set("oauth_r", b64)
		u.RawQuery = q.Encode()

		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
		return nil

	})
}
