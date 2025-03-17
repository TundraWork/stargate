package railgun_cdn

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tundrawork/stargate/app/common"
	"github.com/tundrawork/stargate/app/common/matomo"
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
	tenant, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	hlog.CtxInfof(ctx, "[RailgunCDN][Request] Method=%s AppID=%s", "GetBucket", tenantRequest.AppID)
	prefix := tenant.RootPath
	resp, err := api.GetBucket(ctx, prefix)
	if err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			hlog.CtxErrorf(ctx, "[RailgunCDN][Error] Method=%s AppID=%s StatusCode=%d", "GetBucket", tenantRequest.AppID, cosErr.Response.StatusCode)
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
	}
	matomo.ReportEvent(ctx, matomo.Event{
		SiteID:     tenant.SiteID,
		ActionName: "railgun-cdn:server:GetBucket",
		URL:        config.Conf.Services.RailgunCDN.CDN.Endpoint + tenantRequest.ObjectPath,
		UserAgent:  string(c.UserAgent()),
		ClientIP:   c.ClientIP(),
		ClientTime: time.Now(),
	})
	c.JSON(consts.StatusOK, common.APIResponseSuccess(resp))
}

// HeadObject returns the metadata of an object.
func HeadObject(ctx context.Context, c *app.RequestContext) {
	tenantRequest := &CommonTenantRequest{}
	if err := tenantRequest.FromRequestContext(c); err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
		return
	}
	tenant, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	hlog.CtxInfof(ctx, "[RailgunCDN][Request] Method=%s AppID=%s ObjectPath=%s", "HeadObject", tenantRequest.AppID, tenantRequest.ObjectPath)
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	objectKey := tenant.RootPath + tenantRequest.ObjectPath
	resp, err := api.HeadObject(ctx, objectKey)
	if err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			hlog.CtxErrorf(ctx, "[RailgunCDN][Error] Method=%s AppID=%s StatusCode=%d", "HeadObject", tenantRequest.AppID, cosErr.Response.StatusCode)
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
	}
	matomo.ReportEvent(ctx, matomo.Event{
		SiteID:     tenant.SiteID,
		ActionName: "railgun-cdn:server:HeadObject",
		URL:        config.Conf.Services.RailgunCDN.CDN.Endpoint + tenantRequest.ObjectPath,
		UserAgent:  string(c.UserAgent()),
		ClientIP:   c.ClientIP(),
		ClientTime: time.Now(),
	})
	c.JSON(consts.StatusOK, common.APIResponseSuccess(resp))
}

// PutObject uploads an object.
func PutObject(ctx context.Context, c *app.RequestContext) {
	tenantRequest := &CommonTenantRequest{}
	if err := tenantRequest.FromRequestContext(c); err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
		return
	}
	tenant, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	hlog.CtxInfof(ctx, "[RailgunCDN][Request] Method=%s AppID=%s ObjectPath=%s", "PutObject", tenantRequest.AppID, tenantRequest.ObjectPath)
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	objectKey := tenant.RootPath + tenantRequest.ObjectPath
	contentType := string(c.GetHeader("Content-Type"))
	resp, err := api.PutObject(ctx, objectKey, c.RequestBodyStream(), contentType, tenantRequest.TTL)
	if err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			hlog.CtxErrorf(ctx, "[RailgunCDN][Error] Method=%s AppID=%s StatusCode=%d", "PutObject", tenantRequest.AppID, cosErr.Response.StatusCode)
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
	}
	matomo.ReportEvent(ctx, matomo.Event{
		SiteID:     tenant.SiteID,
		ActionName: "railgun-cdn:server:PutObject",
		URL:        config.Conf.Services.RailgunCDN.CDN.Endpoint + tenantRequest.ObjectPath,
		UserAgent:  string(c.UserAgent()),
		ClientIP:   c.ClientIP(),
		ClientTime: time.Now(),
	})
	c.JSON(consts.StatusOK, common.APIResponseSuccess(resp))
}

