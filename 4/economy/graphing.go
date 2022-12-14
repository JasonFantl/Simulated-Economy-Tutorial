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
		for good, market := range actor.markets {
			if _, ok := previousDataPoints[good]; !ok {
				previousDataPoints[good] = make(map[*Actor][]dataPoint)
			}
			if _, ok := previousDataPoints[good][actor]; !ok {
				previousDataPoints[good][actor] = make([]dataPoint, 0)
			}

			previousDataPoints[good][actor] = append(previousDataPoints[good][actor], dataPoint{
				market.expectedMarketPrice,
				actor.valueToPrice(actor.currentPersonalValue(good)),
				actor.isBuyer(good),
			})
		}
	}
}

var previousDataPoints map[Good]map[*Actor][]dataPoint = map[Good]map[*Actor][]dataPoint{}

func GraphExpectedValues(screen *ebiten.Image, title string, good Good, drawXOff, drawYOff, drawXZoom, drawYZoom float64, jumpXAxis, jumpYAxis int) {

	// setup persistent screen
	if persistentScreen == nil {
		persistentScreen = ebiten.NewImage(screen.Size())
	}

	maxX := 0.0
	maxY := 0.0

	// add new data points (and get max/min X and Y values)
	for actor := range actors {
		v := float64(len(previousDataPoints[good][actor]))
		if v > maxX {
			maxX = v
		}
		for _, v := range previousDataPoints[good][actor] {
			if v.expected > maxY {
				maxY = v.expected
			}
		}
	}

	// title
	ebitenutil.DebugPrintAt(persistentScreen, title, int(drawXOff), int(drawYOff)+20)
	ebitenutil.DebugPrintAt(persistentScreen, "Expected Market Price (Green = Buyer, Red = Seller)", int(drawXOff), int(drawYOff)+35)
	ebitenutil.DebugPrintAt(persistentScreen, "Personal Price (Pink)", int(drawXOff), int(drawYOff)+50)

	// X axis
	ebitenutil.DrawLine(persistentScreen, drawXOff, drawYOff, drawXOff+drawXZoom*maxX, drawYOff, color.White)
	for i := 0; i < int(maxX)+jumpXAxis; i += jumpXAxis {
		x := int(drawXOff + drawXZoom*(float64(i)))
		y := int(drawYOff)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), x, y)
	}

	// Y axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*maxY, color.White)
	for i := 0; i < int(maxY)+jumpYAxis; i += jumpYAxis {
		x := int(drawXOff)
		y := int(drawYOff - drawYZoom*(float64(i)))
		ebitenutil.DebugPrintAt(persistentScreen, fmt.Sprintf("%d", i), x-20, y-5)
	}

	// graph data
	for actor := range actors {
		iteration := len(previousDataPoints[good][actor])
		if iteration > 1 {
			// expected values
			x0, y0 := drawXOff+drawXZoom*float64(iteration), drawYOff-drawYZoom*previousDataPoints[good][actor][iteration-2].expected
			x1, y1 := drawXOff+drawXZoom*float64(iteration+1), drawYOff-drawYZoom*previousDataPoints[good][actor][iteration-1].expected
			col := color.RGBA{143, 12, 3, 100}
			if previousDataPoints[good][actor][iteration-1].buyer {
				col = color.RGBA{50, 135, 0, 100}
			}
			ebitenutil.DrawLine(persistentScreen, x0, y0, x1, y1, col)

			// personal values
			x0, y0 = drawXOff+drawXZoom*float64(iteration), drawYOff-drawYZoom*previousDataPoints[good][actor][iteration-2].personal
			x1, y1 = drawXOff+drawXZoom*float64(iteration+1), drawYOff-drawYZoom*previousDataPoints[good][actor][iteration-1].personal
			col = color.RGBA{166, 0, 191, 100}
			ebitenutil.DrawLine(persistentScreen, x0, y0, x1, y1, col)
		}
	}

	screen.DrawImage(persistentScreen, nil)
}

func GraphGoodsVMoney(screen *ebiten.Image, title string, good Good, drawXOff, drawYOff, drawXZoom, drawYZoom float64, jumpXAxis, jumpYAxis int) {

	minX, maxX := 0.0, 0.0
	minY, maxY := 0.0, 0.0

	points := make([][]float64, len(actors))
	i := 0
	for actor := range actors {
		x := actor.money
		y := float64(actor.markets[good].ownedGoods)
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
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX), drawYOff, color.White)
	for i := int(minX); i <= int(float64(maxX+1)); i += jumpXAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff+drawXZoom*float64(i-int(minX))), int(drawYOff))
	}

	// Y axis
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
