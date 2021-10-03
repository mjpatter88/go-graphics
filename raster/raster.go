package main

type color struct {
	r byte
	g byte
	b byte
}

type point struct {
	x int
	y int
}

const windowWidth int = 1000
const windowHeight int = 1000

var backgroundColor = color{0x00, 0x00, 0x00}

// X and Y are canvas coordinates
// (0,0 in middle, -(Width/2), -(Hight/2) in bottom left).
func putPixel(screen []byte, color color, x int, y int) {
	screenX := (windowWidth / 2) + x
	screenY := (windowHeight / 2) - y - 1
	base := (screenY*windowWidth + screenX) * 4
	screen[base] = color.r
	screen[base+1] = color.g
	screen[base+2] = color.b
	screen[base+3] = 0xFF
	screen[0] = 0xFF
}

func rasterizeFrame(screen *[windowWidth * windowHeight * 4]byte) {

	for x := -(windowWidth / 2); x < (windowWidth / 2); x++ {
		for y := -(windowHeight / 2); y < (windowHeight / 2); y++ {
			putPixel(screen[:], backgroundColor, x, y)
		}
	}

	drawWireframeTriangle(screen[:], color{0xff, 0x00, 0xff}, point{-200, -250}, point{200, 50}, point{20, 250})
}

func drawLine(screen []byte, color color, start point, end point) {
	if abs(end.x-start.x) > abs(end.y-start.y) {
		// line is more horizontal than vertical
		if start.x > end.x {
			// Always start at the leftmost point
			start, end = swap(start, end)
		}

		yValues := interpolate(start.x, start.y, end.x, end.y)
		for x := start.x; x <= end.x; x++ {
			putPixel(screen, color, x, int(yValues[x-start.x]))
		}
	} else {
		// line is more vertical than horizontal
		if start.y > end.y {
			// Always start at the bottom point
			start, end = swap(start, end)
		}
		xValues := interpolate(start.y, start.x, end.y, end.x)
		for y := start.y; y <= end.y; y++ {
			putPixel(screen, color, int(xValues[y-start.y]), y)
		}

	}
}

func drawWireframeTriangle(screen []byte, color color, p0, p1, p2 point) {
	drawLine(screen, color, p0, p1)
	drawLine(screen, color, p1, p2)
	drawLine(screen, color, p2, p0)

}

func swap(x point, y point) (point, point) {
	return y, x
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x

}

// Calculate a series of values for a dependent variable by stepping betwween the start
// and end of an independent variable.
//
// Example: Given the start and end points of a line, return the series of y values
// calculated by stepping from x start to x end.
func interpolate(iStart int, dStart int, iEnd int, dEnd int) []float64 {
	// If there is only one point
	if iStart == iEnd {
		return []float64{float64(dStart)}
	}

	values := make([]float64, (iEnd-iStart)+1)

	slope := float64(dEnd-dStart) / float64(iEnd-iStart)
	d := float64(dStart)

	for i := iStart; i <= iEnd; i++ {
		values[(i - iStart)] = d
		d += slope
	}

	return values
}
