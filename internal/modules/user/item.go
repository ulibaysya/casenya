package user

type Item struct {
	ID     int64
	ItemID int64
	Price  int64
	Status string
}

type ItemStatus string

var (
	ItemStatusActive       ItemStatus = "active"
	ItemStatusSold         ItemStatus = "sold"
	ItemStatusInWithdrawal ItemStatus = "in_withdrawal"
	ItemStatusWithdrawn    ItemStatus = "withdrawn"
)
