package oauth

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

var (
	invalidReturnUri = errors.New("10003 returnUri is not valid")
)

type Service struct {
	SessionManager *session.Manager
	OauthServer
}

func Run(routerGin *gin.Engine, service *Service) {
	routerGin.GET("/authorize", service.Authorize)
	routerGin.POST("/authorize", service.Authorize)

	routerGin.GET("/token", service.Token)
	routerGin.POST("/token", service.Token)
}

func (service *Service) Authorize(c *gin.Context) {
	cw := c.Writer
	cr := c.Request

	store, err := service.SessionManager.Start(context.Background(), cw, cr)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	params, ok := store.Get("ReturnUri")
	if !ok {
		params = ""
	}

	form, err := url.ParseQuery(params.(string))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(form) != 0 {
		cr.Form = form
	}

	store.Delete("ReturnUri")
	store.Save()

	err = service.OauthServer.HandleAuthorizeRequest(cw, cr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func (service *Service) Token(c *gin.Context) {
	cw := c.Writer
	cr := c.Request

	err := service.OauthServer.HandleTokenRequest(cw, cr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
