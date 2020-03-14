package auth

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"

	"github.com/tskippervold/golang-base-server/internal/utils/log"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"

	"github.com/gorilla/mux"

	env "github.com/tskippervold/golang-base-server/internal/app"
)

type contextKey string
type Error string

func (e Error) Error() string { return string(e) }

const (
	requestAuth = contextKey("_rRequestAuth")
)

const (
	ErrInvalidCredentials = Error("Invalid credentials")
)

func Setup(r *mux.Router, env *env.Env) {
	authHandlers(r, env)
}

func ForRequest(r *http.Request) *jwt.Token {
	if t, ok := r.Context().Value(requestAuth).(*jwt.Token); ok {
		return t
	}

	return nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.ForRequest(r)

		bearerToken := r.Header.Get("Authorization")
		bearerToken = strings.TrimSpace(bearerToken)
		bearerToken = strings.TrimPrefix(strings.TrimPrefix(bearerToken, "bearer "), "Bearer ")

		if bearerToken == "" {
			respond.Unauthorized(w)
			return
		}

		token, err := parseAndVerifyJWT(bearerToken)
		if err != nil {
			log.Info(err.Error())
			respond.Forbidden(w)
			return
		}

		ctx := context.WithValue(r.Context(), requestAuth, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func privateKey() (*rsa.PrivateKey, error) {
	path := "../keys/private.pem"
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func publicKey() (*rsa.PublicKey, error) {
	path := "../keys/public.pem"
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func defaultClaims() jwt.Claims {
	expiresAt := time.Now().Add(time.Second * 120).Unix()

	return jwt.StandardClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "me",
	}
}

func signedJWTWithClaims(c jwt.Claims) (string, error) {
	ecdsaKey, err := privateKey()
	if err != nil {
		return "", err
	}

	signingMethod := jwt.SigningMethodRS256
	token := jwt.NewWithClaims(signingMethod, c)
	signedToken, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func parseAndVerifyJWT(str string) (*jwt.Token, error) {
	ecdsaKey, err := publicKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		return ecdsaKey, nil
	})

	// Verify token...

	if err != nil {
		return nil, err
	}

	return token, nil
}

func hashPassword(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}

func compareHashAndPassword(h []byte, p string) error {
	if err := bcrypt.CompareHashAndPassword(h, []byte(p)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrInvalidCredentials
		}

		return err
	}

	return nil
}