package railgun_cdn

import (
	"errors"
	"fmt"

	"github.com/tundrawork/stargate/app/railgun-cdn/api"
	"github.com/tundrawork/stargate/config"
)

type TenantBusinessData struct {
	AppID    string
	RootPath string
	SiteID   string
}

// isValidObjectPath checks if the object path is valid.
func isValidObjectPath(objectPath string) bool {
	// The object path must start with a slash and not end with a slash.
	return len(objectPath) > 0 && objectPath[0] == '/' && objectPath[len(objectPath)-1] != '/'
}

// authTenant authenticates the tenant from the common tenant request and returns the tenant's root path.
func authTenant(req *CommonTenantRequest) (*TenantBusinessData, error) {
	if tenant, ok := config.Conf.Services.RailgunCDN.Tenants[req.AppID]; ok {
		if tenant.AppKey == req.AppKey {
			return &TenantBusinessData{
				AppID:    req.AppID,
				RootPath: tenant.RootPath,
				SiteID:   tenant.SiteID,
			}, nil
		}
	}
	return nil, errors.New("tenant authorization failed")
}

// getObjectPrivateURL gets the private CDN URL of an object.
func getObjectPrivateURL(tenant *TenantBusinessData, tenantRequest *CommonTenantRequest) (privateURL string, expires int64, err error) {
	objectKey := "/" + tenant.RootPath + tenantRequest.ObjectPath
	sign, timestamp, expires, err := api.SignObject(objectKey, tenantRequest.TTL)
	if err != nil {
		return "", -1, err
	}
	privateURL = fmt.Sprintf(
		"%s?a=%so=%s&s=%s&t=%d",
		config.Conf.Services.RailgunCDN.Private.Endpoint,
		tenant.AppID,
		tenantRequest.ObjectPath,
		sign,
		timestamp,
	)
	return privateURL, expires, nil
}
