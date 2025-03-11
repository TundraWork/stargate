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
	railgun_.GET("/bucket", railgunCDN.GetBucket)
	railgun_.GET("/object", railgunCDN.HeadObject)
	railgun_.PUT("/object", railgunCDN.PutObject)
	railgun_.DELETE("/object", railgunCDN.DeleteObject)
	railgun_.GET("/url", railgunCDN.GetURL)
}
