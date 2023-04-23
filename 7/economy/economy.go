package economy

// EconomicAgent is an interface that requires the minimum methods to interact in the economy
type EconomicAgent interface {
	isSelling(Good) (bool, float64)
	transact(Good, bool, float64)
	gossip(Good) float64
}
