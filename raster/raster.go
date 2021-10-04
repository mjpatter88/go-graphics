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

	drawTriangle(screen[:], color{0x99, 0xB8, 0x98}, point{-200, -250}, point{200, 50}, point{20, 250})
	drawWireframeTriangle(screen[:], color{0xe8, 0x4a, 0x5f}, point{-200, -250}, point{200, 50}, point{20, 250})
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

func drawTriangle(screen []byte, color color, p0, p1, p2 point) {
	// First sort the vertices so p0 is the lowest point and p2 the highest.
	// This sorting also guarantees that p0 -> p2 is always a "tall" side.
	// The other two sides are "short" sides.

	if p1.y < p0.y {
		p0, p1 = swap(p0, p1)
	}
	if p2.y < p0.y {
		p0, p2 = swap(p0, p2)
	}
	if p2.y < p1.y {
		p1, p2 = swap(p1, p2)
	}

	// Calculate the sets of x values for each side.
	x01 := interpolate(p0.y, p0.x, p1.y, p1.x)
	x12 := interpolate(p1.y, p1.x, p2.y, p2.x)
	x02 := interpolate(p0.y, p0.x, p2.y, p2.x)

	// Since x02 is a long side, we don't have to do anything to it.
	// For the other side, we want to concat x01 and x12.
	// There is a single value of overlap, so drop the last item from x01.
	x012 := append(x01[:len(x01)-1], x12...)

	// Determine which set of x values is the left and which is the right
	mid := len(x02) / 2
	xLeft := x02
	xRight := x012
	if x02[mid] >= x012[mid] {
		xLeft = x012
		xRight = x02
	}

	// Draw horizontal lines from top to bottom
	drawHorizontalLines(screen, color, p0.y, p2.y, xLeft, xRight)

}

func drawHorizontalLines(screen []byte, color color, yStart int, yEnd int, leftXs []float64, rightXs []float64) {
	for y := yStart; y <= yEnd; y++ {
		for x := leftXs[y-yStart]; x <= rightXs[y-yStart]; x++ {
			putPixel(screen, color, int(x), y)
		}
	}
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
