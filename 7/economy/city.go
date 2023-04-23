package economy

type cityName string

// Locations is a temporary variable to track the cities (later we will have cities that can appear and disappear during runtime)
var Locations = []cityName{"RIVERWOOD", "SEASIDE", "WINTERHOLD", "PORTSVILLE"}

// City separates economies from each other and managers all of its residents
type City struct {
	name      cityName
	locals    map[*Local]bool
	merchants map[*Merchant]bool
}

// NewCity creates a city
func NewCity(name cityName, size int) *City {
	city := City{
		name:      name,
		locals:    make(map[*Local]bool),
		merchants: make(map[*Merchant]bool),
	}

	for i := 0; i < size; i++ {
		city.locals[NewLocal()] = true
	}
	for i := 0; i < size/4; i++ {
		city.merchants[NewMerchant(name, WOOD)] = true
	}

	return &city
}

// Update will take a time step. All residents will get their own Update method called
func (city *City) Update() {

	allAgents := city.allEconomicAgents()
	for i := 0; i < 100; i++ {
		for local := range city.locals {
			local.update(allAgents)
		}
		for merchant := range city.merchants {
			// TODO: merchants need to move between cities, currently they can't move
			merchant.update(city, allAgents)
		}
	}

	updateGraph(*city)
}

func (city *City) allEconomicAgents() map[EconomicAgent]bool {
	merged := make(map[EconomicAgent]bool)
	for k, v := range city.locals {
		merged[k] = v
	}
	for k, v := range city.merchants {
		merged[k] = v
	}
	return merged
}

// Influence will make some change to the city, hopefully allowing you to run experiments on the economy
func Influence(location cityName, value float64) {
	// do something to influence the economy

	// if location == RIVERWOOD {
	// 	movingCosts[RIVERWOOD][SEASIDE] *= (value+1)/2 + 0.5
	// 	fmt.Printf("cost from Riverwood to Seaside now: %f\n", movingCosts[RIVERWOOD][SEASIDE])
	// } else {
	// 	movingCosts[SEASIDE][RIVERWOOD] *= (value+1)/2 + 0.5
	// 	fmt.Printf("cost from Seaside to Riverwood now: %f\n", movingCosts[SEASIDE][RIVERWOOD])
	// }
}
