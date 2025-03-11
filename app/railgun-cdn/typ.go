package railgun_cdn

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

type CommonTenantRequest struct {
	AppID      string
	AppKey     string
	ObjectPath string
	TTL        int64
}

type GetURLResponse struct {
	URL     string `json:"url"`
	Expires int64  `json:"expires"`
}

// FromRequestContext extracts the common tenant request fields from the request context.
func (req *CommonTenantRequest) FromRequestContext(c *app.RequestContext) error {
	appID := c.GetHeader("X-App-Id")
	appKey := c.GetHeader("X-App-Key")
	objectPath := c.GetHeader("X-Object-Path")

	if len(appID) == 0 || len(appKey) == 0 {
		return errors.New("missing common tenant request fields")
	}
	if len(objectPath) > 0 && !isValidObjectPath(string(objectPath)) {
		return errors.New("invalid object path")
	}
	ttlStr := string(c.GetHeader("X-TTL"))
	var ttl int64
	var err error
	if len(ttlStr) > 0 {
		ttl, err = strconv.ParseInt(ttlStr, 10, 64)
		if err != nil {
			return errors.New("invalid X-TTL value")
		}
		if ttl <= 0 {
			return errors.New("X-TTL value must be greater than 0")
		}
	}

	req.AppID = string(appID)
	req.AppKey = string(appKey)
	req.ObjectPath = string(objectPath)
	req.TTL = ttl

	return nil
}
