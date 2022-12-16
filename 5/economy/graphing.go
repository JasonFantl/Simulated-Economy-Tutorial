package economy

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type dataPoint struct {
	min, max float64
}

var previousDataPoints map[Good]map[Location][]*dataPoint = map[Good]map[Location][]*dataPoint{}

func updateGraph() {

	datapoints := make(map[Location]map[Good]*dataPoint)
	for actor := range actors {
		for good, market := range actor.markets {
			if _, ok := datapoints[actor.location]; !ok {
				datapoints[actor.location] = make(map[Good]*dataPoint)
			}
			if _, ok := datapoints[actor.location][good]; !ok {
				datapoints[actor.location][good] = &dataPoint{market.expectedMarketPrice, market.expectedMarketPrice}
			} else {
				if market.expectedMarketPrice < datapoints[actor.location][good].min {
					datapoints[actor.location][good].min = market.expectedMarketPrice
				} else if market.expectedMarketPrice > datapoints[actor.location][good].max {
					datapoints[actor.location][good].max = market.expectedMarketPrice
				}
			}
		}
	}

	for location, goods := range datapoints {
		for good, datapoint := range goods {
			if _, ok := previousDataPoints[good]; !ok {
				previousDataPoints[good] = make(map[Location][]*dataPoint)
			}
			if _, ok := previousDataPoints[good][location]; !ok {
				previousDataPoints[good][location] = make([]*dataPoint, 0)
			}

			previousDataPoints[good][location] = append(previousDataPoints[good][location], datapoint)
		}
	}

}

func GraphExpectedValues(screen *ebiten.Image, title string, good Good, drawXOff, drawYOff, drawXZoom, drawYZoom float64, xRange, jumpXAxis, jumpYAxis int) {

	minX, maxX := math.MaxInt, 0
	maxY := 0.0

	for location := range previousDataPoints[good] {
		// add new data points (and get max/min X and Y values)
		v := len(previousDataPoints[good][location])
		if v > maxX {
			maxX = v
		}
		if len(previousDataPoints[good][location])-xRange < minX {
			minX = len(previousDataPoints[good][location]) - xRange
		}
		if minX < 0 {
			minX = 0
		}

		for i, v := range previousDataPoints[good][location] {
			if i > len(previousDataPoints[good][location])-xRange && v.max > maxY {
				maxY = v.max
			}
		}
	}

	// title
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff)+20)
	ebitenutil.DebugPrintAt(screen, "(Green = Riverwood, Blue = Seaside)", int(drawXOff), int(drawYOff)+35)

	// X axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX), drawYOff, color.White)
	for i := 0; i < (maxX-minX)+jumpXAxis; i += jumpXAxis {
		lowerRounded := ((i + minX) / jumpXAxis) * jumpXAxis
		if lowerRounded < minX {
			continue
		}
		x := int(drawXOff + drawXZoom*(float64(lowerRounded-minX)))
		y := int(drawYOff)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", lowerRounded), x, y)
	}

	// Y axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*maxY, color.White)
	for i := 0; i < int(maxY)+jumpYAxis; i += jumpYAxis {
		x := int(drawXOff)
		y := int(drawYOff - drawYZoom*(float64(i)))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), x-20, y-5)
	}

	// graph data
	for location := range previousDataPoints[good] {
		i := 0
		if len(previousDataPoints[good][location]) > xRange {
			i = len(previousDataPoints[good][location]) - xRange
		}
		for ; i < len(previousDataPoints[good][location]); i++ {
			// expected values
			datapoint := previousDataPoints[good][location][i]
			x, y := drawXOff+drawXZoom*float64(i-minX), drawYOff-drawYZoom*(datapoint.min+datapoint.max)/2.0

			col := color.RGBA{58, 158, 33, 100}
			if location == SEASIDE {
				col = color.RGBA{10, 159, 227, 100}
			}

			w := 1.0
			h := datapoint.max - datapoint.min
			if h < 3 {
				h = 3
			}
			ebitenutil.DrawRect(screen, x-w/2.0, y-h/2.0, w, h, col)

			// ebitenutil.DrawLine(screen, x0, y0, x1, y1, col)
		}
	}
}

func GraphGoodsVMoney(screen *ebiten.Image, title string, good Good, drawXOff, drawYOff, drawXZoom, drawYZoom float64, jumpXAxis, jumpYAxis int) {

	type dataPoint struct {
		x, y float64
		col  color.Color
	}

	minX, maxX := 0.0, 0.0
	minY, maxY := 0.0, 0.0

	points := make([]dataPoint, len(actors))
	i := 0
	for actor := range actors {
		x := actor.money
		y := float64(actor.markets[good].ownedGoods)
		col := color.RGBA{58, 158, 33, 100}
		if actor.location == SEASIDE {
			col = color.RGBA{10, 159, 227, 100}
		}
		points[i] = dataPoint{
			x:   x,
			y:   y,
			col: col,
		}
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
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff)+20)
	ebitenutil.DebugPrintAt(screen, "(Green = Riverwood, Blue = Seaside)", int(drawXOff), int(drawYOff)+35)

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
		x := drawXOff + drawXZoom*point.x
		y := drawYOff - drawYZoom*point.y
		w := 5.0
		h := 5.0

		ebitenutil.DrawRect(screen, x-w/2.0, y-h/2.0, w, h, point.col)
	}
}
