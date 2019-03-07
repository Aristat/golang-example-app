package oauth

import (
	"github.com/go-oauth2/oauth2"
	"github.com/go-oauth2/oauth2/models"
	"gopkg.in/oauth2.v3/store"
)

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
