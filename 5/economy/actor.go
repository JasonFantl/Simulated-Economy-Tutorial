package economy

import (
	"math"
	"math/rand"
)

type Location string

const (
	RIVERWOOD Location = "Riverwood"
	SEASIDE   Location = "Seaside"
)

var shippingCosts map[Location]map[Location]float64 = map[Location]map[Location]float64{
	RIVERWOOD: {
		RIVERWOOD: 0,
		SEASIDE:   1,
	},
	SEASIDE: {
		RIVERWOOD: 10,
		SEASIDE:   0,
	},
}

type EconomicActor interface {
	isSelling(Good) (bool, float64)
	transact(Good, bool, float64)
	gossip(Good) float64
}

type Actor struct {
	money    float64
	location Location
	markets  map[Good]*Market
}

func NewActor() *Actor {
	actor := &Actor{
		money:    1000,
		location: RIVERWOOD,
		markets: map[Good]*Market{
			WOOD:    NewMarket(rand.Intn(20), 4+rand.Float64()*4, 15),
			CHAIR:   NewMarket(rand.Intn(20), 30+rand.Float64()*20, 10),
			LEISURE: NewMarket(0, 2+rand.Float64()*4, 50),
		},
	}

	// set expected prices to match our current value
	for good, market := range actor.markets {
		market.expectedMarketPrice = actor.currentPersonalValue(good)
	}

	if rand.Float64() < 0.5 {
		actor.location = SEASIDE
	}

	return actor
}

func (actor *Actor) update() {
	// usually people don't try to buy or sell things
	if rand.Float64() > 0.1 {
		return
	}

	// we sometimes break a chair
	if rand.Float64() < 0.01 {
		if actor.markets[CHAIR].ownedGoods > 0 {
			actor.markets[CHAIR].ownedGoods--
		}
	}

	if float64(iteration)/300.0 > rand.Float64() { // slow start the economy since initial conditions are all over the place
		// evaluate all your actions
		doNothingValue := actor.potentialPersonalValue(LEISURE)

		cutWoodValue := math.Max(actor.potentialPersonalValue(WOOD), actor.priceToValue(actor.markets[WOOD].expectedMarketPrice))
		if actor.location == RIVERWOOD {
			cutWoodValue *= 2
		}

		buildChairValue := 0.0
		materialCount := 4
		if actor.location == SEASIDE {
			materialCount = 2
		}
		if actor.markets[WOOD].ownedGoods > materialCount {
			potentialChairValue := math.Max(actor.potentialPersonalValue(CHAIR), actor.priceToValue(actor.markets[CHAIR].expectedMarketPrice))
			materialValue := math.Max(actor.currentPersonalValue(WOOD), actor.priceToValue(actor.markets[WOOD].expectedMarketPrice)) * float64(materialCount)
			buildChairValue = potentialChairValue - materialValue
		}

		// act out the best action
		maxValueAction := math.Max(math.Max(doNothingValue, cutWoodValue), buildChairValue)
		if maxValueAction == doNothingValue {
			actor.markets[LEISURE].ownedGoods++ // we value doing nothing less and less the more we do it (diminishing utility)
		} else {
			if maxValueAction == cutWoodValue {
				actor.markets[WOOD].ownedGoods++
				if actor.location == RIVERWOOD {
					actor.markets[WOOD].ownedGoods++
				}
			} else if maxValueAction == buildChairValue {
				actor.markets[WOOD].ownedGoods -= materialCount
				actor.markets[CHAIR].ownedGoods++
			}
			actor.markets[LEISURE].ownedGoods = 0 // make sure we have renewed value for doing nothing since we just did something
		}
	}

	for good := range actor.markets {
		if good != LEISURE {
			actor.updateMarket(good)
		}
	}
}

func (actor *Actor) isSelling(good Good) (bool, float64) {
	if !actor.isSeller(good) || actor.markets[good].ownedGoods <= 0 {
		return false, 0
	}
	return true, actor.markets[good].expectedMarketPrice
}

func (actor *Actor) transact(good Good, buying bool, price float64) {
	actor.markets[good].timeSinceLastTransaction = 0
	if buying {
		actor.money -= price
		actor.markets[good].ownedGoods++
		actor.markets[good].expectedMarketPrice -= actor.markets[good].beliefVolatility
	} else {
		actor.money += price
		actor.markets[good].ownedGoods--
		actor.markets[good].expectedMarketPrice += actor.markets[good].beliefVolatility
	}
}

func (actor *Actor) gossip(good Good) float64 {
	return actor.markets[good].expectedMarketPrice
}
