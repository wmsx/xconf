package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-micro/web"
	"github.com/wmsx/xconf/admin-api/config"
	"github.com/wmsx/xconf/admin-api/handler"
	pconfig "github.com/wmsx/xconf/proto/config"
)

const name  = "wm.sx.web.xconf"

func main() {
	service := web.NewService(
		web.Name(name),
	)

	if err := service.Init(); err != nil {
		panic(err)
	}

	service.Options().Service.Init(micro.Transport(grpc.NewTransport()))
	client := pconfig.NewConfigService("wm.sx.svc.config", service.Options().Service.Client())

	config.Init(client)
	router := Router()
	service.Handle("/", router)

	if err := service.Run(); err != nil {
		panic(err)
	}
}

func Router() *gin.Engine {
	router := gin.Default()
	router.Static("/xconf/admin/ui", "./dist")
	r := router.Group("/xconf/admin/api/v1")
	r.GET("/apps", handler.ListApps)
	r.GET("/app", handler.QueryApp)
	r.POST("/app", handler.CreateApp)
	r.DELETE("/app", handler.DeleteApp)

	r.GET("/clusters", handler.ListClusters)
	r.GET("/cluster", handler.QueryCluster)
	r.POST("/cluster", handler.CreateCluster)
	r.DELETE("/cluster", handler.DeleteCluster)

	r.GET("/namespaces", handler.ListNamespaces)
	r.GET("/namespace", handler.QueryNamespace)
	r.POST("/namespace", handler.CreateNamespace)
	r.DELETE("/namespace", handler.DeleteNamespace)

	r.POST("/config", handler.UpdateConfig)
	r.POST("/release", handler.Release)
	r.POST("/rollback", handler.Rollback)
	r.GET("/release/history", handler.ListReleaseHistory)

	r.GET("/format", handler.ListSupportedFormat)

	return router
}