// DeleteObject deletes an object.
func DeleteObject(ctx context.Context, c *app.RequestContext) {
	tenantRequest := &CommonTenantRequest{}
	if err := tenantRequest.FromRequestContext(c); err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
		return
	}
	tenant, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	hlog.CtxInfof(ctx, "[RailgunCDN][Request] Method=%s AppID=%s ObjectPath=%s", "DeleteObject", tenantRequest.AppID, tenantRequest.ObjectPath)
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	objectKey := tenant.RootPath + tenantRequest.ObjectPath
	if err := api.DeleteObject(ctx, objectKey); err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) {
			hlog.CtxErrorf(ctx, "[RailgunCDN][Error] Method=%s AppID=%s StatusCode=%d", "DeleteObject", tenantRequest.AppID, cosErr.Response.StatusCode)
			c.JSON(cosErr.Response.StatusCode, common.APIResponseError(cosErr.Response.StatusCode, cosErr.Response.Status))
			return
		}
	}
	matomo.ReportEvent(ctx, matomo.Event{
		SiteID:     tenant.SiteID,
		ActionName: "railgun-cdn:server:DeleteObject",
		URL:        config.Conf.Services.RailgunCDN.CDN.Endpoint + tenantRequest.ObjectPath,
		UserAgent:  string(c.UserAgent()),
		ClientIP:   c.ClientIP(),
		ClientTime: time.Now(),
	})
	c.JSON(consts.StatusOK, common.APIResponseSuccess(nil))
}

// GetURL returns the signed URL to access an object.
func GetURL(ctx context.Context, c *app.RequestContext) {
	tenantRequest := &CommonTenantRequest{}
	if err := tenantRequest.FromRequestContext(c); err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, err.Error()))
		return
	}
	tenant, err := authTenant(tenantRequest)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, common.APIResponseError(consts.StatusUnauthorized, err.Error()))
		return
	}
	hlog.CtxInfof(ctx, "[RailgunCDN][Request] Method=%s AppID=%s ObjectPath=%s", "GetURL", tenantRequest.AppID, tenantRequest.ObjectPath)
	if tenantRequest.ObjectPath == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing object path"))
		return
	}
	privateURL, expires, err := getObjectPrivateURL(tenant, tenantRequest)
	if err != nil {
		hlog.CtxErrorf(ctx, "[RailgunCDN][Error] Method=%s AppID=%s Error=%s", "GetURL", tenantRequest.AppID, err.Error())
		c.JSON(consts.StatusInternalServerError, common.APIResponseError(consts.StatusInternalServerError, err.Error()))
		return
	}
	matomo.ReportEvent(ctx, matomo.Event{
		SiteID:     tenant.SiteID,
		ActionName: "railgun-cdn:server:GetURL",
		URL:        config.Conf.Services.RailgunCDN.CDN.Endpoint + tenantRequest.ObjectPath,
		UserAgent:  string(c.UserAgent()),
		ClientIP:   c.ClientIP(),
		ClientTime: time.Now(),
	})
	c.JSON(consts.StatusOK, common.APIResponseSuccess(GetURLResponse{
		URL:     privateURL,
		Expires: expires,
	}))
}

// ClientGateway handles the client access request and redirects it to the actual object URL.
func ClientGateway(ctx context.Context, c *app.RequestContext) {
	appId := c.Query("a")
	objectPath := c.Query("o")
	sign := c.Query("s")
	timestamp := c.Query("t")
	if appId == "" || objectPath == "" || sign == "" || timestamp == "" {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "missing required parameter"))
		return
	}
	var siteId string
	if tenant, ok := config.Conf.Services.RailgunCDN.Tenants[appId]; ok {
		siteId = tenant.SiteID
	} else {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "invalid parameter"))
		return
	}
	timestampParsed, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, common.APIResponseError(consts.StatusBadRequest, "invalid parameter"))
		return
	}
	publicURL := api.GetObjectPublicURL(appId, objectPath, sign, timestampParsed)
	hlog.CtxInfof(ctx, "[RailgunCDN][Request] Method=%s URI=%s", "ClientGateway", objectPath)
	matomo.ReportEvent(ctx, matomo.Event{
		SiteID:     siteId,
		ActionName: "railgun-cdn:client:Gateway",
		URL:        config.Conf.Services.RailgunCDN.CDN.Endpoint + objectPath,
		UserAgent:  string(c.UserAgent()),
		ClientIP:   c.ClientIP(),
		ClientTime: time.Now(),
	})
	c.Redirect(consts.StatusMovedPermanently, []byte(publicURL))
}
