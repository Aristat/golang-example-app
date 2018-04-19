package users

import (
	"context"
	"net/http"

	"fmt"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

type UserRouters struct {
	*common.Env
	common.OauthServer
}

func InitRouters(routerGin *gin.Engine, router *UserRouters) {
	routerGin.GET("/login", router.GetLogin)
	routerGin.POST("/login", router.PostLogin)

	routerGin.GET("/auth", router.Auth)
	routerGin.POST("/auth", router.Auth)

	routerGin.GET("/user", router.User)
}

func (router *UserRouters) GetLogin(c *gin.Context) {
	_, err := session.Start(context.Background(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "users/login.html", gin.H{})
}

func (router *UserRouters) PostLogin(c *gin.Context) {
	store, err := session.Start(context.Background(), c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	currentUser, err := FindByEmail(router.Env, email)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if common.CheckPasswordHash(password, currentUser.EncryptedPassword) == false {
		c.HTML(http.StatusNotFound, "users/login.html", gin.H{"errors": []string{common.UserNotFound.Error()}})
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

func (router *UserRouters) Auth(c *gin.Context) {
	store, err := session.Start(context.Background(), c.Writer, c.Request)
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

func (router *UserRouters) User(c *gin.Context) {
	ti, err := router.OauthServer.ValidationBearerToken(c.Request)
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
