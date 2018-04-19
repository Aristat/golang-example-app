package oauth

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
)

type AuthRouters struct {
	*common.Env
	common.OauthServer
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

	form := url.Values{}
	if v, ok := store.Get("ReturnUri"); ok {

		if rec, ok := v.(map[string]interface{}); ok {
			for key, value := range rec {

				aInterface := value.([]interface{})
				aString := make([]string, len(aInterface))
				for i, v := range aInterface {
					aString[i] = v.(string)
				}

				for _, v := range aString {
					form.Add(key, v)
				}
			}
		} else {
			log.Printf("Value not a map[string]interface{}: %v\n", v)

			c.AbortWithError(http.StatusInternalServerError, common.InvalidReturnUri)
			return
		}
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

		store.Set("ReturnUri", r.Form)
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
