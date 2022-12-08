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

			previousDataPoints[good][actor] = append(previousDataPoints[good][actor], dataPoint{market.expectedMarketPrice, actor.valueToPrice(actor.currentPersonalValue(good)), actor.isBuyer(good)})
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
	ebitenutil.DebugPrintAt(persistentScreen, "Expected Market Values (Green = Buyer, Red = Seller)", int(drawXOff), int(drawYOff)+35)
	ebitenutil.DebugPrintAt(persistentScreen, "Personal Values (Pink)", int(drawXOff), int(drawYOff)+50)

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
