package economy

import (
	"fmt"
	"math/rand"
)

type Location string

const (
	RIVERWOOD  Location = "Riverwood"
	SEASIDE    Location = "Seaside"
	WINTERHOLD Location = "Winterhold"
	PORTSVILLE Location = "Portsville"
)

var movingCosts map[Location]map[Location]float64 = map[Location]map[Location]float64{
	RIVERWOOD: {
		RIVERWOOD:  0,
		SEASIDE:    100,
		WINTERHOLD: 100,
		PORTSVILLE: 100,
	},
	SEASIDE: {
		RIVERWOOD:  100,
		SEASIDE:    0,
		WINTERHOLD: 100,
		PORTSVILLE: 100,
	},
	WINTERHOLD: {
		RIVERWOOD:  100,
		SEASIDE:    100,
		WINTERHOLD: 0,
		PORTSVILLE: 100,
	},
	PORTSVILLE: {
		RIVERWOOD:  100,
		SEASIDE:    100,
		WINTERHOLD: 100,
		PORTSVILLE: 0,
	},
}

type Merchant struct {
	money            float64
	location         Location
	buysSells        Good
	carryingCapacity int
	owned            int
	expectedPrices   map[Good]map[Location]float64 // merchants use this instead of the value in the market
}

func NewMerchant(location Location, good Good) *Merchant {
	merchant := &Merchant{
		money:            1000,
		location:         location,
		buysSells:        good,
		carryingCapacity: 20,
		owned:            0,
		expectedPrices: map[Good]map[Location]float64{
			WOOD: {
				RIVERWOOD: 0,
				SEASIDE:   0,
			},
			CHAIR: {
				RIVERWOOD: 0,
				SEASIDE:   0,
			},
			THREAD: {
				RIVERWOOD: 0,
				SEASIDE:   0,
			},
			BED: {
				RIVERWOOD: 0,
				SEASIDE:   0,
			},
		},
	}

	// immediately get the appropriate expected values
	for local := range locals {
		for good := range merchant.expectedPrices {
			merchant.expectedPrices[good][local.location] = 0.5*merchant.expectedPrices[good][local.location] + 0.5*local.markets[good].expectedMarketPrice
		}
	}

	return merchant
}

func (merchant *Merchant) update() {
	// usually people don't try to buy or sell things
	if rand.Float64() > 0.1 {
		return
	}

	// get some gossip
	for local := range locals {
		for good := range merchant.expectedPrices {
			merchant.expectedPrices[good][local.location] = 0.9*merchant.expectedPrices[good][local.location] + 0.1*local.markets[good].expectedMarketPrice
		}
	}

	// look to buy
	willingBuyPrice := merchant.expectedPrices[merchant.buysSells][merchant.location]
	bestSellLocation, bestSellProfit := merchant.bestSellProfit(merchant.buysSells)

	if bestSellLocation != merchant.location && merchant.owned < merchant.carryingCapacity { // no possible profit by buying and selling in same location
		// try and find someone to buy from
		for otherActor := range locals {
			if merchant.location != otherActor.location { // needs to be in the same city
				continue
			}

			isSeller, sellingPrice := otherActor.isSelling(merchant.buysSells)
			if !isSeller {
				continue
			}

			if willingBuyPrice < sellingPrice || merchant.money < sellingPrice { // merchant is unwilling or unable to buy at this price
				continue
			}

			if bestSellProfit-sellingPrice <= 0 { // merchant wouldn't make a profit buying this good here
				continue
			}

			// made it past all the checks, this is someone we can buy from
			merchant.transact(merchant.buysSells, true, sellingPrice)
			otherActor.transact(merchant.buysSells, false, sellingPrice)
			break
		}
	}

	// randomly switch towns, could be a lot smarter about this
	if rand.Intn(100) == 0 {
		// currently the moving cost is not payed, but it is considered in setting prices
		locations := []Location{RIVERWOOD, SEASIDE}
		merchant.location = locations[rand.Intn(len(locations))]
	}

	// change cities once we bought our good in bulk
	if merchant.owned >= merchant.carryingCapacity && merchant.location != bestSellLocation {
		merchant.location = bestSellLocation
	}

	// consider switching professions if we aren't selling right now
	if merchant.owned == 0 {
		bestGood, bestProfit := WOOD, 0.0
		for good, expectedPrices := range merchant.expectedPrices {
			_, sellProfit := merchant.bestSellProfit(good)
			potentialProfit := sellProfit - expectedPrices[merchant.location] // profit = sell price - buy price
			if potentialProfit > bestProfit {
				bestProfit = potentialProfit
				bestGood = good
			}
		}
		merchant.buysSells = bestGood
	}
}

func (merchant *Merchant) isSelling(good Good) (bool, float64) {
	if good != merchant.buysSells || merchant.owned <= 0 {
		return false, 0
	}

	bestLocation, bestProfit := merchant.bestSellProfit(merchant.buysSells)

	// only sell if we are in the best place to sell
	if merchant.location == bestLocation {
		return true, bestProfit
	}
	return false, 0
}

func (merchant *Merchant) transact(good Good, buying bool, price float64) {
	if good != merchant.buysSells {
		fmt.Printf("Merchant somehow transacted a good they don't deal in")
		return
	}
	if buying {
		// fmt.Printf("Bought \t%s \t%s \t%f\n", good, merchant.location, price)
		merchant.money -= price
		merchant.owned++
	} else {
		// fmt.Printf("Sold \t%s \t%s \t%f\n", good, merchant.location, price)
		merchant.money += price
		merchant.owned--
	}
}

func (merchant *Merchant) gossip(good Good) float64 {
	return merchant.expectedPrices[good][merchant.location]
}

// find the best location to travel to and how much you would make selling a good there minus the travel expense
func (merchant *Merchant) bestSellProfit(good Good) (Location, float64) {
	// find the best place to sell the good
	bestLocation := merchant.location
	bestProfit := 0.0
	for location, sellPrice := range merchant.expectedPrices[good] {
		potentialProfit := sellPrice - movingCosts[merchant.location][location]
		if potentialProfit > bestProfit {
			bestLocation = location
			bestProfit = potentialProfit
		}
	}

	return bestLocation, bestProfit
}
