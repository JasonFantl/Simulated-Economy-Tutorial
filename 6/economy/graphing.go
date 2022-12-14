package economy

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var locationColors = map[Location]color.Color{
	RIVERWOOD:  color.RGBA{58, 158, 33, 100},
	SEASIDE:    color.RGBA{10, 159, 227, 100},
	WINTERHOLD: color.RGBA{255, 255, 255, 100},
	PORTSVILLE: color.RGBA{128, 0, 0, 100},
}

type dataPoint struct {
	min, max float64
}

var previousDataPoints map[Good]map[Location][]*dataPoint = map[Good]map[Location][]*dataPoint{}

func updateGraph() {

	datapoints := make(map[Location]map[Good]*dataPoint)
	for _, location := range locations {
		datapoints[location] = make(map[Good]*dataPoint)
		for _, good := range goods {
			datapoints[location][good] = &dataPoint{math.MaxFloat64, -math.MaxFloat64}
		}
	}
	for local := range locals {
		for good, market := range local.markets {
			if good == LEISURE {
				continue
			}
			if market.expectedMarketPrice < datapoints[local.location][good].min {
				datapoints[local.location][good].min = market.expectedMarketPrice
			}
			if market.expectedMarketPrice > datapoints[local.location][good].max {
				datapoints[local.location][good].max = market.expectedMarketPrice
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
	ebitenutil.DebugPrintAt(screen, "(Green = Riverwood, Blue = Seaside, White=Winterhold, Red=Portsville)", int(drawXOff), int(drawYOff)+35)

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

			col := locationColors[location]

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

	points := make([]dataPoint, 0)
	for local := range locals {
		x := local.money
		y := float64(local.markets[good].ownedGoods)
		col := locationColors[local.location]
		points = append(points, dataPoint{
			x:   x,
			y:   y,
			col: col,
		})
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
	}

	for merchant := range merchants {
		if good != merchant.buysSells {
			continue
		}
		x := merchant.money
		y := float64(merchant.owned)
		r, g, b, _ := locationColors[merchant.location].RGBA()

		points = append(points, dataPoint{
			x:   x,
			y:   y,
			col: color.RGBA{uint8(r), uint8(g), uint8(b), 255},
		})
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

func GraphLeisureVWealth(screen *ebiten.Image, title string, drawXOff, drawYOff, drawXZoom, drawYZoom float64, jumpXAxis, jumpYAxis int) {

	type dataPoint struct {
		x, y float64
		col  color.Color
	}

	minX, maxX := 0.0, 0.0
	minY, maxY := 0.0, 0.0

	points := make([]dataPoint, 0)
	for local := range locals {
		x := local.money
		// for _, market := range local.markets {
		// 	x += market.expectedMarketPrice * float64(market.ownedGoods)
		// }
		y := float64(local.markets[LEISURE].basePersonalValue)
		col := locationColors[local.location]

		points = append(points, dataPoint{
			x:   x,
			y:   y,
			col: col,
		})
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
	}

	// title
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff)+20)

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

func GraphMerchantType(screen *ebiten.Image, title string, drawXOff, drawYOff, drawXZoom, drawYZoom float64) {

	points := make(map[Good]int)
	for _, good := range goods {
		points[good] = 0
	}
	for merchant := range merchants {
		points[merchant.buysSells]++
	}

	// title
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff)+60)

	// 2d plot
	xIndex := 0.0
	for _, good := range goods {
		x := drawXOff + drawXZoom*xIndex
		y := drawYOff
		w := drawXZoom * 0.9
		h := float64(points[good]) * drawYZoom

		ebitenutil.DrawRect(screen, x, y-h, w, h, color.RGBA{100, 100, 100, 255})
		ebitenutil.DebugPrintAt(screen, string(good), int(x), int(y)+20)
		ebitenutil.DebugPrintAt(screen, strconv.Itoa(points[good]), int(x), int(y)+40)

		xIndex++
	}
}
