package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/tundrawork/stargate/app/common"
)

// apiRouteRegister registers all API routes.
func apiRouteRegister(r *server.Hertz) {
	common_ := r.Group("/common/v1")
	common_.GET("/ping", common.Ping)

	railgun_ := r.Group("/railgun/v1")
	railgun_.GET()
}
