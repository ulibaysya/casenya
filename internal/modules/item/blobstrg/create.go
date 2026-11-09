package blobstrg

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *BlobStorage) CreateItemUploadURL(ctx context.Context) (string, string, error) {
	key := path.Join(s.keyFunc(), "img.png")

	req, err := s.pc.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return "", "", nil
	}

	return req.URL, fmt.Sprintf("https://%v/%v/%v", s.host, s.bucket, key), nil
}
