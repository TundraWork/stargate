package api

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/tundrawork/stargate/config"
)

// GetObjectPrivateURL gets the private CDN URL of an object.
func GetObjectPrivateURL(objectKey string, ttl int64) (publicURL string, expires int64, err error) {
	return GetObjectURL(config.Conf.Services.RailgunCDN.Private.Endpoint, objectKey, ttl)
}

// GetObjectPublicURL gets the public CDN URL of an object.
func GetObjectPublicURL(objectKey string, ttl int64) (publicURL string, expires int64, err error) {
	return GetObjectURL(config.Conf.Services.RailgunCDN.CDN.Endpoint, objectKey, ttl)
}

// GetObjectURL gets the URL of an object using specific endpoint.
func GetObjectURL(endpoint string, objectKey string, ttl int64) (url string, expires int64, err error) {
	if ttl <= 0 {
		return "", -1, fmt.Errorf("ttl must be a positive integer")
	}
	if len(objectKey) == 0 || objectKey[len(objectKey)-1] == '/' {
		return "", -1, fmt.Errorf("invalid object key")
	}

	expires = time.Now().Unix() + ttl
	timestamp := expires + config.Conf.Services.RailgunCDN.CDN.TimestampOffset
	hashable := fmt.Sprintf("%s%s%d", config.Conf.Services.RailgunCDN.CDN.PKey, objectKey, timestamp)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(hashable)))

	url = fmt.Sprintf("%s%s?sign=%s&t=%d",
		endpoint,
		objectKey,
		sign,
		timestamp,
	)

	return url, expires, nil
}
