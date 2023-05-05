package economy

import (
	"fmt"
	"math/rand"
)

// Merchant tracks lots of information about each city in order to optimally arbitrage
// As annoying as it is, the JSON package needs access to the fields of Merchant, which it can only do if they are public
type Merchant struct {
	Money            float64
	city             cityName
	BuysSells        Good
	CarryingCapacity int
	Owned            int
	ExpectedPrices   map[Good]map[cityName]float64 // merchants use this instead of the value in the market

	bestSellLocation cityName // helpful to track
}

// NewMerchant creates a merchant
func NewMerchant(city *City, good Good) *Merchant {
	merchant := &Merchant{
		Money:            1000,
		city:             city.name,
		BuysSells:        good,
		CarryingCapacity: 20,
		Owned:            0,
		ExpectedPrices:   make(map[Good]map[cityName]float64),
	}

	// initialize expected prices
	for _, good := range goods {
		merchant.ExpectedPrices[good] = make(map[cityName]float64)
		// only know starting cities values, as we find more cities we can update
		merchant.ExpectedPrices[good][city.name] = 0
	}

	return merchant
}

func (merchant *Merchant) update(city *City) {

	// usually people don't try to buy or sell things
	if rand.Float64() > 0.1 {
		return
	}

	// get some gossip
	for local := range city.allEconomicAgents() {
		for good := range merchant.ExpectedPrices { // make sure its random
			merchant.ExpectedPrices[good][city.name] = 0.9*merchant.ExpectedPrices[good][city.name] + 0.1*local.gossip(good)
		}
	}

	// look to buy
	willingBuyPrice := merchant.ExpectedPrices[merchant.BuysSells][merchant.city]
	merchant.bestSellLocation, _ = merchant.bestDeal(merchant.BuysSells, city)
	bestSellPrice := merchant.ExpectedPrices[merchant.BuysSells][merchant.bestSellLocation]

	if merchant.bestSellLocation != merchant.city && merchant.Owned < merchant.CarryingCapacity { // no possible profit by buying and selling in same location
		// try and find someone to buy from
		for otherAgent := range city.locals {
			isSeller, sellingPrice := otherAgent.isSelling(merchant.BuysSells)
			if !isSeller {
				continue
			}

			if willingBuyPrice < sellingPrice || merchant.Money < sellingPrice { // merchant is unwilling or unable to buy at this price
				continue
			}

			if bestSellPrice-sellingPrice <= 0 { // merchant wouldn't make a profit buying this good here
				continue
			}

			// made it past all the checks, this is someone we can buy from
			merchant.transact(merchant.BuysSells, true, sellingPrice)
			otherAgent.transact(merchant.BuysSells, false, sellingPrice)
			break
		}
	}

	// randomly move cities
	if rand.Intn(100) == 0 {
		city.outboundTravelWays.Range(func(_ cityName, outboundTravelWay chan *Merchant) bool {
			merchant.leaveCity(city, outboundTravelWay)
			return false
		})
		return
	}

	// change cities once we bought our good in bulk
	if merchant.Owned >= merchant.CarryingCapacity && merchant.city != merchant.bestSellLocation {
		// make sure we have a path there, if not, randomly move (smarter merchants could potentially do path finding, Q learning?)
		if outboundTravelWay, ok := city.outboundTravelWays.Load(merchant.bestSellLocation); ok {
			merchant.leaveCity(city, outboundTravelWay)
			return
		} else if len(city.outboundTravelWays.channels) > 0 {
			// get random travelWay by using first in iteration
			city.outboundTravelWays.Range(func(_ cityName, randomOutboundTravelWay chan *Merchant) bool {
				merchant.leaveCity(city, randomOutboundTravelWay)
				return false
			})
			return
		}
	}

	// consider switching professions if we aren't selling right now
	if merchant.Owned == 0 {
		bestGood, bestProfit := BED, 0.0
		for good := range merchant.ExpectedPrices {
			_, potentialProfit := merchant.bestDeal(good, city)
			if potentialProfit > bestProfit {
				bestProfit = potentialProfit
				bestGood = good
			}
		}

		merchant.BuysSells = bestGood
	}
}

func (merchant *Merchant) leaveCity(city *City, outboundTravelWay chan *Merchant) {
	// remove self from city
	delete(city.merchants, merchant)
	// enter travelWay
	merchant.city = "traveling..." // not necessary, gets ignored by JSON serializer
	outboundTravelWay <- merchant
}

func (merchant *Merchant) isSelling(good Good) (bool, float64) {
	if good != merchant.BuysSells || merchant.Owned <= 0 {
		return false, 0
	}

	// only sell if we are in the best place to sell, and only sell for the market price
	if merchant.city == merchant.bestSellLocation {
		return true, merchant.ExpectedPrices[good][merchant.bestSellLocation]
	}
	return false, 0
}

func (merchant *Merchant) transact(good Good, buying bool, price float64) {
	if good != merchant.BuysSells {
		fmt.Printf("Merchant somehow transacted a good they don't deal in")
		return
	}
	if buying {
		merchant.Money -= price
		merchant.Owned++
	} else {
		merchant.Money += price
		merchant.Owned--
	}
}

func (merchant *Merchant) gossip(good Good) float64 {
	return merchant.ExpectedPrices[good][merchant.city]
}

// find the best location to travel to and how much you would make selling a good there minus the travel expense.
// returns sell location, expected sell price
func (merchant *Merchant) bestDeal(good Good, city *City) (cityName, float64) {

	// considers nearby (connected) cities. Later, merchants can be more intelligent
	bestSellLocation := merchant.city
	bestProfit := 0.0
	possibleCities := []cityName{merchant.city}
	city.outboundTravelWays.Range(func(cityName cityName, _ chan *Merchant) bool {
		possibleCities = append(possibleCities, cityName)
		return true
	})
	for _, buyLocation := range possibleCities {
		buyPrice := merchant.ExpectedPrices[good][buyLocation]
		for _, sellLocation := range possibleCities {
			sellPrice := merchant.ExpectedPrices[good][sellLocation]
			// get moving costs, pretend it's 1 for now
			movingCost := 1.0

			potentialProfit := sellPrice - (buyPrice + movingCost)
			if potentialProfit > bestProfit {
				bestSellLocation = sellLocation
				bestProfit = potentialProfit
			}
		}
	}

	return bestSellLocation, bestProfit
}
