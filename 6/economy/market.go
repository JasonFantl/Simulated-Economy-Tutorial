package economy

import (
	"math"
	"math/rand"
)

type Good string

const (
	WOOD    Good = "wood"
	CHAIR   Good = "chair"
	THREAD  Good = "thread"
	BED     Good = "bed"
	LEISURE Good = "leisure"
)

type Market struct {
	ownedGoods int

	basePersonalValue   float64
	halfPersonalValueAt float64
	beliefVolatility    float64
	gossipFrequency     float64 // how likely an local is to gossip with someone each frame. 0.01-0.1 seems to be a good range\

	timeSinceLastTransaction    int
	maxTimeSinceLastTransaction int

	expectedMarketPrice float64
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
		expectedMarketPrice:         (rand.Float64() - 0.5) + baseValue,
	}

	return market
}

func (local *Local) updateMarket(good Good, nearbyActors map[EconomicActor]bool) {

	// gossip, hear about other economies as well
	if rand.Float64() < local.markets[good].gossipFrequency {
		for otherActor := range nearbyActors {
			otherExpectedPrice := otherActor.gossip(good)
			if otherExpectedPrice > local.markets[good].expectedMarketPrice {
				local.markets[good].expectedMarketPrice += local.markets[good].beliefVolatility
			} else if otherExpectedPrice < local.markets[good].expectedMarketPrice {
				local.markets[good].expectedMarketPrice -= local.markets[good].beliefVolatility
			}
			break
		}
	}
	willingBuyPrice := local.markets[good].expectedMarketPrice

	// only track failed time for when we could transact but didn't
	if local.isBuyer(good) && local.money >= willingBuyPrice {
		local.markets[good].timeSinceLastTransaction++
	} else if local.isSeller(good) && local.markets[good].ownedGoods > 0 {
		local.markets[good].timeSinceLastTransaction++
	}

	// only buyers initiate transactions (usually buyers come to sellers, not the other way around)
	if local.isBuyer(good) && local.money >= willingBuyPrice {

		// look for a seller, simulates going from shop to shop
		for otherActor := range nearbyActors { // randomly iterates through everyone

			isSeller, sellingPrice := otherActor.isSelling(good)
			if !isSeller {
				continue
			}
			if willingBuyPrice < sellingPrice || local.money < sellingPrice { // the buyer is unwilling or unable to buy at this price
				continue
			}

			// made it past all the checks, this is someone we can buy from
			local.transact(good, true, sellingPrice)
			otherActor.transact(good, false, sellingPrice)
			break
		}
	}

	// if we haven't transacted in a while then update expected values
	if local.markets[good].timeSinceLastTransaction > local.markets[good].maxTimeSinceLastTransaction {
		local.markets[good].timeSinceLastTransaction = 0
		if local.isBuyer(good) {
			// need to be willing to pay more
			local.markets[good].expectedMarketPrice += local.markets[good].beliefVolatility
		} else if local.isSeller(good) {
			// need to be willing to sell for lower
			local.markets[good].expectedMarketPrice -= local.markets[good].beliefVolatility
		}
	}
}

// should not be called anywhere except from potentialValue and currentValue
func (local Local) personalValue(good Good, x int) float64 {
	S := local.markets[good].basePersonalValue
	D := local.markets[good].halfPersonalValueAt
	// simulates diminishing returns
	return S / (math.Pow(float64(x)/D, 3) + 1.0)
}

// returns how much utility you would get from buying another good
func (local Local) potentialPersonalValue(good Good) float64 {
	return local.personalValue(good, local.markets[good].ownedGoods+1)
}

// how much utility you currently get from your good
func (local Local) currentPersonalValue(good Good) float64 {
	return local.personalValue(good, local.markets[good].ownedGoods)
}

func (local Local) isSeller(good Good) bool {
	return local.priceToValue(local.markets[good].expectedMarketPrice) > local.currentPersonalValue(good)
}

func (local Local) isBuyer(good Good) bool {
	return local.priceToValue(local.markets[good].expectedMarketPrice) < local.potentialPersonalValue(good)
}

func (local *Local) priceToValue(price float64) float64 {
	return price * local.utilityPerDollar()
}

func (local *Local) valueToPrice(value float64) float64 {
	return value / local.utilityPerDollar()
}

func (local *Local) utilityPerDollar() float64 {
	// utility per dollar has diminishing returns
	return 1000.0 / (local.money + 1.0)
}
