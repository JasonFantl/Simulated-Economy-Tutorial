package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jasonfantl/SimulatedEconomy6/economy"
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
			economy.Influence(economy.RIVERWOOD, 1)
		} else if p == ebiten.KeyArrowDown {
			economy.Influence(economy.RIVERWOOD, -1)
		} else if p == ebiten.KeyArrowLeft {
			economy.Influence(economy.SEASIDE, 1)
		} else if p == ebiten.KeyArrowRight {
			economy.Influence(economy.SEASIDE, -1)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	economy.GraphExpectedValues(screen, "Price of Wood", economy.WOOD, 100, 300, 0.5, 40.0, 800, 100, 1)
	economy.GraphExpectedValues(screen, "Price of Chairs", economy.CHAIR, 100, 600, 0.5, 8.0, 800, 100, 5)
	economy.GraphExpectedValues(screen, "Price of Thread", economy.THREAD, 600, 300, 0.5, 30.0, 800, 100, 1)
	economy.GraphExpectedValues(screen, "Price of Bed", economy.BED, 600, 600, 0.5, 4.0, 800, 100, 10)

	// economy.GraphGoodsVMoney(screen, "Wood V Money", economy.WOOD, 600, 200, 0.1, 4.0, 250, 10)
	// economy.GraphGoodsVMoney(screen, "Chair V Money", economy.CHAIR, 600, 400, 0.1, 4.0, 250, 10)
	// economy.GraphGoodsVMoney(screen, "Thread V Money", economy.THREAD, 600, 600, 0.1, 2.0, 250, 10)
	// economy.GraphGoodsVMoney(screen, "Bed V Money", economy.BED, 600, 800, 0.1, 20.0, 250, 1)

	economy.GraphLeisureVWealth(screen, "Leisure V Wealth", 500, 800, 0.1, 10, 250, 2)

	economy.GraphMerchantType(screen, "Merchant types", 100, 800, 50, 5)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	rand.Seed(time.Now().Unix())
	game := &Game{}

	ebiten.SetWindowSize(1440, 940)
	ebiten.SetWindowTitle("Economy Simulation")

	economy.Initialize(100)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
