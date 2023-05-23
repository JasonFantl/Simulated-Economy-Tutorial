package main

import (
	"errors"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jasonfantl/SimulatedEconomy8/economy"
	"golang.org/x/exp/maps"
)

// Game is required by ebiten
type Game struct {
}

var previousTime time.Time

var cities []*economy.City

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
		} else if p == ebiten.KeyAlt && inpututil.IsKeyJustPressed(p) {
			cities[0].CreateTravelWayToCity("127.0.0.1:55555")
		}
	}

	return nil
}

// Draw is called after Update to display to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	economy.GraphExpectedValues(screen, "Price of Wood", economy.WOOD, 100, 200, 0.2, 20.0, 800, 200, 1)
	economy.GraphExpectedValues(screen, "Price of Chairs", economy.CHAIR, 100, 400, 0.2, 4.0, 800, 200, 5)
	economy.GraphExpectedValues(screen, "Price of Fur", economy.FUR, 350, 200, 0.2, 20.0, 800, 200, 1)
	economy.GraphExpectedValues(screen, "Price of Bed", economy.BED, 350, 400, 0.2, 2.0, 800, 200, 10)

	economy.GraphMerchantType(screen, cities, "Merchant types", 80, 600, 40, 5)
	for _, city := range cities {
		economy.GraphLeisureVWealth(screen, *city, "Leisure V Wealth", 300, 600, 0.1, 10, 250, 2)
	}
}

// Layout determins the window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

var locationColors = map[string]color.Color{
	"RIVERWOOD":  color.RGBA{58, 158, 33, 100},
	"SEASIDE":    color.RGBA{10, 159, 227, 100},
	"WINTERHOLD": color.RGBA{255, 255, 255, 100},
	"PORTSVILLE": color.RGBA{128, 0, 0, 100},
}

func main() {
	rand.Seed(time.Now().Unix())
	game := &Game{}

	ebiten.SetWindowSize(650, 750)
	ebiten.SetWindowTitle("Economy Simulation")

	cityNames := os.Args[1:]

	cities = make([]*economy.City, len(cityNames))
	for i, name := range cityNames {
		name = strings.ToUpper(name)
		if col, ok := locationColors[name]; ok {
			cities[i] = economy.NewCity(name, col, 20)
		} else {
			fmt.Println("Currently only support cities: " + strings.Join(maps.Keys(locationColors), ", "))
		}
	}

	// add connections between cities
	economy.RegisterTravelWay(cities[0], cities[1])
	economy.RegisterTravelWay(cities[1], cities[0])

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
