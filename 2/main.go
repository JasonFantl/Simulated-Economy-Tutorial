package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jasonfantl/SimulatedEconomy2/economy"
)

type Game struct {
}

var previousTime time.Time

func (g *Game) Update() error {

	now := time.Now()
	elapsed := now.Sub(previousTime)
	if elapsed.Milliseconds() > 10 {
		economy.Update()
		previousTime = now
	}

	for _, p := range inpututil.AppendPressedKeys(make([]ebiten.Key, 1)) {
		if p == ebiten.KeyEscape {
			return errors.New("user quit")
		} else if p == ebiten.KeyArrowUp {
			economy.Influence(1)
		} else if p == ebiten.KeyArrowDown {
			economy.Influence(-1)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	economy.GraphExpectedValues(screen, 100, 450, 1.0, 20.0)
	// economy.GraphSupplyDemand(screen, 400, 750, 10.0, 1.0)
	economy.GraphGoodsVMoney(screen, 150, 750, 2.0, 5.0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	rand.Seed(time.Now().Unix())
	game := &Game{}

	ebiten.SetWindowSize(1040, 840)
	ebiten.SetWindowTitle("Economy Simulation")

	economy.Initialize(100)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
