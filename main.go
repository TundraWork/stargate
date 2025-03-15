package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/requestid"
	"github.com/tundrawork/stargate/app/common/matomo"
	railgunCDN "github.com/tundrawork/stargate/app/railgun-cdn"
	"github.com/tundrawork/stargate/config"
	"github.com/tundrawork/stargate/router"
)

func main() {
	config.Init()
	h := server.Default(
		server.WithHostPorts(":"+config.Conf.ListenPort),
		server.WithHandleMethodNotAllowed(true),
		server.WithStreamBody(true),
		server.WithMaxRequestBodySize(config.Conf.MaxRequestBodySize),
	)
	h.Use(
		requestid.New(),
	)
	h.LoadHTMLGlob("docs/*")
	initServices(h)
	router.Register(h)
	h.Spin()
}

func initServices(server *server.Hertz) {
	matomo.InitClient(
		config.Conf.Matomo.Endpoint,
		config.Conf.Matomo.SiteID,
		config.Conf.Matomo.AuthToken,
		config.Conf.Matomo.NumWorkers,
		config.Conf.Matomo.BatchSize,
		config.Conf.Matomo.EventBufferSize,
	)
	server.OnShutdown = append(server.OnShutdown, func(ctx context.Context) {
		matomo.Shutdown(ctx)
	})
	railgunCDN.Init()
}
