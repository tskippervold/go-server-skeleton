package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

var (
	googleOAuthConfig    *oauth2.Config
	microsoftOAuthConfig *oauth2.Config
)

type OAuthProvider string

const (
	OAuthProviderGoogle    OAuthProvider = "google"
	OAuthProviderMicrosoft OAuthProvider = "microsoft"
)

const (
	keyGoogleOauthStateCookie string = "oauthstate"
	keyOauthPassthroughCookie string = "passthrough"
	urlGoogleAPI              string = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	urlMicrosoftAPI           string = "https://graph.microsoft.com/v1.0"
)

type credentials struct {
	id     string `json:"cid"`
	secret string `json:"csecret"`
}

type oAuthGoogle struct {
	ID              string `json:"id`
	PictureURL      string `json:"picture"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"verified_email"`
}

type oAuthMicrosoft struct {
	ID    string `json:"id"`
	Email string `json:"mail"`
}

type OAuthResult struct {
	Provider        OAuthProvider `json:"provider"`
	ID              string        `json:"id"`
	Email           string        `json:"email"`
	JSONPassthrough string
}

func configureOAuthWithGoogle() {
	cred := credentials{
		id:     "208084384794-h2dr3ri4g8ouhun61gd36koh7i76eig4.apps.googleusercontent.com",
		secret: "bbCTXSdrg35rAINUMgMlBtbg",
	}

	googleOAuthConfig = &oauth2.Config{
		ClientID:     cred.id,
		ClientSecret: cred.secret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:3000/auth/oauth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}
}

func configureOAuthWithMicrosoft() {
	cred := credentials{
		id:     "44a0f725-cee8-412c-b6c9-6db6d8602574",
		secret: "3c=gmK?i531o-XMORSmM/rUb4[@l]0Uz",
	}

	tenant := "342a0539-37c3-4c36-8188-d8bd916d7de5"

	microsoftOAuthConfig = &oauth2.Config{
		ClientID:     cred.id,
		ClientSecret: cred.secret,
		Endpoint:     microsoft.AzureADEndpoint(tenant),
		RedirectURL:  "http://localhost:3000/auth/oauth/microsoft/callback",
		Scopes: []string{
			"openid",
			"offline_access",
		},
	}
}

func generateStateOauthCookie(w http.ResponseWriter, passthroughValue interface{}) string {
	// TODO: Is this waaaay to long?
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	stateCookie := http.Cookie{
		Name:    keyGoogleOauthStateCookie,
		Value:   state,
		Expires: expiration,
	}

	jsonString, _ := json.Marshal(passthroughValue)
	passthrough := base64.URLEncoding.EncodeToString(jsonString)

	passthroghCookie := http.Cookie{
		Name:    keyOauthPassthroughCookie,
		Value:   passthrough,
		Expires: expiration,
	}

	http.SetCookie(w, &stateCookie)
	http.SetCookie(w, &passthroghCookie)

	return state
}

func OAuthAuthenticate(provider OAuthProvider, w http.ResponseWriter, r *http.Request, passthroughValue interface{}) error {
	switch provider {
	case OAuthProviderGoogle:
		state := generateStateOauthCookie(w, passthroughValue)
		url := googleOAuthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	case OAuthProviderMicrosoft:
		state := generateStateOauthCookie(w, passthroughValue)
		url := microsoftOAuthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	default:
		return errors.New("Invalid oauth provider")
	}
}

func OAuthCallback(provider OAuthProvider, r *http.Request) (*OAuthResult, error) {
	state, err := r.Cookie(keyGoogleOauthStateCookie)
	if err != nil {
		return nil, err
	}

	if r.FormValue("state") != state.Value {
		// Invalid oauth state.
		return nil, errors.New("Invalid OAuth state")
	}

	errorCode := r.FormValue("error")
	if errorCode != "" {
		err := errors.New(r.FormValue("error_description"))
		return nil, err
	}

	code := r.FormValue("code")
	var oauthResult *OAuthResult

	switch provider {
	case OAuthProviderGoogle:
		token, err := googleOAuthConfig.Exchange(context.Background(), code)
		if err != nil {
			return nil, err
		}

		resp, err := http.Get(urlGoogleAPI + token.AccessToken)
		if err != nil {
			return nil, err
		}

		var g oAuthGoogle
		if err := request.Decode(resp.Body, &g); err != nil {
			return nil, err
		}

		oauthResult = &OAuthResult{
			Provider: provider,
			ID:       g.ID,
			Email:    g.Email,
		}
	case OAuthProviderMicrosoft:
		token, err := microsoftOAuthConfig.Exchange(context.Background(), code)
		if err != nil {
			return nil, err
		}

		req, _ := http.NewRequest("GET", urlMicrosoftAPI+"/me/", nil)
		req.Header.Add("Authorization", "bearer "+token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			return nil, err
		}

		var m oAuthMicrosoft
		if err := request.Decode(resp.Body, &m); err != nil {
			return nil, err
		}

		oauthResult = &OAuthResult{
			Provider: provider,
			ID:       m.ID,
			Email:    m.Email,
		}
	default:
		return nil, errors.New("Invalid oauth provider")
	}

	if oauthResult.Email == "" {
		return nil, errors.New("Missing oauth email from provider " + string(provider))
	}

	if oauthResult.ID == "" {
		return nil, errors.New("Missing oauth ID from provider " + string(provider))
	}

	if passthroughValue, _ := r.Cookie(keyOauthPassthroughCookie); passthroughValue.Value != "" {
		decoded, err := base64.URLEncoding.DecodeString(passthroughValue.Value)
		if err != nil {
			return nil, err
		}

		oauthResult.JSONPassthrough = string(decoded)
	}

	return oauthResult, nil
}
