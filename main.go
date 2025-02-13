package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/requestid"
	railgunCDN "github.com/tundrawork/stargate/app/railgun-cdn"
	"github.com/tundrawork/stargate/config"
	"github.com/tundrawork/stargate/router"
)

func main() {
	config.Init()
	initServices()
	h := server.Default(
		server.WithHostPorts(":"+config.Conf.ListenPort),
		server.WithHandleMethodNotAllowed(true),
		server.WithStreamBody(true),
		server.WithMaxRequestBodySize(config.Conf.MaxRequestBodySize),
	)
	h.Use(
		requestid.New(),
	)
	router.Register(h)
	h.Spin()
}

func initServices() {
	railgunCDN.Init()
}
