package http

import (
	"context"
	"fmt"
	"log"
	"syscall"

	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"

	"github.com/aristat/golang-gin-oauth2-example-app/app/db"

	"github.com/aristat/golang-gin-oauth2-example-app/app/session"

	"github.com/aristat/golang-gin-oauth2-example-app/users"
	"github.com/fvbock/endless"

	"github.com/gin-gonic/gin"
)

// Config
type Config struct {
	Debug bool
	Bind  string
}

// Http
type Http struct {
	ctx   context.Context
	cfg   Config
	oauth *oauth.OAuth
}

// ListenAndServe
func (m *Http) ListenAndServe(bind ...string) (err error) {
	bindAdrr := m.cfg.Bind

	if len(bind) > 0 && len(bind[0]) > 0 {
		bindAdrr = bind[0]
	}

	dbManager, _, e := db.Build()
	if e != nil {
		return e
	}

	sessionManager, _, e := session.Build()
	if e != nil {
		return e
	}

	usersService := &users.Service{
		SessionManager: sessionManager,
		DB:             dbManager.DB,
		OauthServer:    m.oauth.OauthServer,
	}

	fmt.Println("init gin")
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/**/*")

	users.Run(r, usersService)
	oauth.Run(r, m.oauth.OauthService)

	fmt.Printf("start listen and serve http at %v", bindAdrr)

	server := endless.NewServer(bindAdrr, r)
	server.BeforeBegin = func(add string) {
		log.Printf("[INFO] Actual pid is %d", syscall.Getpid())
	}

	go func() {
		<-m.ctx.Done()
		fmt.Println("context cancelled, shutdown is raised")
		if e := server.Shutdown(context.Background()); e != nil {
			fmt.Printf("graceful shutdown error, %v", e)
		}
	}()

	err = server.ListenAndServe()
	if err != nil {
		log.Printf("[ERROR] Server err: %v", err)
	}

	return
}

// New
func New(ctx context.Context, cfg Config, oauth *oauth.OAuth) *Http {
	return &Http{
		ctx:   ctx,
		cfg:   cfg,
		oauth: oauth,
	}
}
