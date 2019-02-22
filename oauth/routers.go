package oauth

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
)

var (
	invalidReturnUri = errors.New("10003 returnUri is not valid")
)

type AuthRouters struct {
	*common.Env
	OauthServer
}

func InitRouters(routerGin *gin.Engine, router *AuthRouters) {
	router.OauthServer.UserAuthorizationHandler(userAuthorization)

	routerGin.GET("/authorize", router.Authorize)
	routerGin.POST("/authorize", router.Authorize)

	routerGin.GET("/token", router.Token)
	routerGin.POST("/token", router.Token)
}

func (router *AuthRouters) Authorize(c *gin.Context) {
	cw := c.Writer
	cr := c.Request

	store, err := session.Start(context.Background(), cw, cr)
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

	err = router.OauthServer.HandleAuthorizeRequest(cw, cr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func (router *AuthRouters) Token(c *gin.Context) {
	cw := c.Writer
	cr := c.Request

	err := router.OauthServer.HandleTokenRequest(cw, cr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func userAuthorization(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form.Encode())
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)

	store.Delete("LoggedInUserID")
	store.Save()

	return
}
