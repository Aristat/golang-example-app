package common

import (
	"net/http"

	oauth2 "gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/server"
)

type H map[string]interface{}

type OauthServer interface {
	UserAuthorizationHandler(handler server.UserAuthorizationHandler)
	HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) (err error)
	HandleTokenRequest(w http.ResponseWriter, r *http.Request) (err error)
	ValidationBearerToken(r *http.Request) (ti oauth2.TokenInfo, err error)
}

type ExternalService interface {
	Login(email, password string) (*http.Response, error)
}
