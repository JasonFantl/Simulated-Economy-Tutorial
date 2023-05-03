package economy

import (
	"image/color"
)

// EconomicAgent is an interface that requires the minimum methods to interact in the economy
type EconomicAgent interface {
	isSelling(Good) (bool, float64)
	transact(Good, bool, float64)
	gossip(Good) float64
}

type cityName string

// City separates economies from each other and manages all of its residents
type City struct {
	name  cityName
	color color.Color

	locals    map[*Local]bool
	merchants map[*Merchant]bool

	inboundTravelWays  map[cityName]chan *Merchant
	outboundTravelWays map[cityName]chan *Merchant

	networkPorts *networkedTravelWays
}

// NewCity creates a city
func NewCity(name string, col color.Color, size int) *City {
	city := &City{
		name:      cityName(name),
		color:     col,
		locals:    make(map[*Local]bool),
		merchants: make(map[*Merchant]bool),

		inboundTravelWays:  make(map[cityName]chan *Merchant),
		outboundTravelWays: make(map[cityName]chan *Merchant),
	}

	for i := 0; i < size; i++ {
		city.locals[NewLocal()] = true
	}
	for i := 0; i < size/2; i++ {
		city.merchants[NewMerchant(city, CHAIR)] = true
	}

	city.networkPorts = setupNetworkedTravelWay(55555, city)

	return city
}

// Update will take a time step. All residents will get their own Update method called
func (city *City) Update() {

	// speed up the simulation
	for i := 0; i < 100; i++ {
		// check for new merchants
		for _, channel := range city.inboundTravelWays {
			if existNewMerchant, newMerchant := city.receiveImmigrant(channel); existNewMerchant {
				city.merchants[newMerchant] = true
				newMerchant.city = city.name // let the merchant know they arrived

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

	updateGraph(city)
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

func (city *City) receiveImmigrant(channel chan *Merchant) (bool, *Merchant) {
	select { // makes this non-blocking
	case merchant := <-channel:
		return true, merchant
	default:
		return false, nil
	}
}

// CreateTravelWayToCity will make a bidirectional networked connection to another city over which merchants can travel
func (city *City) CreateTravelWayToCity(address string) {
	city.networkPorts.requestConnection(address, BIDIRECTIONAL)
}

// Influence will make some change to the city, hopefully allowing you to run experiments on the economy
func Influence(location cityName, value float64) {
	// do something to influence the economy

}
