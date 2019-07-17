package oauth

import (
	"net/http"
	"time"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

type IServer interface {
	HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) (err error)
	HandleTokenRequest(w http.ResponseWriter, r *http.Request) (err error)
	ValidationBearerToken(r *http.Request) (ti oauth2.TokenInfo, err error)
}

type oauth2Server struct {
	*server.Server
}

func NewServer(srv *server.Server) IServer {
	return &oauth2Server{Server: srv}
}

func NewOauthServer(oauth2Service *Service, log logger.Logger) IServer {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(
		&manage.Config{
			AccessTokenExp:    time.Hour * 24 * 7,
			RefreshTokenExp:   time.Hour * 24 * 14,
			IsGenerateRefresh: true,
		},
	)

	manager.MapTokenStorage(oauth2Service.TokenStore)
	manager.MapClientStorage(oauth2Service.ClientStore)

	server := server.NewDefaultServer(manager)
	server.UserAuthorizationHandler = userAuthorization(oauth2Service)
	server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Error("Internal Error: %s", logger.Args(err.Error()))
		return
	})
	server.SetResponseErrorHandler(func(re *errors.Response) {
		log.Error("Response Error: %s", logger.Args(re.Error.Error()))
	})

	return NewServer(server)
}
