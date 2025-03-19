package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/tundrawork/stargate/app/common"
	"github.com/tundrawork/stargate/app/railgun_cdn"
)

// apiRouteRegister registers all API routes.
func apiRouteRegister(r *server.Hertz) {
	r.NoMethod(common.InvalidAPIPathHandler)

	common_ := r.Group("/common/v1")
	common_.GET("/ping", common.Ping)

	railgun_ := r.Group("/railgun/v1")
	railgun_.GET("/bucket", railgun_cdn.GetBucket)
	railgun_.GET("/object", railgun_cdn.HeadObject)
	railgun_.PUT("/object", railgun_cdn.PutObject)
	railgun_.DELETE("/object", railgun_cdn.DeleteObject)
	railgun_.GET("/url", railgun_cdn.GetURL)
	railgun_.GET("/gateway", railgun_cdn.ClientGateway)
}
