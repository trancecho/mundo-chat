package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-chat/initialize"
	"github.com/trancecho/mundo-chat/initialize/router"
	"github.com/trancecho/mundo-chat/server"
	"github.com/trancecho/mundo-chat/server/middleware"
)

func GinInit() *gin.Engine {
	r := gin.Default()
	initialize.InitRedis()
	router.GenerateRouter(r)
	middleware.InitSecret()
	//启动房间管理器
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Managers.Start(ctx)
	return r
}
