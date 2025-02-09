package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/requestid"
	"github.com/tundrawork/stargate/app/common/service"
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
	c.JSON(consts.StatusOK, service.APIResponseSuccess(PingResponseData{
		Timestamp: timestamp,
		RequestId: requestId,
	}))
}
