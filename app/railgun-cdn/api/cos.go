package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"github.com/tundrawork/stargate/app/common"
)

type ObjectMetadata struct {
	ContentType   *string `json:"content-type"`
	ContentLength *int64  `json:"content-length"`
	ETag          *string `json:"etag"`
	LastModified  *string `json:"last-modified"`
	CRC64         *string `json:"crc64"`
}

type ObjectKey string

type ListObjectsResponse map[ObjectKey]ObjectMetadata

type PutObjectResponse struct {
	ETag  string `json:"etag"`
	CRC64 string `json:"crc64"`
}

type HeadObjectResponse ObjectMetadata

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

// GetBucket lists objects in a COS bucket.
func GetBucket(ctx context.Context, prefix string) (ListObjectsResponse, error) {
	opt := &cos.BucketGetOptions{
		Prefix: prefix,
	}
	resp, _, err := cosClient.Bucket.Get(ctx, opt)
	if err != nil {
		return ListObjectsResponse{}, err
	}
	if resp == nil {
		return ListObjectsResponse{}, errors.New("empty response from storage")
	}
	res := make(ListObjectsResponse)
	for _, obj := range resp.Contents {
		// trim prefix from key
		obj.Key = obj.Key[len(prefix):]
		// trim quotes from ETag
		obj.ETag = obj.ETag[1 : len(obj.ETag)-1]
		res[ObjectKey(obj.Key)] = ObjectMetadata{
			ContentType:   nil,
			ContentLength: common.ToPtr(obj.Size),
			ETag:          common.ToPtr(obj.ETag),
			LastModified:  common.ToPtr(obj.LastModified),
			CRC64:         nil,
		}
	}
	return res, err
}

// PutObject puts a streamable object to COS.
func PutObject(ctx context.Context, objectKey string, dataStream io.Reader, contentType string, ttl int64) (PutObjectResponse, error) {
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
	resp, err := cosClient.Object.Put(ctx, objectKey, dataStream, opt)
	if err != nil {
		return PutObjectResponse{}, err
	}
	if resp == nil {
		return PutObjectResponse{}, errors.New("empty response from storage")
	}
	res := PutObjectResponse{
		ETag:  resp.Header.Get("ETag"),
		CRC64: resp.Header.Get("x-cos-hash-crc64ecma"),
	}
	return res, err
}

// HeadObject retrieves the metadata of an object from COS.
func HeadObject(ctx context.Context, objectKey string) (HeadObjectResponse, error) {
	resp, err := cosClient.Object.Head(ctx, objectKey, nil)
	if err != nil {
		return HeadObjectResponse{}, err
	}
	if resp == nil {
		return HeadObjectResponse{}, errors.New("empty response from storage")
	}
	// trim quotes from ETag
	eTag := resp.Header.Get("ETag")
	eTag = eTag[1 : len(eTag)-1]
	res := HeadObjectResponse{
		ContentType:   common.ToPtr(resp.Header.Get("Content-Type")),
		ContentLength: common.ToPtr(resp.ContentLength),
		ETag:          common.ToPtr(eTag),
		LastModified:  common.ToPtr(resp.Header.Get("Last-Modified")),
		CRC64:         common.ToPtr(resp.Header.Get("x-cos-hash-crc64ecma")),
	}
	return res, err
}

// DeleteObject deletes an object from COS.
func DeleteObject(ctx context.Context, objectKey string) error {
	_, err := cosClient.Object.Delete(ctx, objectKey)
	return err
}
