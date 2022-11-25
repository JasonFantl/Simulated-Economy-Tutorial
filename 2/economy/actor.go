package economy

import (
	"fmt"
	"math"
	"math/rand"
)

type Actor struct {
	ownedGoods int
	money      float64

	basePersonalValue   float64
	halfPersonalValueAt float64
	expectedMarketValue float64
	beliefVolatility    float64

	timeSinceLastTransaction    int
	maxTimeSinceLastTransaction int
}

func NewActor() *Actor {
	actor := &Actor{
		ownedGoods:                  10,
		money:                       100,
		basePersonalValue:           10 + rand.Float64()*5,
		halfPersonalValueAt:         15,
		expectedMarketValue:         rand.Float64() * 20,
		beliefVolatility:            0.1,
		timeSinceLastTransaction:    0,
		maxTimeSinceLastTransaction: 10,
	}
	return actor
}

func (actor *Actor) update() {
	// usually people don't try to buy or sell things
	if rand.Float64() > 0.1 {
		return
	}

	// gossip
	for otherActor := range actors {
		gossipPrice := otherActor.expectedMarketValue
		if gossipPrice > actor.expectedMarketValue {
			actor.expectedMarketValue += actor.beliefVolatility
		} else if gossipPrice < actor.expectedMarketValue {
			actor.expectedMarketValue -= actor.beliefVolatility
		}
		break
	}

	// only track failed time for when we could transact but didn't
	if actor.isBuyer() && actor.money >= actor.expectedMarketValue {
		actor.timeSinceLastTransaction++
	} else if actor.isSeller() && actor.ownedGoods > 0 {
		actor.timeSinceLastTransaction++
	}

	// only buyers initiate transactions (usually buyers come to sellers, not the other way around)
	if actor.isBuyer() && actor.money >= actor.expectedMarketValue {
		willingBuyPrice := actor.expectedMarketValue

		// look for a seller, simulates going from shop to shop
		for otherActor := range actors { // randomly iterates through everyone
			if !otherActor.isSeller() || otherActor.ownedGoods == 0 { // must be a seller with goods to sell
				continue
			}
			sellingPrice := otherActor.expectedMarketValue // looking at the price tag

			if willingBuyPrice < sellingPrice || actor.money < sellingPrice { // the buyer is unwilling or unable to buy at this price
				continue
			}

			// made it past all the checks, this is someone we can buy from
			actor.money -= sellingPrice
			otherActor.money += sellingPrice
			actor.ownedGoods++
			otherActor.ownedGoods--
			actor.timeSinceLastTransaction, otherActor.timeSinceLastTransaction = 0, 0
			actor.expectedMarketValue -= actor.beliefVolatility
			otherActor.expectedMarketValue += actor.beliefVolatility
			fmt.Println(sellingPrice)
			break
		}
	}

	// if we haven't transacted in a while then update expected values
	if actor.timeSinceLastTransaction > actor.maxTimeSinceLastTransaction {
		actor.timeSinceLastTransaction = 0
		if actor.isBuyer() {
			// need to be willing to pay more
			actor.expectedMarketValue += actor.beliefVolatility
		} else if actor.isSeller() {
			// need to be willing to sell for lower
			actor.expectedMarketValue -= actor.beliefVolatility
		}
	}
}

// should not be called anywhere except from potentialValue and currentValue
func (actor Actor) personalValue(x int) float64 {
	S := actor.basePersonalValue
	D := actor.halfPersonalValueAt
	// simulates diminishing returns
	return S / (math.Pow(float64(x)/D, 3) + 1.0)
}

// returns how much utility you would get from buying another good
func (actor Actor) potentialValue() float64 {
	return actor.personalValue(actor.ownedGoods + 1)

}

// how much utility you currently get from your good
func (actor Actor) currentValue() float64 {
	return actor.personalValue(actor.ownedGoods)
}

func (actor Actor) isSeller() bool {
	return actor.expectedMarketValue > actor.currentValue()
}

func (actor Actor) isBuyer() bool {
	return actor.expectedMarketValue < actor.potentialValue()
}
