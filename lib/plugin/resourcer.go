package plugin

// Resourcer behaves as a stock information resource
type Resourcer interface {
	GetStocker(string) (Stocker, error)
	GetStockers([]string) (Stockers, error)
}
