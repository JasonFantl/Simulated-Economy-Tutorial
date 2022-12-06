package economy

import (
	"math"
	"math/rand"
)

type Actor struct {
	money   float64
	markets map[Good]*Market
}

func NewActor() *Actor {
	actor := &Actor{
		money: 1000,
		markets: map[Good]*Market{
			WOOD:    NewMarket(4 + rand.Float64()*4),
			CHAIR:   NewMarket(20 + rand.Float64()*10),
			LEISURE: NewMarket(8 + rand.Float64()*4),
		},
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

	// evaluate all your actions
	doNothingValue := actor.potentialPersonalValue(LEISURE)

	cutWoodValue := math.Max(actor.potentialPersonalValue(WOOD), actor.markets[WOOD].expectedMarketValue)

	buildChairValue := 0.0
	materialCount := 4
	if actor.markets[WOOD].ownedGoods > materialCount {
		potentialChairValue := math.Max(actor.potentialPersonalValue(CHAIR), actor.markets[CHAIR].expectedMarketValue)
		materialValue := math.Max(actor.currentPersonalValue(WOOD), actor.markets[WOOD].expectedMarketValue) * float64(materialCount)
		buildChairValue = potentialChairValue - materialValue
	}

	// act out the best action
	maxValueAction := math.Max(math.Max(doNothingValue, cutWoodValue), buildChairValue)
	if maxValueAction == doNothingValue {
		actor.markets[LEISURE].ownedGoods++ // we value doing nothing less and less the more we do it (diminishing utility)
	} else {
		if maxValueAction == cutWoodValue {
			actor.markets[WOOD].ownedGoods++
		} else if maxValueAction == buildChairValue {
			actor.markets[WOOD].ownedGoods -= materialCount
			actor.markets[CHAIR].ownedGoods++
		}
		actor.markets[LEISURE].ownedGoods = 0 // make sure we have renewed value for doing nothing since we just did something
	}

	for good := range actor.markets {
		if good != LEISURE {
			actor.updateMarket(good)
		}
	}
}
