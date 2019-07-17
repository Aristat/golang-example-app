package oauth

import (
	"context"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/go-oauth2/oauth2/models"

	"github.com/go-session/session"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/store"

	oauthRedis "gopkg.in/go-oauth2/redis.v3"
)

const prefix = "app.oauth"

// OAuth
type OAuth struct {
	ctx          context.Context
	cfg          Config
	log          logger.Logger
	OauthServer  IServer
	OauthService *Routers
}

var ClientsConfig = map[string]oauth2.ClientInfo{
	"123456": &models.Client{
		ID:     "123456",
		Secret: "12345678",
		Domain: "http://localhost:9094",
	},
}

func NewClientStore(config map[string]oauth2.ClientInfo) *store.ClientStore {
	clientStore := store.NewClientStore()
	for key, value := range config {
		clientStore.Set(key, value)
	}

	return clientStore
}

// New
func New(ctx context.Context, log logger.Logger, cfg Config, session *session.Manager) *OAuth {
	oauthConfig := oauthRedis.Options{
		Addr: cfg.RedisUrl,
		DB:   cfg.RedisDB,
	}

	oauth2Service := &Service{
		TokenStore:     oauthRedis.NewRedisStore(&oauthConfig),
		ClientStore:    NewClientStore(ClientsConfig),
		SessionManager: session,
	}

	oauthServer := NewOauthServer(oauth2Service, log)
	authService := &Routers{
		SessionManager: session,
		IServer:        oauthServer,
	}

	return &OAuth{
		ctx:          ctx,
		cfg:          cfg,
		log:          log.WithFields(logger.Fields{"service": prefix}),
		OauthServer:  oauthServer,
		OauthService: authService,
	}
}
