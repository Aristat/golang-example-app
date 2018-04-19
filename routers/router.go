package routers

import (
	"github.com/aristat/golang-gin-oauth2-example-app/common"
	"github.com/aristat/golang-gin-oauth2-example-app/oauth"
	"github.com/aristat/golang-gin-oauth2-example-app/users"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init(env *common.Env) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob(viper.GetString("TEMPLATE_PATH"))

	userRouters := &users.UserRouters{
		Env:         env,
		OauthServer: oauth.GetIOauthServer(),
	}
	users.InitRouters(r, userRouters)

	authRouters := &oauth.AuthRouters{
		Env:         env,
		OauthServer: oauth.GetIOauthServer(),
	}
	oauth.InitRouters(r, authRouters)

	return r
}
