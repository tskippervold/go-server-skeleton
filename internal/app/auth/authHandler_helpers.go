package auth

import (
	"net/http"

	"github.com/tskippervold/golang-base-server/internal/utils/respond"
)

type Tokens struct {
	AccessToken  *string `json:"accessToken"`
	RefreshToken *string `json:"refreshToken"`
}

type LoginResponse struct {
	Tokens Tokens `json:"tokens"`
}

func loginResponse(accountID string) *respond.Response {
	claims := defaultClaims(accountID)

	jwt, err := signedJWTWithClaims(claims)
	if err != nil {
		return respond.GenericServerError(err)
	}

	return respond.Success(http.StatusOK, LoginResponse{
		Tokens: Tokens{
			AccessToken: &jwt,
		},
	})
}
