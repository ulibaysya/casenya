package item

import "github.com/ulibaysya/casenya/internal/modules/itemtag"

type Item struct {
	ID       int64
	Name     string
	ImageURL string
	Tags     []itemtag.Tag
	Prices   []Price
}

type PublicItem struct {
	ID       int64
	Name     string
	ImageURL string
	Tags     []Tag
	Price    int64
}

type Tag struct {
	Key    string
	Values []string
}
