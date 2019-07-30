package oauth

import (
	"context"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
	"github.com/go-oauth2/oauth2/models"

	"github.com/go-session/session"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/store"
)

const prefix = "app.oauth"

// OAuth
type OAuth struct {
	ctx          context.Context
	Logger       logger.Logger
	OauthService *Routers
}

var ClientsConfig = map[string]oauth2.ClientInfo{
	"123456": &models.Client{
		ID:     "123456",
		Secret: "12345678",
		Domain: "http://127.0.0.1:8090",
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
func New(ctx context.Context, log logger.Logger, tokenStore oauth2.TokenStore, session *session.Manager) *OAuth {
	oauth2Service := &Service{
		TokenStore:     tokenStore,
		ClientStore:    NewClientStore(ClientsConfig),
		SessionManager: session,
	}

	oauthServer := NewOauthServer(oauth2Service, log)
	authService := &Routers{
		Server:        oauthServer,
		OauthService2: oauth2Service,
	}

	return &OAuth{
		ctx:          ctx,
		Logger:       log.WithFields(logger.Fields{"service": prefix}),
		OauthService: authService,
	}
}
