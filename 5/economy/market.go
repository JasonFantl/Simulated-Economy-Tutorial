package economy

import (
	"math"
	"math/rand"
)

type Good string

const (
	WOOD    = "wood"
	CHAIR   = "chair"
	LEISURE = "leisure"
)

type Market struct {
	ownedGoods int

	basePersonalValue   float64
	halfPersonalValueAt float64
	beliefVolatility    float64
	gossipFrequency     float64 // how likely an actor is to gossip with someone each frame. 0.01-0.1 seems to be a good range\

	timeSinceLastTransaction    int
	maxTimeSinceLastTransaction int

	expectedMarketPrices map[Location]float64
}

func NewMarket(owned int, baseValue, halfValueAt float64) *Market {
	market := &Market{
		ownedGoods:                  owned,
		basePersonalValue:           baseValue,
		halfPersonalValueAt:         halfValueAt,
		beliefVolatility:            baseValue / 50,
		timeSinceLastTransaction:    0,
		maxTimeSinceLastTransaction: 10,
		gossipFrequency:             0.01,
		expectedMarketPrices: map[Location]float64{
			RIVERWOOD: (rand.Float64() - 0.5) + baseValue,
			SEASIDE:   (rand.Float64() - 0.5) * baseValue,
		},
	}

	return market
}

func (actor *Actor) updateMarket(good Good) {

	// gossip, hear about other economies as well
	if rand.Float64() < actor.markets[good].gossipFrequency {
		for otherActor := range actors {
			for location, otherExpectedPrice := range otherActor.markets[good].expectedMarketPrices {
				if otherExpectedPrice > actor.markets[good].expectedMarketPrices[location] {
					actor.markets[good].expectedMarketPrices[location] += actor.markets[good].beliefVolatility
				} else if otherExpectedPrice < actor.markets[good].expectedMarketPrices[location] {
					actor.markets[good].expectedMarketPrices[location] -= actor.markets[good].beliefVolatility
				}
			}
			break
		}
	}

	willingBuyPrice := actor.markets[good].expectedMarketPrices[actor.location]

	// only track failed time for when we could transact but didn't
	if actor.isBuyer(good) && actor.money >= willingBuyPrice {
		actor.markets[good].timeSinceLastTransaction++
	} else if actor.isSeller(good) && actor.markets[good].ownedGoods > 0 {
		actor.markets[good].timeSinceLastTransaction++
	}

	// only buyers initiate transactions (usually buyers come to sellers, not the other way around)
	if actor.isBuyer(good) && actor.money >= willingBuyPrice {

		// look for a seller, simulates going from shop to shop
		for otherActor := range actors { // randomly iterates through everyone
			if actor.location != otherActor.location { // must be in the same economy
				continue
			}

			if !otherActor.isSeller(good) || otherActor.markets[good].ownedGoods == 0 { // must be a seller with goods to sell
				continue
			}
			sellingPrice := otherActor.markets[good].expectedMarketPrices[actor.location] // looking at the price tag

			if willingBuyPrice < sellingPrice || actor.money < sellingPrice { // the buyer is unwilling or unable to buy at this price
				continue
			}

			// made it past all the checks, this is someone we can buy from
			actor.money -= sellingPrice
			otherActor.money += sellingPrice
			actor.markets[good].ownedGoods++
			otherActor.markets[good].ownedGoods--
			actor.markets[good].timeSinceLastTransaction, otherActor.markets[good].timeSinceLastTransaction = 0, 0
			actor.markets[good].expectedMarketPrices[actor.location] -= actor.markets[good].beliefVolatility
			otherActor.markets[good].expectedMarketPrices[actor.location] += actor.markets[good].beliefVolatility
			break
		}
	}

	// if we haven't transacted in a while then update expected values
	if actor.markets[good].timeSinceLastTransaction > actor.markets[good].maxTimeSinceLastTransaction {
		actor.markets[good].timeSinceLastTransaction = 0
		if actor.isBuyer(good) {
			// need to be willing to pay more
			actor.markets[good].expectedMarketPrices[actor.location] += actor.markets[good].beliefVolatility
		} else if actor.isSeller(good) {
			// need to be willing to sell for lower
			actor.markets[good].expectedMarketPrices[actor.location] -= actor.markets[good].beliefVolatility
		}
	}
}

// should not be called anywhere except from potentialValue and currentValue
func (actor Actor) personalValue(good Good, x int) float64 {
	S := actor.markets[good].basePersonalValue
	D := actor.markets[good].halfPersonalValueAt
	// simulates diminishing returns
	return S / (math.Pow(float64(x)/D, 3) + 1.0)
}

// returns how much utility you would get from buying another good
func (actor Actor) potentialPersonalValue(good Good) float64 {
	return actor.personalValue(good, actor.markets[good].ownedGoods+1)

}

// how much utility you currently get from your good
func (actor Actor) currentPersonalValue(good Good) float64 {
	return actor.personalValue(good, actor.markets[good].ownedGoods)
}

func (actor Actor) isSeller(good Good) bool {
	return actor.priceToValue(actor.markets[good].expectedMarketPrices[actor.location]) > actor.currentPersonalValue(good)
}

func (actor Actor) isBuyer(good Good) bool {
	return actor.priceToValue(actor.markets[good].expectedMarketPrices[actor.location]) < actor.potentialPersonalValue(good)
}

func (actor *Actor) priceToValue(price float64) float64 {
	return price * actor.utilityPerDollar()
}

func (actor *Actor) valueToPrice(value float64) float64 {
	return value / actor.utilityPerDollar()
}

func (actor *Actor) utilityPerDollar() float64 {
	// utility per dollar has diminishing returns
	return 1000.0 / (actor.money + 1.0)
}
