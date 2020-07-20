package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-micro/v2/web"
	"github.com/wmsx/xconf/agent-api/config"
	"github.com/wmsx/xconf/agent-api/handler"
	pconfig "github.com/wmsx/xconf/proto/config"
)

func main() {
	var cacheSize int

	service := web.NewService(
		web.Name("wm.sx.web.agent"),
		web.Flags(
			&cli.IntFlag{
				Name:        "cache_size",
				Usage:       "cache size (k)",
				EnvVars:     []string{"XCONF_CACHE_SIZE"},
				Value:       1024 * 1024,
				Destination: nil,
			},
		),
		web.Action(func(context *cli.Context) {
			cacheSize = context.Int("cache_size")
		}),
	)

	if err := service.Init(); err != nil {
		panic(err)
	}

	service.Options().Service.Init(micro.Transport(grpc.NewTransport()))
	client := pconfig.NewConfigService("wm.sx.svc.config", service.Options().Service.Client())

	config.Init(client, cacheSize*1024)
	router := Router()
	service.Handle("/", router)

	if err := service.Run(); err != nil {
		panic(err)
	}
}

func Router() *gin.Engine {
	router := gin.Default()
	r := router.Group("/agent/api/v1")
	r.GET("/config", handler.ReadConfig)
	r.GET("/config/raw", handler.ReadConfigRaw)
	r.GET("/watch", handler.WatchUpdate)
	r.GET("/watch/raw", handler.WatchUpdateRaw)

	return router
}
