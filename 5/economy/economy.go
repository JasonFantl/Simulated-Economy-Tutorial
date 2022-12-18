package economy

import "fmt"

var locals map[*Local]bool
var merchants map[*Merchant]bool

func Initialize(size int) {
	locals = make(map[*Local]bool)
	merchants = make(map[*Merchant]bool)

	for i := 0; i < size; i++ {
		locals[NewLocal()] = true
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

	if iteration == 200 {
		for i := 0; i < 10; i++ {
			merchants[NewMerchant()] = true
		}
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
