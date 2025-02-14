package railgun_cdn

import (
	"errors"
	"github.com/tundrawork/stargate/config"
)

// isValidObjectPath checks if the object path is valid.
func isValidObjectPath(objectPath string) bool {
	// The object path must start with a slash and not end with a slash.
	return len(objectPath) > 0 && objectPath[0] == '/' && objectPath[len(objectPath)-1] != '/'
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
