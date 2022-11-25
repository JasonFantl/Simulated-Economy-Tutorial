package economy

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var persistentScreen *ebiten.Image = nil

type dataPoint struct {
	expected, personal float64
	buyer              bool
}

func updateGraph() {
	for actor := range actors {
		if _, ok := previousDataPoints[actor]; !ok {
			previousDataPoints[actor] = make([]dataPoint, 0)
		}

		previousDataPoints[actor] = append(previousDataPoints[actor], dataPoint{actor.expectedMarketValue, actor.currentValue(), actor.expectedMarketValue < actor.currentValue()})
	}
}

var previousDataPoints map[*Actor][]dataPoint = map[*Actor][]dataPoint{}
var theoreticalPrices []float64 = []float64{}

func GraphExpectedValues(screen *ebiten.Image, drawXOff, drawYOff, drawXZoom, drawYZoom float64) {

	// setup persistent screen
	if persistentScreen == nil {
		persistentScreen = ebiten.NewImage(screen.Size())
	}

	maxX := 0.0
	maxY := 0.0

	// add new data points (and get max/min X and Y values)
	for actor := range actors {
		v := float64(len(previousDataPoints[actor]))
		if v > maxX {
			maxX = v
		}
		for _, v := range previousDataPoints[actor] {
			if v.expected > maxY {
				maxY = v.expected
			}
		}
	}

	// title
	ebitenutil.DebugPrintAt(persistentScreen, "Expected Market Values (Green = Buyer, Red = Seller)", int(drawXOff), int(drawYOff)+20)
	ebitenutil.DebugPrintAt(persistentScreen, "Personal Values (Pink)", int(drawXOff), int(drawYOff)+35)

	// X axis
	ebitenutil.DrawLine(persistentScreen, drawXOff, drawYOff, drawXOff+drawXZoom*maxX, drawYOff, color.White)
	for i := 0; i <= int(maxX); i += 50 {
		x := int(drawXOff + drawXZoom*(float64(i)))
		y := int(drawYOff)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), x, y)
	}

	// Y axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*maxY, color.White)
	jumpY := 5
	for i := 0; i < int(maxY)+jumpY; i += jumpY {
		x := int(drawXOff)
		y := int(drawYOff - drawYZoom*(float64(i)))
		ebitenutil.DebugPrintAt(persistentScreen, fmt.Sprintf("%d", i), x-20, y-5)
	}

	// graph data
	for actor := range actors {
		iteration := len(previousDataPoints[actor])
		if iteration > 1 {
			// expected values
			x0, y0 := drawXOff+drawXZoom*float64(iteration), drawYOff-drawYZoom*previousDataPoints[actor][iteration-2].expected
			x1, y1 := drawXOff+drawXZoom*float64(iteration+1), drawYOff-drawYZoom*previousDataPoints[actor][iteration-1].expected
			col := color.RGBA{143, 12, 3, 100}
			if previousDataPoints[actor][iteration-1].buyer {
				col = color.RGBA{50, 135, 0, 100}
			}
			ebitenutil.DrawLine(persistentScreen, x0, y0, x1, y1, col)

			// personal values
			x0, y0 = drawXOff+drawXZoom*float64(iteration), drawYOff-drawYZoom*previousDataPoints[actor][iteration-2].personal
			x1, y1 = drawXOff+drawXZoom*float64(iteration+1), drawYOff-drawYZoom*previousDataPoints[actor][iteration-1].personal
			col = color.RGBA{166, 0, 191, 100}
			ebitenutil.DrawLine(persistentScreen, x0, y0, x1, y1, col)
		}
	}

	// graph theoretical values
	if len(theoreticalPrices) > 0 {
		li := len(theoreticalPrices) - 1
		x0, y0 := drawXOff+drawXZoom*float64(li), drawYOff-drawYZoom*theoreticalPrices[len(theoreticalPrices)-1]
		col := color.RGBA{24, 100, 222, 100}
		ebitenutil.DrawRect(persistentScreen, x0, y0-5, 10, 10, col)
	}

	screen.DrawImage(persistentScreen, nil)
}

func GraphGoodsVMoney(screen *ebiten.Image, drawXOff, drawYOff, drawXZoom, drawYZoom float64) {

	minX, maxX := 0.0, 0.0
	minY, maxY := 0.0, 0.0

	points := make([][]float64, len(actors))
	i := 0
	for actor := range actors {
		x := actor.money
		y := float64(actor.ownedGoods)
		points[i] = []float64{x, y}
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}

		i++
	}

	// title
	ebitenutil.DebugPrintAt(screen, "Money", int(drawXOff)+40, int(drawYOff)+20)
	ebitenutil.DebugPrintAt(screen, "Goods", int(drawXOff)-80, int(drawYOff)-60)

	// X axis
	jumpXAxis := 10
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX), drawYOff, color.White)
	for i := int(minX); i <= int(float64(maxX+1)); i += jumpXAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff+drawXZoom*float64(i-int(minX))), int(drawYOff))
	}

	// Y axis
	jumpYAxis := 5
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*float64(maxY-minY), color.White)
	for i := int(minY); i <= int(maxY+1); i += jumpYAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff)-20, int(drawYOff-drawYZoom*float64(i-int(minY))))
	}

	// 2d plot
	for _, point := range points {
		x := drawXOff + drawXZoom*point[0]
		y := drawYOff - drawYZoom*point[1]
		w := 10.0
		h := 10.0
		col := color.RGBA{143, 12, 3, 100}
		ebitenutil.DrawRect(screen, x-w/2.0, y-h/2.0, w, h, col)
	}
}
