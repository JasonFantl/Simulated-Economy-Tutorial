package economy

import (
	"math"
	"math/rand"
)

type EconomicActor interface {
	isSelling(Good) (bool, float64)
	transact(Good, bool, float64)
	gossip(Good) float64
}

type Local struct {
	money    float64
	location Location
	markets  map[Good]*Market
}

func NewLocal(location Location) *Local {
	local := &Local{
		money:    1000,
		location: location,
		markets: map[Good]*Market{
			WOOD:    NewMarket(rand.Intn(20), 4+rand.Float64()*4, 15),
			CHAIR:   NewMarket(rand.Intn(20), 30+rand.Float64()*20, 10),
			LEISURE: NewMarket(0, 2+rand.Float64()*4, 50),
		},
	}

	// set expected prices to match our current value
	for good, market := range local.markets {
		market.expectedMarketPrice = local.currentPersonalValue(good)
	}

	return local
}

func (local *Local) update() {
	// usually people don't try to buy or sell things
	if rand.Float64() > 0.1 {
		return
	}

	// we sometimes break a chair
	if rand.Float64() < 0.01 {
		if local.markets[CHAIR].ownedGoods > 0 {
			local.markets[CHAIR].ownedGoods--
		}
	}

	if float64(iteration)/250.0 > rand.Float64() { // slow start the economy since initial conditions are all over the place
		// evaluate all your actions
		doNothingValue := local.potentialPersonalValue(LEISURE)

		cutWoodValue := math.Max(local.potentialPersonalValue(WOOD), local.priceToValue(local.markets[WOOD].expectedMarketPrice))
		if local.location == RIVERWOOD {
			cutWoodValue *= 2
		}

		buildChairValue := 0.0
		materialCount := 4
		if local.location == SEASIDE {
			materialCount = 2
		}
		if local.markets[WOOD].ownedGoods > materialCount {
			potentialChairValue := math.Max(local.potentialPersonalValue(CHAIR), local.priceToValue(local.markets[CHAIR].expectedMarketPrice))
			materialValue := math.Max(local.currentPersonalValue(WOOD), local.priceToValue(local.markets[WOOD].expectedMarketPrice)) * float64(materialCount)
			buildChairValue = potentialChairValue - materialValue
		}

		// act out the best action
		maxValueAction := math.Max(math.Max(doNothingValue, cutWoodValue), buildChairValue)
		if maxValueAction == doNothingValue {
			local.markets[LEISURE].ownedGoods++ // we value doing nothing less and less the more we do it (diminishing utility)
		} else {
			if maxValueAction == cutWoodValue {
				local.markets[WOOD].ownedGoods++
				if local.location == RIVERWOOD {
					local.markets[WOOD].ownedGoods++
				}
			} else if maxValueAction == buildChairValue {
				local.markets[WOOD].ownedGoods -= materialCount
				local.markets[CHAIR].ownedGoods++
			}
			local.markets[LEISURE].ownedGoods = 0 // make sure we have renewed value for doing nothing since we just did something
		}
	}

	nearbyActors := make(map[EconomicActor]bool)
	for otherActor := range locals {
		if local.location == otherActor.location {
			nearbyActors[otherActor] = true
		}
	}
	for otherActor := range merchants {
		if local.location == otherActor.location {
			nearbyActors[otherActor] = true
		}
	}

	for good := range local.markets {
		if good != LEISURE {
			local.updateMarket(good, nearbyActors)
		}
	}
}

func (local *Local) isSelling(good Good) (bool, float64) {
	if !local.isSeller(good) || local.markets[good].ownedGoods <= 0 {
		return false, 0
	}
	return true, local.markets[good].expectedMarketPrice
}

func (local *Local) transact(good Good, buying bool, price float64) {
	local.markets[good].timeSinceLastTransaction = 0
	if buying {
		local.money -= price
		local.markets[good].ownedGoods++
		local.markets[good].expectedMarketPrice -= local.markets[good].beliefVolatility
	} else {
		local.money += price
		local.markets[good].ownedGoods--
		local.markets[good].expectedMarketPrice += local.markets[good].beliefVolatility
	}
}

func (local *Local) gossip(good Good) float64 {
	return local.markets[good].expectedMarketPrice
}
