package auth

import (
	"errors"
	"net/http"

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
}

func signup(env *env.Env) handler.Handler {

	type Request struct {
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Password string `json:"password"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		var body Request
		if err := request.Decode(r, &body); err != nil {
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

		account := model.NewAccount(body.Email)
		uIID, err := account.Insert(tx)
		if err != nil {
			tx.Rollback()
			return respond.GenericServerError(err)
		}

		hash, _ := hashPassword(body.Password)
		iden := model.NewIdentityEmail(body.Email, uIID, hash)
		if err := iden.Insert(tx); err != nil {
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

	type Tokens struct {
		AccessToken  *string `json:"accessToken"`
		RefreshToken *string `json:"refreshToken"`
	}

	type Response struct {
		Tokens Tokens `json:"tokens"`
	}

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		logger := log.ForRequest(r)

		var body Request
		if err := request.Decode(r, &body); err != nil {
			return respond.GenericServerError(err)
		}

		account, err := model.GetAccount(env.DB, body.Email)
		if err != nil {
			logger.Error(err)
			return respond.Error(err, http.StatusNotFound, "Account not found", "no_account")
		}

		ident, err := model.GetIdentityEmail(env.DB, account.Email)
		if err != nil {
			logger.Error(err)
			return respond.Error(err, http.StatusNotFound, "Account not found", "no_account")
		}

		if err := compareHashAndPassword(ident.PWHash, body.Password); err != nil {
			logger.Error(err)
			return respond.Error(err, http.StatusUnauthorized, "Wrong credentials", "invalid_credentials")
		}

		claims := defaultClaims()
		jwt, err := signedJWTWithClaims(claims)
		if err != nil {
			logger.Error(err)
			return respond.GenericServerError(err)

		}

		return respond.Success(http.StatusOK, Response{
			Tokens: Tokens{
				AccessToken: &jwt,
			},
		})
	})
}
