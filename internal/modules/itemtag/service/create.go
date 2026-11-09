package service

import "context"

func (s *Service) CreateTag(ctx context.Context, name string, isPublic bool, values []string) error {
	return s.strg.CreateTag(ctx, name, isPublic, values)
}
