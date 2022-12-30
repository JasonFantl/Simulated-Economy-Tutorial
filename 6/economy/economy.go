package economy

import (
	"fmt"
	"math/rand"
)

var locals map[*Local]bool
var merchants map[*Merchant]bool

var specialized = false

func Initialize(size int) {
	locals = make(map[*Local]bool)
	merchants = make(map[*Merchant]bool)

	for i := 0; i < size; i++ {
		locals[NewLocal(locations[i%len(locations)])] = true
	}
	for i := 0; i < size/4; i++ {
		merchants[NewMerchant(locations[rand.Intn(len(locations))], WOOD)] = true
	}
}

var iteration = 0

func Update() {

	for i := 0; i < 100; i++ {
		for local := range locals {
			local.update()
		}
		for merchant := range merchants {
			merchant.update()
		}
	}

	iteration++

	// taxing merchants
	taxSums := make(map[Location]float64)
	taxThreshold := 1000.0
	taxPercent := 0.1
	if iteration%10 == 0 {
		for merchant := range merchants {
			if merchant.money > taxThreshold {
				tax := (merchant.money - taxThreshold) * taxPercent
				merchant.money -= tax
				taxSums[merchant.location] += tax
			}
		}
	}
	// now redistribute to the locals equally
	localCount := make(map[Location]int)
	for local := range locals {
		localCount[local.location]++
	}
	for location, taxSum := range taxSums {
		for local := range locals {
			if local.location == location {
				local.money += taxSum / float64(localCount[location])
			}
		}
	}

	// technology advances
	if iteration == 500 {
		specialized = true
	}
	if iteration == 2000 {
		setTravelingCost(RIVERWOOD, SEASIDE, 0.5)
		setTravelingCost(PORTSVILLE, WINTERHOLD, 0.5)
	}
	if iteration == 2500 {
		setTravelingCost(RIVERWOOD, PORTSVILLE, 0.5)
		setTravelingCost(SEASIDE, WINTERHOLD, 0.5)
		setTravelingCost(RIVERWOOD, WINTERHOLD, 1)
		setTravelingCost(SEASIDE, PORTSVILLE, 1)
	}
	if iteration == 3000 {
		setTravelingCost(RIVERWOOD, SEASIDE, 100)
		setTravelingCost(RIVERWOOD, PORTSVILLE, 100)
		setTravelingCost(RIVERWOOD, WINTERHOLD, 100)
		newLocals := make(map[*Local]bool)
		for local := range locals {
			if local.location != RIVERWOOD {
				newLocals[local] = true
			}
		}
		locals = newLocals
	}
	if iteration == 4000 {
		specialized = false
	}

	updateGraph()
}

func Influence(location Location, value float64) {
	// do something to influence the economy

	if location == RIVERWOOD {
		movingCosts[RIVERWOOD][SEASIDE] *= (value+1)/2 + 0.5
		fmt.Printf("cost from Riverwood to Seaside now: %f\n", movingCosts[RIVERWOOD][SEASIDE])
	} else {
		movingCosts[SEASIDE][RIVERWOOD] *= (value+1)/2 + 0.5
		fmt.Printf("cost from Seaside to Riverwood now: %f\n", movingCosts[SEASIDE][RIVERWOOD])
	}
}
