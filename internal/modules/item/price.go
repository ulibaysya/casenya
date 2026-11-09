package item

type PriceProvider string

const (
	PPCustom PriceProvider = "custom"
)

type Price struct {
	Provider PriceProvider
	Value    int64
}
