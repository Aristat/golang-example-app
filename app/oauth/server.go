package oauth

import (
	"fmt"
	"time"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

func NewOauthServer(oauth2Service *Service, log logger.Logger) *server.Server {
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
		fmt.Println("logger.Logger1", &log)
		log.Error("Response Error: %s", logger.Args(re.Error.Error()))
	})

	return server
}
