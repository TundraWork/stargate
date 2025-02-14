package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/tundrawork/stargate/app/common"
	railgunCDN "github.com/tundrawork/stargate/app/railgun-cdn"
)

// apiRouteRegister registers all API routes.
func apiRouteRegister(r *server.Hertz) {
	r.NoMethod(common.InvalidAPIPathHandler)

	common_ := r.Group("/common/v1")
	common_.GET("/ping", common.Ping)

	railgun_ := r.Group("/railgun/v1")
	railgun_.PUT("/object", railgunCDN.Put)
	railgun_.DELETE("/object", railgunCDN.Delete)
	railgun_.GET("/object", railgunCDN.GetURL)
}
