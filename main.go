package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/aristat/golang-gin-oauth2-example-app/routers"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
	"github.com/fvbock/endless"
	"github.com/spf13/viper"
)

var (
	env *common.Env
)

func init() {
	viper.AutomaticEnv()
	common.InitConfig()
	env = common.InitEnv()
	common.InitSession()
}

func main() {
	port := viper.GetInt("HTTP_SERVER_PORT")
	endPoint := fmt.Sprintf(":%d", port)
	routersInit := routers.Init(env)

	server := endless.NewServer(endPoint, routersInit)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
