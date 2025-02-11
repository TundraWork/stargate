package railgun_cdn

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/tundrawork/stargate/app/common"
	"github.com/tundrawork/stargate/app/railgun-cdn/api"
	"github.com/tundrawork/stargate/config"
)

// Init initializes the Railgun CDN service.
func Init() {
	api.InitCosClient(
		config.Conf.Services.RailgunCDN.COS.Bucket,
		config.Conf.Services.RailgunCDN.COS.Region,
		config.Conf.Services.RailgunCDN.COS.SecretID,
		config.Conf.Services.RailgunCDN.COS.SecretKey,
	)
}

// Put uploads an object.
func Put(ctx context.Context, c *app.RequestContext) {
	tenantRequest, err := parseCommonTenantRequest(c)
	if err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
	}
	rootPath, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
	}
	objectKey := rootPath + tenantRequest.ObjectPath
	contentType := string(c.GetHeader("Content-Type"))
	if err := api.PutObject(ctx, objectKey, c.RequestBodyStream(), contentType); err != nil {
		c.JSON(consts.StatusInternalServerError, common.APIResponseError(consts.StatusInternalServerError, err.Error()))
	}
	c.JSON(consts.StatusOK, common.APIResponseSuccess(nil))
}

// parseCommonTenantRequest parses the common request fields of a tenant.
func parseCommonTenantRequest(c *app.RequestContext) (*CommonTenantRequest, error) {
	appID := c.GetHeader("X-App-Id")
	appKey := c.GetHeader("X-App-Key")
	objectPath := c.GetHeader("X-Object-Path")
	if len(appID) == 0 || len(appKey) == 0 || len(objectPath) == 0 {
		return &CommonTenantRequest{}, errors.New("missing common tenant request fields")
	}
	return &CommonTenantRequest{
		AppID:      string(appID),
		AppKey:     string(appKey),
		ObjectPath: string(objectPath),
	}, nil
}

// authTenant authenticates the tenant from the common tenant request and returns the tenant's root path.
func authTenant(req *CommonTenantRequest) (string, error) {
	for _, tenant := range config.Conf.Services.RailgunCDN.Tenants {
		if tenant.AppID == req.AppID && tenant.AppKey == req.AppKey {
			return tenant.RootPath, nil
		}
	}
	return "", errors.New("cannot authorize tenant")
}
