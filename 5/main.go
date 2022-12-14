package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jasonfantl/SimulatedEconomy5/economy"
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
			economy.Influence(economy.CHAIR, 1)
		} else if p == ebiten.KeyArrowDown {
			economy.Influence(economy.CHAIR, -1)
		} else if p == ebiten.KeyArrowLeft {
			economy.Influence(economy.WOOD, 1)
		} else if p == ebiten.KeyArrowRight {
			economy.Influence(economy.WOOD, -1)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	economy.GraphExpectedValues(screen, "Price of Wood", economy.RIVERWOOD, economy.WOOD, 50, 200, 0.5, 20.0, 100, 1)
	economy.GraphExpectedValues(screen, "Price of Chairs", economy.RIVERWOOD, economy.CHAIR, 50, 500, 0.5, 4.0, 100, 5)

	economy.GraphExpectedValues(screen, "Price of Wood", economy.SEASIDE, economy.WOOD, 500, 200, 0.5, 20.0, 100, 1)
	economy.GraphExpectedValues(screen, "Price of Chairs", economy.SEASIDE, economy.CHAIR, 500, 500, 0.5, 4.0, 100, 5)

	economy.GraphGoodsVMoney(screen, "Wood V Money", economy.WOOD, 1000, 200, 0.1, 2.0, 250, 10)
	economy.GraphGoodsVMoney(screen, "Chair V Money", economy.CHAIR, 1000, 500, 0.1, 2.0, 250, 10)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	rand.Seed(time.Now().Unix())
	game := &Game{}

	ebiten.SetWindowSize(1440, 840)
	ebiten.SetWindowTitle("Economy Simulation")

	economy.Initialize(100)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
