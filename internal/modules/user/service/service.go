package service

import "context"

type Storage interface {
	CreateUser(ctx context.Context, name string) (int64, error)
	SaveUserOpaqueToken(ctx context.Context, id int64, token string) error
	ExchangeUserOpaqueToken(ctx context.Context, token string) (int64, error)
}

type Hasher interface {
	Hash([]byte) string
}

type Service struct {
	strg   Storage
	hasher Hasher
}

func New(strg Storage, hasher Hasher) *Service {
	return &Service{
		strg:   strg,
		hasher: hasher,
	}
}
