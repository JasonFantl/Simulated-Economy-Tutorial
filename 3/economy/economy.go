package economy

import "math/rand"

var actors map[*Actor]bool

func Initialize(size int) {
	actors = make(map[*Actor]bool)

	for i := 0; i < size; i++ {
		actors[NewActor()] = true
	}
}

var iteration = 0

func Update() {

	for i := 0; i < 100; i++ {
		for actor := range actors {
			actor.update()
		}
	}

	iteration++
	if iteration >= 200 && iteration < 600 {
		for actor := range actors {
			if rand.Float64() < 0.5 && actor.markets[WOOD].ownedGoods > 0 {
				actor.markets[WOOD].ownedGoods--
			}
		}
	}

	if iteration == 600 {
		for actor := range actors {
			actor.markets[WOOD].ownedGoods += 50
		}
	}

	// switch iteration {
	// // case 200:
	// // 	for actor := range actors {
	// // 		actor.markets[CHAIR].ownedGoods = 0
	// // 	}
	// case 500:
	// 	for actor := range actors {
	// 		actor.markets[CHAIR].ownedGoods *= 2
	// 	}
	// case 1000:
	// 	for actor := range actors {
	// 		actor.markets[WOOD].ownedGoods += 50
	// 	}
	// }
	updateGraph()
}

func Influence(good Good, value float64) {
	// do something to influence the economy
	influenceNumPeople := 100
	for actor := range actors {
		newValue := actor.markets[good].ownedGoods + int(value)
		if newValue >= 0 {
			actor.markets[good].ownedGoods = newValue
		}
		influenceNumPeople--
		if influenceNumPeople == 0 {
			break
		}
	}
}
