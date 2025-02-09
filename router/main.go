package router

import "github.com/cloudwego/hertz/pkg/app/server"

// Register registers all routes.
func Register(r *server.Hertz) {
	apiRouteRegister(r)
	webRouteRegister(r)
}
