package userstorage

import (
	"time"

	"github.com/valkey-io/valkey-go"
)

type UserStorage struct {
	client   valkey.Client
	timeFunc func() time.Time
}

func New(client valkey.Client, timeFunc func() time.Time) *UserStorage {
	return &UserStorage{
		client:   client,
		timeFunc: timeFunc,
	}
}
