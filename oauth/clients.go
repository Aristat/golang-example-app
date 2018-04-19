package oauth

import (
	"github.com/go-oauth2/oauth2"
	"github.com/go-oauth2/oauth2/models"
)

var clientsConfig = map[string]oauth2.ClientInfo{
	"123456": &models.Client{
		ID:     "123456",
		Secret: "12345678",
		Domain: "http://localhost:9094",
	},
}
