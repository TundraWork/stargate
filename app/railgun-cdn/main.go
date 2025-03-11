package railgun_cdn

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/tencentyun/cos-go-sdk-v5"
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

// GetBucket lists all objects in a bucket.
func GetBucket(ctx context.Context, c *app.RequestContext) {
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
	prefix := rootPath
	resp, err := api.GetBucket(ctx, prefix)
	if err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
	}
	c.JSON(consts.StatusOK, common.APIResponseSuccess(resp))
}

// HeadObject returns the metadata of an object.
func HeadObject(ctx context.Context, c *app.RequestContext) {
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
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	objectKey := rootPath + tenantRequest.ObjectPath
	resp, err := api.HeadObject(ctx, objectKey)
	if err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
	}
	c.JSON(consts.StatusOK, common.APIResponseSuccess(resp))
}

// PutObject uploads an object.
func PutObject(ctx context.Context, c *app.RequestContext) {
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
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	objectKey := rootPath + tenantRequest.ObjectPath
	contentType := string(c.GetHeader("Content-Type"))
	resp, err := api.PutObject(ctx, objectKey, c.RequestBodyStream(), contentType, tenantRequest.TTL)
	if err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
	}
	c.JSON(consts.StatusOK, common.APIResponseSuccess(resp))
}

// DeleteObject deletes an object.
func DeleteObject(ctx context.Context, c *app.RequestContext) {
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
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	objectKey := rootPath + tenantRequest.ObjectPath
	if err := api.DeleteObject(ctx, objectKey); err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
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
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	objectKey := "/" + rootPath + tenantRequest.ObjectPath
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
