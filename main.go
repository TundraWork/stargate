package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/requestid"
	"github.com/tundrawork/stargate/config"
	"github.com/tundrawork/stargate/router"
)

func main() {
	config.Init()
	h := server.Default(
		server.WithHostPorts(":"+config.Conf.ListenPort),
		server.WithHandleMethodNotAllowed(true),
	)
	h.Use(
		requestid.New(),
	)
	router.Register(h)
	h.Spin()
}
