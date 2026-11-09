package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"path"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type presignClient interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type presigner struct {
	bucket string
	c      presignClient
}

func (p *presigner) UploadItemImage(ctx context.Context) (string, error) {
	keyBytes := make([]byte, 32)
	rand.Read(keyBytes)
	req, err := p.c.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &p.bucket,
		Key:    new(path.Join(base64.RawURLEncoding.EncodeToString(keyBytes), "img.png")),
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	c := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	pc := s3.NewPresignClient(c, func(po *s3.PresignOptions) {
		po.Expires = time.Second * 60
	})
	presigner := &presigner{
		bucket: "items",
		c:      pc,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	url, err := presigner.UploadItemImage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
}
