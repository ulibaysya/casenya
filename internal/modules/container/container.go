package container

import "github.com/ulibaysya/casenya/internal/modules/item"

type Container struct {
	ID    int64
	Price int64
	Items []item.WithWeight
}
