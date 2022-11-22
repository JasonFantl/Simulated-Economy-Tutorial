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

var previousDataPoints map[*Actor][]dataPoint = map[*Actor][]dataPoint{}
var averages []float64 = []float64{}

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
		for i := intMax(0, len(previousDataPoints[actor])-2); i < len(previousDataPoints[actor])-1; i++ {
			// expected values
			x0, y0 := drawXOff+drawXZoom*float64(i), drawYOff-drawYZoom*previousDataPoints[actor][i].expected
			x1, y1 := drawXOff+drawXZoom*float64(i+1), drawYOff-drawYZoom*previousDataPoints[actor][i+1].expected
			col := color.RGBA{143, 12, 3, 100}
			if previousDataPoints[actor][i+1].buyer {
				col = color.RGBA{50, 135, 0, 100}
			}
			ebitenutil.DrawLine(persistentScreen, x0, y0, x1, y1, col)

			// personal values
			x0, y0 = drawXOff+drawXZoom*float64(i), drawYOff-drawYZoom*previousDataPoints[actor][i].personal
			x1, y1 = drawXOff+drawXZoom*float64(i+1), drawYOff-drawYZoom*previousDataPoints[actor][i+1].personal
			col = color.RGBA{166, 0, 191, 100}
			ebitenutil.DrawLine(persistentScreen, x0, y0, x1, y1, col)
		}
	}

	screen.DrawImage(persistentScreen, nil)

}
