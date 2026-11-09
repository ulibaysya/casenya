package service

import (
	"context"
)

func (s *Service) UploadItemImage(ctx context.Context, id int64) (string, error) {
	uploadURL, actualURL, err := s.blobStrg.CreateItemImageUploadURL(ctx)
	if err != nil {
		return "", err
	}
	if err = s.strg.SaveItemImageURL(ctx, id, actualURL); err != nil {
		return "", err
	}
	return uploadURL, nil
}
