package oauth

import (
	"time"

	"gopkg.in/oauth2.v3/store"

	"gopkg.in/oauth2.v3"

	"github.com/aristat/golang-example-app/app/logger"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

func NewServer(log logger.Logger, tokenStore oauth2.TokenStore, clientStore *store.ClientStore) *server.Server {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(
		&manage.Config{
			AccessTokenExp:    time.Hour * 24 * 7,
			RefreshTokenExp:   time.Hour * 24 * 14,
			IsGenerateRefresh: true,
		},
	)

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	s := server.NewDefaultServer(manager)

	s.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Error("Internal Error: %s", logger.Args(err.Error()))
		return
	})
	s.SetResponseErrorHandler(func(re *errors.Response) {
		log.Error("Response Error: %s", logger.Args(re.Error.Error()))
	})

	return s
}
