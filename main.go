package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/requestid"

	"github.com/tundrawork/stargate/app/common"
	"github.com/tundrawork/stargate/app/common/matomo"
	"github.com/tundrawork/stargate/app/railgun_cdn"
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
	server.OnShutdown = append(server.OnShutdown, func(ctx context.Context) {
		matomo.Shutdown(ctx)
	})
	common.Init()
	railgun_cdn.Init()
}
