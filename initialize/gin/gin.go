package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-chat/initialize"
	"github.com/trancecho/mundo-chat/initialize/router"
	"github.com/trancecho/mundo-chat/server/middleware"
)

func GinInit() *gin.Engine {
	r := gin.Default()
	initialize.InitRedis()
	router.GenerateRouter(r)
	middleware.InitSecret()
	return r
}
