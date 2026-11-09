package watchlist_enums

type WatchItemType string

const (
	PRODUCT WatchItemType = "PRODUCT"
	// Future types can be added here
)

func (t WatchItemType) String() string {
	return string(t)
}
