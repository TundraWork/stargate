package api

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/tundrawork/stargate/config"
)

// GetObjectPublicURL gets the public CDN URL of an object.
func GetObjectPublicURL(appId, objectPath, sign string, timestamp int64) string {
	return fmt.Sprintf("%s/%s%s?sign=%s&t=%d",
		config.Conf.Services.RailgunCDN.CDN.Endpoint,
		appId,
		objectPath,
		sign,
		timestamp,
	)
}

// SignObject gets the URL of an object using specific endpoint.
func SignObject(objectKey string, ttl int64) (sign string, timestamp int64, expires int64, err error) {
	if ttl <= 0 {
		return "", -1, -1, fmt.Errorf("ttl must be a positive integer")
	}
	if len(objectKey) == 0 || objectKey[len(objectKey)-1] == '/' {
		return "", -1, -1, fmt.Errorf("invalid object key")
	}

	expires = time.Now().Unix() + ttl
	timestamp = expires + config.Conf.Services.RailgunCDN.CDN.TimestampOffset
	hashable := fmt.Sprintf("%s%s%d", config.Conf.Services.RailgunCDN.CDN.PKey, objectKey, timestamp)
	sign = fmt.Sprintf("%x", md5.Sum([]byte(hashable)))

	return sign, timestamp, expires, nil
}
