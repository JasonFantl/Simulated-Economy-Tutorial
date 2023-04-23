package economy

import (
	"math"
	"math/rand"
)

// Local tracks each market to buy and sell what they need
type Local struct {
	money   float64
	markets map[Good]*Market
}

// NewLocal creates a new local
func NewLocal() *Local {
	local := &Local{
		money: 1000,
		markets: map[Good]*Market{
			WOOD:    NewMarket(rand.Intn(20), 4+rand.Float64()*4, 15),
			CHAIR:   NewMarket(rand.Intn(10), 30+rand.Float64()*20, 5),
			THREAD:  NewMarket(rand.Intn(30), 2+rand.Float64()*2, 50),
			BED:     NewMarket(rand.Intn(2), 50+rand.Float64()*10, 2),
			LEISURE: NewMarket(0, 2+rand.Float64()*4, 50),
		},
	}

	// set expected prices to match our current value
	for good, market := range local.markets {
		market.expectedMarketPrice = local.currentPersonalValue(good)
	}

	return local
}

func (local *Local) update(citiesAgents map[EconomicAgent]bool) {
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

	// we sometimes break a bed
	if rand.Float64() < 0.01 {
		if local.markets[BED].ownedGoods > 0 {
			local.markets[BED].ownedGoods--
		}
	}

	// evaluate all your actions

	doNothingValue := local.potentialPersonalValue(LEISURE)

	cutWoodValue := math.Max(local.potentialPersonalValue(WOOD), local.priceToValue(local.markets[WOOD].expectedMarketPrice))
	spinThreadValue := math.Max(local.potentialPersonalValue(THREAD), local.priceToValue(local.markets[THREAD].expectedMarketPrice))

	buildChairValue := 0.0
	materialCount := 4
	if local.markets[WOOD].ownedGoods > materialCount {
		potentialChairValue := math.Max(local.potentialPersonalValue(CHAIR), local.priceToValue(local.markets[CHAIR].expectedMarketPrice))
		materialValue := math.Max(local.currentPersonalValue(WOOD), local.priceToValue(local.markets[WOOD].expectedMarketPrice)) * float64(materialCount)
		buildChairValue = potentialChairValue - materialValue
	}

	buildBedValue := 0.0
	materialWoodCount := 2
	materialThreadCount := 10
	if local.markets[WOOD].ownedGoods > materialWoodCount && local.markets[THREAD].ownedGoods > materialThreadCount {
		potentialBedValue := math.Max(local.potentialPersonalValue(BED), local.priceToValue(local.markets[BED].expectedMarketPrice))
		materialValue := math.Max(local.currentPersonalValue(WOOD), local.priceToValue(local.markets[WOOD].expectedMarketPrice))*float64(materialWoodCount) +
			math.Max(local.currentPersonalValue(THREAD), local.priceToValue(local.markets[THREAD].expectedMarketPrice))*float64(materialThreadCount)
		buildBedValue = potentialBedValue - materialValue
	}

	// act out the best action
	maxValueAction := math.Max(math.Max(math.Max(math.Max(doNothingValue, cutWoodValue), spinThreadValue), buildChairValue), buildBedValue)
	if maxValueAction == doNothingValue {
		local.markets[LEISURE].ownedGoods++ // we value doing nothing less and less the more we do it (diminishing utility)
	} else {
		if maxValueAction == cutWoodValue {
			local.markets[WOOD].ownedGoods++
		} else if maxValueAction == spinThreadValue {
			local.markets[THREAD].ownedGoods++
		} else if maxValueAction == buildChairValue {
			local.markets[WOOD].ownedGoods -= materialCount
			local.markets[CHAIR].ownedGoods++
		} else if maxValueAction == buildBedValue {
			local.markets[WOOD].ownedGoods -= materialWoodCount
			local.markets[THREAD].ownedGoods -= materialThreadCount
			local.markets[BED].ownedGoods++
		}
		local.markets[LEISURE].ownedGoods = 0 // make sure we have renewed value for doing nothing since we just did something
	}

	for good := range local.markets {
		if good != LEISURE {
			local.updateMarket(good, citiesAgents)
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
