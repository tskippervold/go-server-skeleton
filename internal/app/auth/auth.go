package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/okta/okta-sdk-golang/v2/okta"

	"github.com/tskippervold/golang-base-server/internal/utils/respond"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"

	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/utils/log"
)

type accessToken *jwtverifier.Jwt
type contextKey string
type Error string

func (e Error) Error() string { return string(e) }

const (
	requestAuth = contextKey("_rRequestAuth")
)

const (
	ErrInvalidCredentials = Error("Invalid credentials")
)

var (
	oktaClient *okta.Client
)

func Setup(r *mux.Router, env *env.Env) {
	logger := log.NewLogger()

	configureOAuthWithGoogle()
	configureOAuthWithMicrosoft()
	authHandlers(r, env)

	logger.Info("Configuring Okta client.")
	_, client, err := okta.NewClient(context.TODO())
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	oktaClient = client
}

func ForRequest(r *http.Request) accessToken {
	if t, ok := r.Context().Value(requestAuth).(accessToken); ok {
		return t
	}

	return nil
}

func AuthenticatedAccount(r *http.Request, db *sqlx.DB) (*okta.User, error) {
	log := log.ForRequest(r)

	t, ok := r.Context().Value(requestAuth).(accessToken)
	if ok == false {
		return nil, ErrInvalidCredentials
	}

	uid := fmt.Sprintf("%v", t.Claims["uid"])
	log.Debug("Getting Okta user", uid)

	user, _, err := oktaClient.User.GetUser(context.TODO(), uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.ForRequest(r)

		bearerToken := r.Header.Get("Authorization")
		bearerToken = strings.TrimSpace(bearerToken)
		bearerToken = strings.TrimPrefix(strings.TrimPrefix(bearerToken, "bearer "), "Bearer ")

		if bearerToken == "" {
			err := errors.New("Bearer token is malformed or missing")

			res := respond.Error(err, http.StatusUnauthorized, "Authorization required", "missing_auth")
			res.Write(w)
			return
		}

		token, err := parseAndVerifyJWT(bearerToken)

		if err != nil {
			log.Info(err)
			res := respond.Error(err, http.StatusUnauthorized, "Invalid credentials", "invalid_credentials")
			res.Write(w)
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

func defaultClaims(subject string) jwt.Claims {
	expiresAt := time.Now().Add(time.Minute * 30).Unix()

	return jwt.StandardClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "me",
		Subject:   subject,
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

func parseAndVerifyJWT(tokenString string) (accessToken, error) {
	claims := map[string]string{}
	claims["aud"] = "api://default"
	claims["cid"] = "0oaocuc380xeJPqoh4x6" // ClientID

	setup := jwtverifier.JwtVerifier{
		Issuer:           "https://auth.01pay.no/oauth2/default",
		ClaimsToValidate: claims,
	}

	verifier := setup.New()
	token, err := verifier.VerifyAccessToken(tokenString)
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
