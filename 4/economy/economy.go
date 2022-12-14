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

	if iteration == 200 {
		for actor := range actors {
			actor.money += 1000
		}
	}
	if iteration == 400 {
		for actor := range actors {
			actor.money -= 1000
		}
	}

	if iteration > 600 {
		for actor := range actors {
			actor.money += 10
		}
	}

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
		newValue := actor.money + value*10
		if newValue >= 0 {
			actor.money = newValue
		}
	}
}
