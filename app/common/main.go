package common

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/requestid"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type PingResponseData struct {
	Timestamp int64  `json:"timestamp"`
	RequestId string `json:"requestId"`
}

// Ping returns environment information of the server.
func Ping(ctx context.Context, c *app.RequestContext) {
	timestamp := time.Now().Unix()
	requestId := requestid.Get(c)
	c.JSON(consts.StatusOK, APIResponseSuccess(PingResponseData{
		Timestamp: timestamp,
		RequestId: requestId,
	}))
}

// DocsHandler handles the request for the API documentation.
func DocsHandler(ctx context.Context, c *app.RequestContext) {
	template := fmt.Sprintf("%s.tmpl", c.Param("file"))
	filename := fmt.Sprintf("docs/%s", template)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		c.SetStatusCode(http.StatusNotFound)
		c.SetBodyString("404 Not Found: The requested documentation file does not exist.")
		return
	}
	c.HTML(http.StatusOK, template, utils.H{
		// custom data
	})
}

// InvalidAPIPathHandler handles the request for invalid API paths.
func InvalidAPIPathHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusNotFound, APIResponseError(consts.StatusNotFound, "The requested API path does not exist."))
}
