package api

import (
	"crypto/md5"
	"fmt"
	"github.com/tundrawork/stargate/config"
	"time"
)

// GetObjectPublicURL gets the public CDN URL of an object.
func GetObjectPublicURL(objectKey string, ttl int64) (publicURL string, expires int64, err error) {
	if ttl <= 0 {
		return "", -1, fmt.Errorf("ttl must be a positive integer")
	}
	if len(objectKey) == 0 || objectKey[len(objectKey)-1] == '/' {
		return "", -1, fmt.Errorf("invalid object key")
	}
	timestamp := time.Now().Unix() + config.Conf.Services.RailgunCDN.CDN.TimestampOffset + ttl
	hashable := fmt.Sprintf("%s%s%d", config.Conf.Services.RailgunCDN.CDN.PKey, objectKey, timestamp)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(hashable)))

	publicURL = fmt.Sprintf("%s%s?sign=%s&t=%d",
		config.Conf.Services.RailgunCDN.CDN.Endpoint,
		objectKey,
		sign,
		timestamp,
	)
	return publicURL, timestamp, nil
}
