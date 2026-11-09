package service

import "context"

type Storage interface {
	CreateUser(ctx context.Context, name string) (int64, error)
}

type Service struct {
	strg Storage
}

func New(strg Storage) *Service {
	return &Service{
		strg: strg,
	}
}
