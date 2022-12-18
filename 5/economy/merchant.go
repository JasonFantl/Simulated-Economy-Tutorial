package economy

import (
	"math/rand"
)

type Location string

const (
	RIVERWOOD Location = "Riverwood"
	SEASIDE   Location = "Seaside"
)

var movingCosts map[Location]map[Location]float64 = map[Location]map[Location]float64{
	RIVERWOOD: {
		RIVERWOOD: 0,
		SEASIDE:   0.5,
	},
	SEASIDE: {
		RIVERWOOD: 10,
		SEASIDE:   0,
	},
}

type Merchant struct {
	money          float64
	location       Location
	ownedGoods     map[Good]int
	expectedPrices map[Good]map[Location]float64 // merchants use this instead of the value in the market
}

func NewMerchant() *Merchant {
	merchant := &Merchant{
		money:    1000,
		location: RIVERWOOD,
		ownedGoods: map[Good]int{
			WOOD:  0,
			CHAIR: 0,
		},
		expectedPrices: map[Good]map[Location]float64{
			WOOD: {
				RIVERWOOD: 0,
				SEASIDE:   0,
			},
			CHAIR: {
				RIVERWOOD: 0,
				SEASIDE:   0,
			},
		},
	}

	if rand.Float64() < 0.5 {
		merchant.location = SEASIDE
	}

	// immediately get the appropriate expected values
	for good, prices := range merchant.expectedPrices {
		for local := range locals {
			merchant.expectedPrices[good][local.location] = 0.5*prices[local.location] + 0.5*local.markets[good].expectedMarketPrice
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
		for good, prices := range merchant.expectedPrices {
			merchant.expectedPrices[good][local.location] = 0.9*prices[local.location] + 0.1*local.markets[good].expectedMarketPrice
		}
	}

	// look to buy goods
	for good, prices := range merchant.expectedPrices {
		willingBuyPrice := prices[merchant.location]
		bestSellLocation, bestSellProfit := merchant.bestSellProfit(good)

		if bestSellLocation == merchant.location { // no possible profit by buying and selling in same location
			continue
		}

		// try and find someone to buy from
		for otherActor := range locals {
			if merchant.location != otherActor.location { // needs to be in the same city
				continue
			}

			isSeller, sellingPrice := otherActor.isSelling(good)
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
			merchant.transact(good, true, sellingPrice)
			otherActor.transact(good, false, sellingPrice)
			break
		}
	}

	// randomly switch towns, could be a lot smarter about this
	if rand.Intn(100) == 0 {
		// currently the moving cost is not payed, but it is considered in setting prices
		locations := []Location{RIVERWOOD, SEASIDE}
		newLocation := locations[rand.Intn(len(locations))]
		if newLocation != merchant.location {
			merchant.location = newLocation
		}
	}
	// change cities once we bought our goods in bulk
	if merchant.ownedGoods[WOOD] >= 40 && merchant.location == RIVERWOOD {
		merchant.location = SEASIDE
	}
}

func (merchant *Merchant) isSelling(good Good) (bool, float64) {
	if merchant.ownedGoods[good] <= 0 {
		return false, 0
	}

	bestLocation, bestProfit := merchant.bestSellProfit(good)

	// only sell if we are in the best place to sell
	if merchant.location == bestLocation {
		return true, bestProfit
	}
	return false, 0
}

func (merchant *Merchant) transact(good Good, buying bool, price float64) {
	if buying {
		// fmt.Printf("Bought \t%s \t%s \t%f\n", good, merchant.location, price)
		merchant.money -= price
		merchant.ownedGoods[good]++
	} else {
		// fmt.Printf("Sold \t%s \t%s \t%f\n", good, merchant.location, price)
		merchant.money += price
		merchant.ownedGoods[good]--
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
