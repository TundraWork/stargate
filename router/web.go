package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/tundrawork/stargate/app/common"
)

// webRouteRegister registers all API routes.
func webRouteRegister(r *server.Hertz) {
	docs_ := r.Group("/docs")
	docs_.GET("/:file", common.DocsHandler)
}
