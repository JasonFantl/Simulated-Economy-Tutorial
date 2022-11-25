package economy

import "fmt"

var actors map[*Actor]bool

func Initialize(size int) {
	actors = make(map[*Actor]bool)

	for i := 0; i < size; i++ {
		actors[NewActor()] = true
	}
}

func Update() {

	buyers := 0
	sellers := 0
	for i := 0; i < 100; i++ {
		for actor := range actors {
			actor.update()
			if actor.isBuyer() {
				buyers++
			} else if actor.isSeller() {
				sellers++
			}
		}
		if i == 0 {
			fmt.Println(buyers, sellers)
		}
	}

	updateGraph()
}

func Influence(value float64) {
	// do something to influence the economy
	influenceNumPeople := 30
	for actor := range actors {
		actor.basePersonalValue += value
		influenceNumPeople--
		if influenceNumPeople == 0 {
			break
		}
	}
}
