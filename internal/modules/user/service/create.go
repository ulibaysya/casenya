package service

import "context"

func (s *Service) CreateUser(ctx context.Context, name string) (int64, error) {
	return s.strg.CreateUser(ctx, name)
}
