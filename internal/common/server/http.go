package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RunHTTPServer(serverName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serverName).GetString("http-addr")
	if addr == "" {
		panic("empty http address")
	}
	runHTTPServerOnAddr(addr, wrapper)
}

func runHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	serMiddlewares(apiRouter)
	wrapper(apiRouter)
	apiRouter.Group("/api")
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

func serMiddlewares(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("default_server"))

}
