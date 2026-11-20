package userservice

import (
	"context"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/user/v1"
)

type UserStorage interface {
	SaveUserSession(ctx context.Context, userID int64, session string) error
	FindUserSession(ctx context.Context, session string) (int64, error)
}

type Hasher interface {
	Hash(s string) string
}

type SessionGenerator interface {
	GenerateSession() string
}

type UserService struct {
	ustrg     UserStorage
	hasher    Hasher
	hubClient v1.OpaqueTokenServiceClient
	sg        SessionGenerator
}

func New(ustrg UserStorage, hasher Hasher, hubClient v1.OpaqueTokenServiceClient, sg SessionGenerator) *UserService {
	return &UserService{
		ustrg:     ustrg,
		hasher:    hasher,
		hubClient: hubClient,
		sg:        sg,
	}
}
