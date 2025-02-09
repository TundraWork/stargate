package router

import (
	"github.com/tundrawork/stargate/app/common/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// apiRouteRegister registers all API routes.
func apiRouteRegister(r *server.Hertz) {
	r.GET("/ping", handler.Ping)
}
