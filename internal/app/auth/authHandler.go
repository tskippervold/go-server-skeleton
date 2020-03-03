package auth

import (
	"net/http"

	"github.com/tskippervold/golang-base-server/internal/app/model"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/utils/log"
	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"
)

func authHandlers(r *mux.Router, env *env.Env) {
	r.Handle("/signup", signup(env)).Methods("POST")
	r.Handle("/login", login(env)).Methods("POST")
}

func signup(env *env.Env) http.Handler {

	type Request struct {
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Password string `json:"password"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body Request
		if err := request.Decode(r, &body); err != nil {
			respond.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accountExists, err := model.AccountExists(env.DB, body.Email)
		if err != nil {
			respond.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if accountExists {
			respond.Error(w, "Account already exists", http.StatusConflict)
			return
		}

		tx := env.DB.MustBegin()

		account := model.NewAccount(body.Email)
		uIID, err := account.Insert(tx)
		if err != nil {
			tx.Rollback()
			respond.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hash, _ := hashPassword(body.Password)
		iden := model.NewIdentityEmail(body.Email, uIID, hash)
		if err := iden.Insert(tx); err != nil {
			tx.Rollback()
			respond.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			respond.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respond.Created(w, nil)
	})
}

func login(env *env.Env) http.Handler {

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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.ForRequest(r)

		var body Request
		if err := request.Decode(r, &body); err != nil {
			respond.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		account, err := model.GetAccount(env.DB, body.Email)
		if err != nil {
			logger.Error(err)
			respond.Error(w, "Account does not exist", http.StatusNotFound)
			return
		}

		ident, err := model.GetIdentityEmail(env.DB, account.Email)
		if err != nil {
			logger.Error(err)
			respond.Error(w, "Account does not exist", http.StatusNotFound)
			return
		}

		if err := compareHashAndPassword(ident.PWHash, body.Password); err != nil {
			logger.Error(err.Error())
			respond.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims := defaultClaims()
		jwt, err := signedJWTWithClaims(claims)
		if err != nil {
			logger.Error(err.Error())
			respond.InternalError(w)
			return
		}

		respond.Ok(w, Response{
			Tokens: Tokens{
				AccessToken: &jwt,
			},
		})
	})
}
