package economy

import "math/rand"

type Actor struct {
	personalValue       float64
	expectedMarketValue float64
}

func NewActor() *Actor {
	actor := &Actor{
		personalValue:       rand.Float64()*9 + 2,
		expectedMarketValue: rand.Float64() * 20,
	}
	return actor
}

var actors map[*Actor]bool

func Initialize(size int) {
	actors = make(map[*Actor]bool)

	for i := 0; i < size; i++ {
		actors[NewActor()] = true
	}
}

func Update() {

	// how quickly we should update our beliefs about the market
	beliefVolatility := 0.1

	// find all buyers and sellers
	sellers := make([]*Actor, 0)
	buyers := make([]*Actor, 0)
	for actor := range actors {
		if actor.expectedMarketValue < actor.personalValue {
			buyers = append(buyers, actor)
		} else {
			sellers = append(sellers, actor)
		}
	}

	// try to buy and sell
	matchedCount := intMin(len(buyers), len(sellers))
	for i := 0; i < matchedCount; i++ {

		// buyers and sellers are randomly matched up
		buyer := buyers[i]
		seller := sellers[i]

		// attempt to transact
		willingSellPrice := seller.expectedMarketValue
		willingBuyPrice := buyer.expectedMarketValue
		if willingBuyPrice >= willingSellPrice {
			// transaction made
			buyer.expectedMarketValue -= beliefVolatility
			seller.expectedMarketValue += beliefVolatility
		} else {
			// transaction failed, make a better offer next time
			buyer.expectedMarketValue += beliefVolatility
			seller.expectedMarketValue -= beliefVolatility
		}
	}

	// if you didn't get matched with anyone, offer a better deal next time
	for i := matchedCount; i < len(buyers); i++ {
		buyers[i].expectedMarketValue += beliefVolatility
	}
	for i := matchedCount; i < len(sellers); i++ {
		sellers[i].expectedMarketValue -= beliefVolatility
	}

	// update graph
	updateGraph()
}

func Influence(value float64) {
	// do something to influence the economy
	influenceNumPeople := 30
	for actor := range actors {
		actor.personalValue += value
		influenceNumPeople--
		if influenceNumPeople == 0 {
			break
		}
	}
}

func intMin(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func intMax(i, j int) int {
	if i > j {
		return i
	}
	return j
}
