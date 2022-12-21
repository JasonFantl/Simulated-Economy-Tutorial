package economy

import "fmt"

var locals map[*Local]bool
var merchants map[*Merchant]bool

func Initialize(size int) {
	locals = make(map[*Local]bool)
	merchants = make(map[*Merchant]bool)

	for i := 0; i < size; i++ {
		locations := []Location{RIVERWOOD, SEASIDE, WINTERHOLD, PORTSVILLE}
		locals[NewLocal(locations[i%len(locations)])] = true
	}
	// for i := 0; i < size/4; i++ {
	// 	merchants[NewMerchant(RIVERWOOD, WOOD)] = true
	// }
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

	if iteration == 500 {
		movingCosts[RIVERWOOD][SEASIDE] = 1
		movingCosts[SEASIDE][RIVERWOOD] = 1

	}
	if iteration == 1000 {
		movingCosts[SEASIDE][WINTERHOLD] = 1
		movingCosts[WINTERHOLD][SEASIDE] = 1
	}
	if iteration == 1500 {
		movingCosts[WINTERHOLD][PORTSVILLE] = 1
		movingCosts[PORTSVILLE][WINTERHOLD] = 1
	}
	if iteration == 2000 {
		movingCosts[PORTSVILLE][RIVERWOOD] = 1
		movingCosts[RIVERWOOD][PORTSVILLE] = 1
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
