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

		previousDataPoints[actor] = append(previousDataPoints[actor], dataPoint{actor.expectedMarketValue, actor.personalValue, actor.expectedMarketValue < actor.personalValue})
	}

	// calculate theoretical value
	bucketSize := 0.2
	supply, demand, xOff := calculateSupplyDemand(bucketSize)
	theoreticalValue := 0.0
	for i, v := range supply {
		if v >= demand[i] {
			theoreticalValue = bucketSize * float64(i+xOff)
			break
		}
	}
	theoreticalPrices = append(theoreticalPrices, theoreticalValue)
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

func GraphSupplyDemand(screen *ebiten.Image, drawXOff, drawYOff, drawXZoom, drawYZoom float64) {
	bucketSize := 0.2
	supply, demand, xOff := calculateSupplyDemand(bucketSize)

	minX, maxX := intMin(0, xOff), intMax(0, len(supply)-xOff)
	minY, maxY := intMin(0, supply[0]), demand[0]

	// title
	ebitenutil.DebugPrintAt(persistentScreen, "Demand and Supply (Green = Demand, Red = Supply)", int(drawXOff), int(drawYOff)+20)

	// X axis
	jumpXAxis := 2
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX)*bucketSize, drawYOff, color.White)
	for i := int(minX); i <= int(float64(maxX+1)*bucketSize); i += jumpXAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff+drawXZoom*float64(i-int(minX))), int(drawYOff))
	}

	// Y axis
	jumpYAxis := 50
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*float64(maxY-minY), color.White)
	for i := int(minY); i <= int(maxY+1); i += jumpYAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff)-20, int(drawYOff-drawYZoom*float64(i-int(minY))))
	}

	// histogram
	for i := 0; i < len(supply); i++ {
		x := drawXOff + drawXZoom*float64(i)*bucketSize
		y := drawYOff
		w := drawXZoom * float64(bucketSize)
		// supply
		h := -drawYZoom * float64(supply[i]-minY)
		col := color.RGBA{143, 12, 3, 100}
		ebitenutil.DrawRect(screen, x, y, w, h, col)
		// demand
		h = -drawYZoom * float64(demand[i]-minY)
		col = color.RGBA{50, 135, 0, 100}
		ebitenutil.DrawRect(screen, x, y, w, h, col)

	}
}

// supply, demand, x offset
func calculateSupplyDemand(bucketSize float64) ([]int, []int, int) {
	// bucket everything
	buckets := make(map[int]int)
	for a := range actors {
		bucketIndex := int(a.personalValue / bucketSize)
		buckets[bucketIndex]++
	}

	// find graph info
	minX, maxX := 0, 0
	for x, _ := range buckets {
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
	}

	// create supply and demand curves by integrating over personal value
	// not an efficient solution, but works
	supply := make([]int, maxX-minX+1)
	demand := make([]int, maxX-minX+1)
	for x, y := range buckets {
		for i := x; i <= maxX; i++ {
			supply[i-minX] += y
		}
		for i := x; i >= minX; i-- {
			demand[i-minX] += y
		}
	}
	return supply, demand, minX
}
