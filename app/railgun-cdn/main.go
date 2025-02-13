package railgun_cdn

import (
	"context"
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
	tenantRequest := &CommonTenantRequest{}
	if err := tenantRequest.FromRequestContext(c); err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
		return
	}
	rootPath, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	objectKey := rootPath + tenantRequest.ObjectPath
	contentType := string(c.GetHeader("Content-Type"))
	if err := api.PutObject(ctx, objectKey, c.RequestBodyStream(), contentType, tenantRequest.TTL); err != nil {
		c.JSON(consts.StatusInternalServerError, common.APIResponseError(consts.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, common.APIResponseSuccess(nil))
}

// Delete deletes an object.
func Delete(ctx context.Context, c *app.RequestContext) {
	tenantRequest := &CommonTenantRequest{}
	if err := tenantRequest.FromRequestContext(c); err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
		return
	}
	rootPath, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	objectKey := rootPath + tenantRequest.ObjectPath
	if err := api.DeleteObject(ctx, objectKey); err != nil {
		c.JSON(consts.StatusInternalServerError, common.APIResponseError(consts.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, common.APIResponseSuccess(nil))
}

// GetURL returns the signed URL to access an object.
func GetURL(ctx context.Context, c *app.RequestContext) {
	tenantRequest := &CommonTenantRequest{}
	if err := tenantRequest.FromRequestContext(c); err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
		return
	}
	rootPath, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	if !isValidObjectPath(tenantRequest.ObjectPath) {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "invalid object path"))
		return
	}
	objectKey := rootPath + tenantRequest.ObjectPath
	url, expires, err := api.GetObjectPublicURL(objectKey, tenantRequest.TTL)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, common.APIResponseError(consts.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(consts.StatusOK, common.APIResponseSuccess(GetURLResponse{
		URL:     url,
		Expires: expires,
	}))
}
