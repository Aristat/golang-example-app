package cmd

import (
	"fmt"
	"log"
	"syscall"

	"github.com/fvbock/endless"

	"github.com/gin-gonic/gin"

	"github.com/aristat/golang-gin-oauth2-example-app/users"
	goSession "github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
	"github.com/aristat/golang-gin-oauth2-example-app/oauth"
	"github.com/aristat/golang-gin-oauth2-example-app/oauth/session"
	"github.com/go-session/redis"
	oauthRedis "gopkg.in/go-oauth2/redis.v3"
)

type ServerCommand struct {
	DB    DatabaseGroup `group:"database" namespace:"database" env-namespace:"DATABASE" required:"true"`
	Redis RedisGroup    `group:"redis" namespace:"redis" env-namespace:"REDIS" required:"true"`

	Port int `long:"port" env:"PORT" default:"9096" description:"port"`
}

type DatabaseGroup struct {
	URL string `long:"url" env:"URL" description:"database url" default:"postgresql://localhost:5432/oauth2_development?sslmode=disable"`
}

type RedisGroup struct {
	URL       string `long:"url" env:"URL" description:"redis url" default:"127.0.0.1:6379"`
	SessionDB int    `long:"session.db" env:"SESSION_DB" description:"redis session db" default:"10"`
	TokenDB   int    `long:"token.db" env:"TOKEN_DB" description:"redis token db" default:"10"`
}

func (cmd *ServerCommand) Execute(args []string) error {
	db := common.InitDB(cmd.DB.URL)

	sessionStore := redis.NewRedisStore(&redis.Options{
		Addr: cmd.Redis.URL,
		DB:   cmd.Redis.SessionDB,
	})
	sessionManager := session.Init(sessionStore)
	oauthServer := cmd.makeOAuthServer(sessionManager)

	usersService := &users.Service{
		SessionManager: sessionManager,
		DB:             db,
		OauthServer:    oauthServer,
	}

	authService := &oauth.Service{
		SessionManager: sessionManager,
		OauthServer:    oauthServer,
	}

	endPoint := fmt.Sprintf(":%d", cmd.Port)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/**/*")

	users.Run(r, usersService)
	oauth.Run(r, authService)

	server := endless.NewServer(endPoint, r)
	server.BeforeBegin = func(add string) {
		log.Printf("[INFO] Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("[ERROR] Server err: %v", err)
	}

	return nil
}

func (cmd *ServerCommand) makeOAuthServer(sessionManager *goSession.Manager) oauth.OauthServer {
	oauthConfig := oauthRedis.Options{
		Addr: cmd.Redis.URL,
		DB:   cmd.Redis.TokenDB,
	}

	oauth2Service := &oauth.Oauth2Service{
		TokenStore:     oauthRedis.NewRedisStore(&oauthConfig),
		ClientStore:    oauth.NewClientStore(oauth.ClientsConfig),
		SessionManager: sessionManager,
	}

	return oauth.NewOauthServer(oauth2Service)
}
