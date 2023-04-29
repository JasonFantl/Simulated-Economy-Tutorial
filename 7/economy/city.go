package economy

import "fmt"

type cityName string

// City separates economies from each other and managers all of its residents
type City struct {
	name      cityName
	locals    map[*Local]bool
	merchants map[*Merchant]bool

	inboundTravelWays  map[cityName]travelWayInbound  // currently only supports one travelWay per city
	outboundTravelWays map[cityName]travelWayOutbound // currently only supports one travelWay per city

	networkPorts *NetworkedTravelWays
}

// NewCity creates a city
func NewCity(name string, size int) *City {
	city := &City{
		name:      cityName(name),
		locals:    make(map[*Local]bool),
		merchants: make(map[*Merchant]bool),

		inboundTravelWays:  make(map[cityName]travelWayInbound),
		outboundTravelWays: make(map[cityName]travelWayOutbound),
	}

	for i := 0; i < size; i++ {
		city.locals[NewLocal()] = true
	}
	for i := 0; i < size/4; i++ {
		city.merchants[NewMerchant(city, CHAIR)] = true
	}

	city.networkPorts = SetupNetworkedTravelWay(55555, city)

	return city
}

// Update will take a time step. All residents will get their own Update method called
func (city *City) Update() {

	// speed up the simulation
	for i := 0; i < 100; i++ {
		// check for new merchants
		for _, travelWay := range city.inboundTravelWays {
			if existNewMerchant, newMerchant := travelWay.receiveImmigrant(); existNewMerchant {
				city.merchants[newMerchant] = true
				newMerchant.Location = city.name

				// if the merchant is rich, tax them and distribute amongst the locals
				if newMerchant.Money > 1000.0 {
					tax := (newMerchant.Money - 1000) / 10
					newMerchant.Money -= tax
					for local := range city.locals {
						local.money += tax / float64(len(city.locals))
					}
				}
			}
		}

		// run all the agents
		for local := range city.locals {
			local.update(city)
		}
		for merchant := range city.merchants {
			merchant.update(city)
		}
	}

	// TMP:: after some time connect one city to another
	if len(previousDataPoints[WOOD][city.name]) == 2 && city.name == cityName("WINTERHOLD") {
		city.networkPorts.RequestConnection("[::]:55555")
	}
	// TMP:: after some time connect one city to another
	if len(previousDataPoints[WOOD][city.name]) == 200 && city.name == cityName("RIVERWOOD") {
		city.networkPorts.RequestConnection("[::]:55557")
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

func (city *City) addEnteringTravelWay(travelWay travelWayInbound) {
	if travelWay.startCity() == city.name {
		fmt.Printf("Cannot create inbound travelWay from a city to itself (%s)\n", city.name)
		return
	}
	city.inboundTravelWays[travelWay.startCity()] = travelWay
}

func (city *City) addLeavingTravelWay(travelWay travelWayOutbound) {
	if travelWay.endCity() == city.name {
		fmt.Printf("Cannot create inbound travelWay from a city to itself (%s)\n", city.name)
		return
	}
	city.outboundTravelWays[travelWay.endCity()] = travelWay
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
