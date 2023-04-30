package economy

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type dataPoint struct {
	min, max float64
	color    color.Color
}

var previousDataPoints map[Good]map[cityName][]*dataPoint = make(map[Good]map[cityName][]*dataPoint)

func updateGraph(city *City) {

	datapoints := make(map[Good]*dataPoint)
	for _, good := range goods {
		datapoints[good] = &dataPoint{math.MaxFloat64, -math.MaxFloat64, city.color}
	}

	for local := range city.locals {
		for good, market := range local.markets {
			if good == LEISURE {
				continue
			}
			if market.expectedMarketPrice < datapoints[good].min {
				datapoints[good].min = market.expectedMarketPrice
			}
			if market.expectedMarketPrice > datapoints[good].max {
				datapoints[good].max = market.expectedMarketPrice
			}
		}
	}

	for good, datapoint := range datapoints {
		if _, ok := previousDataPoints[good]; !ok {
			previousDataPoints[good] = make(map[cityName][]*dataPoint)
		}
		if _, ok := previousDataPoints[good][city.name]; !ok {
			previousDataPoints[good][city.name] = make([]*dataPoint, 0)
		}
		previousDataPoints[good][city.name] = append(previousDataPoints[good][city.name], datapoint)
	}
}

// GraphExpectedValues will graph the expected values of a specific city
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

			w := 1.0
			h := datapoint.max - datapoint.min
			if h < 3 {
				h = 3
			}
			ebitenutil.DrawRect(screen, x-w/2.0, y-h/2.0, w, h, datapoint.color)
		}
	}
}

// GraphGoodsVMoney will graph a point for each resident, comparing their goods to money
func GraphGoodsVMoney(screen *ebiten.Image, city City, title string, good Good, drawXOff, drawYOff, drawXZoom, drawYZoom float64, jumpXAxis, jumpYAxis int) {

	type dataPoint struct {
		x, y float64
		col  color.Color
	}

	minX, maxX := 0.0, 0.0
	minY, maxY := 0.0, 0.0

	points := make([]dataPoint, 0)
	for local := range city.locals {
		x := local.money
		y := float64(local.markets[good].ownedGoods)
		points = append(points, dataPoint{
			x:   x,
			y:   y,
			col: city.color,
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

	for merchant := range city.merchants {
		if good != merchant.BuysSells {
			continue
		}
		x := merchant.Money
		y := float64(merchant.Owned)
		r, g, b, _ := city.color.RGBA()

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

// GraphLeisureVWealth will graph a point for each resident, comparing their value of leisure to their wealth
func GraphLeisureVWealth(screen *ebiten.Image, city City, title string, drawXOff, drawYOff, drawXZoom, drawYZoom float64, jumpXAxis, jumpYAxis int) {

	type dataPoint struct {
		x, y float64
		col  color.Color
	}

	minX, maxX := 0.0, 0.0
	minY, maxY := 0.0, 0.0

	points := make([]dataPoint, 0)
	for local := range city.locals {
		x := local.money
		// for _, market := range local.markets {
		// 	x += market.expectedMarketPrice * float64(market.ownedGoods)
		// }
		y := float64(local.markets[LEISURE].basePersonalValue)

		points = append(points, dataPoint{
			x:   x,
			y:   y,
			col: city.color,
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

// GraphMerchantType will graph the number of all the different merchant types
func GraphMerchantType(screen *ebiten.Image, cities []City, title string, drawXOff, drawYOff, drawXZoom, drawYZoom float64) {

	points := make(map[Good]map[cityName]int)
	totals := make(map[Good]int)

	for _, good := range goods {
		points[good] = make(map[cityName]int)
		for _, city := range cities {
			points[good][city.name] = 0
		}
	}
	for _, city := range cities {
		for merchant := range city.merchants {
			points[merchant.BuysSells][city.name]++
			totals[merchant.BuysSells]++
		}
	}

	// title
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff)+60)

	// 2d plot
	xIndex := 0.0
	for _, good := range goods {
		yOff := 0.0
		x := drawXOff + drawXZoom*xIndex
		w := drawXZoom * 0.9
		for _, city := range cities {
			y := drawYOff + yOff
			h := float64(points[good][city.name]) * drawYZoom

			ebitenutil.DrawRect(screen, x, y-h, w, h, city.color)

			yOff -= h
		}
		ebitenutil.DebugPrintAt(screen, string(good), int(x), int(drawYOff)+20)
		ebitenutil.DebugPrintAt(screen, strconv.Itoa(totals[good]), int(x), int(drawYOff)+40)
		xIndex++
	}
}
