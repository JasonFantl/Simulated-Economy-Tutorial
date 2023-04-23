package economy

import (
	"fmt"
	"math/rand"
)

// Merchant tracks lots of information about each city in order to optimally arbitrage
type Merchant struct {
	money               float64
	location            cityName
	buysSells           Good
	carryingCapacity    int
	owned               int
	expectedPrices      map[Good]map[cityName]float64 // merchants use this instead of the value in the market
	expectedMovingCosts map[cityName]map[cityName]float64
}

// NewMerchant creates a merchant
func NewMerchant(location cityName, good Good) *Merchant {
	merchant := &Merchant{
		money:               1000,
		location:            location,
		buysSells:           good,
		carryingCapacity:    20,
		owned:               0,
		expectedPrices:      make(map[Good]map[cityName]float64),
		expectedMovingCosts: make(map[cityName]map[cityName]float64),
	}

	// initialize expected prices, only know starting location to begin
	merchant.expectedPrices[good] = make(map[cityName]float64)
	merchant.expectedPrices[good][location] = 0

	for _, cityName1 := range Locations {
		merchant.expectedMovingCosts[cityName1] = make(map[cityName]float64)
		for _, cityName2 := range Locations {
			merchant.expectedMovingCosts[cityName1][cityName2] = 1
		}
	}
	return merchant
}

func (merchant *Merchant) update(city *City, citiesAgents map[EconomicAgent]bool) {
	// usually people don't try to buy or sell things
	if rand.Float64() > 0.1 {
		return
	}

	// get some gossip
	for local := range citiesAgents {
		for good := range merchant.expectedPrices {
			merchant.expectedPrices[good][city.name] = 0.9*merchant.expectedPrices[good][city.name] + 0.1*local.gossip(good)
			break
		}
	}

	// look to buy
	willingBuyPrice := merchant.expectedPrices[merchant.buysSells][merchant.location]
	bestSellLocation, bestSellProfit := merchant.bestSellProfit(merchant.buysSells)

	if bestSellLocation != merchant.location && merchant.owned < merchant.carryingCapacity { // no possible profit by buying and selling in same location
		// try and find someone to buy from
		for otherAgent := range city.locals {
			isSeller, sellingPrice := otherAgent.isSelling(merchant.buysSells)
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
			otherAgent.transact(merchant.buysSells, false, sellingPrice)
			break
		}
	}

	// TODO: integrate with networked travel routes
	// randomly move cities, but weight by moving costs so high cost cities are less likely to go to
	if rand.Intn(100) == 0 {

		// currently the moving cost is not payed, but it is considered in setting prices
		totalCosts := 0.0
		for _, cost := range merchant.expectedMovingCosts[merchant.location] {
			totalCosts += 1 / cost
		}
		for location, cost := range merchant.expectedMovingCosts[merchant.location] {
			if rand.Float64()*totalCosts < 1/cost {
				merchant.location = location
				break
			} else {
				totalCosts -= 1 / cost
			}
		}
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
		merchant.money -= price
		merchant.owned++
	} else {
		merchant.money += price
		merchant.owned--
	}
}

func (merchant *Merchant) gossip(good Good) float64 {
	return merchant.expectedPrices[good][merchant.location]
}

// find the best location to travel to and how much you would make selling a good there minus the travel expense
func (merchant *Merchant) bestSellProfit(good Good) (cityName, float64) {
	// find the best place to sell the good
	bestLocation := merchant.location
	bestProfit := 0.0
	for location, sellPrice := range merchant.expectedPrices[good] {
		potentialProfit := sellPrice - merchant.expectedMovingCosts[merchant.location][location]
		if potentialProfit > bestProfit {
			bestLocation = location
			bestProfit = potentialProfit
		}
	}

	return bestLocation, bestProfit
}
