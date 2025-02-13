package api

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	cosClient *cos.Client
)

// InitCosClient initializes the COS client.
func InitCosClient(bucket, region, secretID, secretKey string) {
	bucketURLRaw := "https://" + bucket + ".cos." + region + ".myqcloud.com"
	bucketURL, err := url.Parse(bucketURLRaw)
	if err != nil {
		return
	}
	baseURL := &cos.BaseURL{BucketURL: bucketURL}
	cosClient = cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})
}

// PutObject puts a streamable object to COS.
func PutObject(ctx context.Context, objectKey string, dataStream io.Reader, contentType string, ttl int64) error {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	headerOptions := &cos.ObjectPutHeaderOptions{
		ContentType: contentType,
	}
	if ttl > 0 {
		timestamp := time.Now().Unix() + ttl
		headerOptions.Expires = time.Unix(timestamp, 0).Format(time.RFC1123)
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: headerOptions,
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "private", // "private" | "public-read" | "public-read-write" | "authenticated-read"
		},
	}
	_, err := cosClient.Object.Put(ctx, objectKey, dataStream, opt)
	return err
}

// DeleteObject deletes an object from COS.
func DeleteObject(ctx context.Context, objectKey string) error {
	_, err := cosClient.Object.Delete(ctx, objectKey)
	return err
}
