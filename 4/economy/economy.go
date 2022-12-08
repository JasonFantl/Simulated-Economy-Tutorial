package economy

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

	// // termites
	// if iteration >= 200 && iteration < 400 {
	// 	for actor := range actors {
	// 		if rand.Float64() < 0.5 && actor.markets[WOOD].ownedGoods > 0 {
	// 			actor.markets[WOOD].ownedGoods--
	// 		}
	// 	}
	// }

	// // new forest
	// if iteration == 600 {
	// 	for actor := range actors {
	// 		actor.markets[WOOD].ownedGoods += 100
	// 	}
	// }

	// // chairs not liked so much anymore
	// if iteration == 200 {
	// 	for actor := range actors {
	// 		if rand.Float64() < 0.5 {
	// 			actor.markets[CHAIR].basePersonalValue *= 0.5
	// 		}
	// 	}
	// }

	// // government regulation
	// if iteration == 400 {
	// 	for actor := range actors {
	// 		actor.markets[CHAIR].ownedGoods = 1
	// 	}
	// }

	// // government reparations
	// if iteration == 600 {
	// 	for actor := range actors {
	// 		actor.markets[CHAIR].ownedGoods = 30
	// 	}
	// }

	updateGraph()
}

func Influence(good Good, value float64) {
	// do something to influence the economy
	// influenceNumPeople := 100
	// for actor := range actors {
	// 	newValue := actor.markets[good].ownedGoods + int(value)
	// 	if newValue >= 0 {
	// 		actor.markets[good].ownedGoods = newValue
	// 	}
	// 	influenceNumPeople--
	// 	if influenceNumPeople == 0 {
	// 		break
	// 	}
	// }

	for actor := range actors {
		actor.money += value * 10
	}
}
