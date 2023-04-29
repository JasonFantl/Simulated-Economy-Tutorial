package main

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jasonfantl/SimulatedEconomy7/economy"
)

// Game is required by ebiten
type Game struct {
}

var previousTime time.Time

var cities []economy.City

// Update will be called at 60 FPS
func (g *Game) Update() error {

	now := time.Now()
	elapsed := now.Sub(previousTime)
	if elapsed.Milliseconds() > 10 {
		for _, city := range cities {
			city.Update()
		}
		previousTime = now
	}

	for _, p := range inpututil.AppendPressedKeys(make([]ebiten.Key, 1)) {
		if p == ebiten.KeyEscape {
			return errors.New("user quit")
		}
	}

	return nil
}

// Draw is called after Update to display to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	economy.GraphExpectedValues(screen, "Price of Wood", economy.WOOD, 100, 300, 0.5, 40.0, 800, 100, 1)
	economy.GraphExpectedValues(screen, "Price of Chairs", economy.CHAIR, 100, 600, 0.5, 8.0, 800, 100, 5)
	economy.GraphExpectedValues(screen, "Price of Thread", economy.THREAD, 600, 300, 0.5, 30.0, 800, 100, 1)
	economy.GraphExpectedValues(screen, "Price of Bed", economy.BED, 600, 600, 0.5, 4.0, 800, 100, 10)

	// economy.GraphGoodsVMoney(screen, "Wood V Money", economy.WOOD, 600, 200, 0.1, 4.0, 250, 10)
	// economy.GraphGoodsVMoney(screen, "Chair V Money", economy.CHAIR, 600, 400, 0.1, 4.0, 250, 10)
	// economy.GraphGoodsVMoney(screen, "Thread V Money", economy.THREAD, 600, 600, 0.1, 2.0, 250, 10)
	// economy.GraphGoodsVMoney(screen, "Bed V Money", economy.BED, 600, 800, 0.1, 20.0, 250, 1)

	for i, city := range cities {
		economy.GraphLeisureVWealth(screen, city, "Leisure V Wealth", 1000, 800, 0.1, 10, 250, 2)

		economy.GraphMerchantType(screen, city, "Merchant types", 200*(float64(i)+1), 800, 50, 5)
	}
}

// Layout determins the window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	rand.Seed(time.Now().Unix())
	game := &Game{}

	ebiten.SetWindowSize(1440, 940)
	ebiten.SetWindowTitle("Economy Simulation")

	arg := os.Args[1]

	cities = make([]economy.City, 2)
	if arg == "1" {
		for i, name := range []string{"RIVERWOOD", "SEASIDE"} {
			cities[i] = *economy.NewCity(name, 20)
		}
	} else {
		for i, name := range []string{"WINTERHOLD", "PORTSVILLE"} {
			cities[i] = *economy.NewCity(name, 20)
		}
	}

	// add connections between cities
	// economy.RegisterChanneledTravelWay(&cities[0], &cities[1])
	// economy.RegisterChanneledTravelWay(&cities[1], &cities[0])

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
