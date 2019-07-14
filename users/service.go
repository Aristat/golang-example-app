package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"

	"fmt"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

var (
	userNotFound = errors.New("10002 user not found")
)

type Service struct {
	SessionManager *session.Manager
	DB             *sql.DB
	oauth.OauthServer
}

func Run(routerGin *gin.Engine, service *Service) {
	routerGin.GET("/login", service.GetLogin)
	routerGin.POST("/login", service.PostLogin)

	routerGin.GET("/auth", service.Auth)
	routerGin.POST("/auth", service.Auth)

	routerGin.GET("/user", service.User)
}

func (service *Service) GetLogin(c *gin.Context) {
	_, err := service.SessionManager.Start(context.Background(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "users/login.html", gin.H{})
}

func (service *Service) PostLogin(c *gin.Context) {
	store, err := service.SessionManager.Start(context.Background(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	currentUser, err := FindByEmail(service.DB, email)
	if err != nil {
		c.HTML(http.StatusNotFound, "users/login.html", gin.H{"errors": []string{userNotFound.Error()}})
		return
	}

	if common.CheckPasswordHash(password, currentUser.EncryptedPassword) == false {
		c.HTML(http.StatusNotFound, "users/login.html", gin.H{"errors": []string{userNotFound.Error()}})
		return
	}

	if currentUser != nil {
		store.Set("LoggedInUserID", fmt.Sprintf("%d", currentUser.ID))
		store.Save()

		c.Header("Location", "/auth")
		c.Status(http.StatusFound)
		return
	}

	c.HTML(http.StatusOK, "users/login.html", gin.H{})
}

func (service *Service) Auth(c *gin.Context) {
	store, err := service.SessionManager.Start(context.Background(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		c.Header("Location", "/login")
		c.Status(http.StatusFound)
		return
	}

	c.HTML(http.StatusOK, "users/auth.html", gin.H{})
}

func (service *Service) User(c *gin.Context) {
	ti, err := service.OauthServer.ValidationBearerToken(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Header("Content-Type", "application/json;charset=UTF-8")
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-store")

	userData := userData{ID: ti.GetUserID(), Scope: ti.GetScope()}
	c.JSON(http.StatusOK, gin.H{"user": userData})
}
