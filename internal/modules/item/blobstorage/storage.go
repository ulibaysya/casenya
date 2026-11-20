package blobstorage

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PresignClient interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type BlobStorage struct {
	host    string
	bucket  string
	pc      PresignClient
	keyFunc func() string
}

func New(host string, bucket string, pc PresignClient, keyFunc func() string) *BlobStorage {
	if keyFunc == nil {
		keyFunc = func() string {
			keyBytes := make([]byte, 32)
			rand.Read(keyBytes)
			return base64.RawURLEncoding.EncodeToString(keyBytes)
		}
	}
	return &BlobStorage{
		host:    host,
		bucket:  bucket,
		pc:      pc,
		keyFunc: keyFunc,
	}
}
