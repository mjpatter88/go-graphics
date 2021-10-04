package main

import (
	"math"
)

type color struct {
	r byte
	g byte
	b byte
}

type point struct {
	x int
	y int
}

type vert struct {
	x int
	y int
	h float64
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

	v1 := vert{-200, -250, 0.5}
	v2 := vert{200, 50, 0.0}
	v3 := vert{20, 250, 1.0}
	drawShadedTriangle(screen[:], color{0xe8, 0x4a, 0x5f}, v1, v2, v3)
}

func drawShadedTriangle(screen []byte, color color, v0, v1, v2 vert) {
	// First sort the vertices so p0 is the lowest point and p2 the highest.
	// This sorting also guarantees that p0 -> p2 is always a "tall" side.
	// The other two sides are "short" sides.

	if v1.y < v0.y {
		v0, v1 = swapVerts(v0, v1)
	}
	if v2.y < v0.y {
		v0, v2 = swapVerts(v0, v2)
	}
	if v2.y < v1.y {
		v1, v2 = swapVerts(v1, v2)
	}

	// Calculate the sets of x values and h values for each side.
	x01 := interpolate(v0.y, float64(v0.x), v1.y, float64(v1.x))
	h01 := interpolate(v0.y, v0.h, v1.y, v1.h)

	x12 := interpolate(v1.y, float64(v1.x), v2.y, float64(v2.x))
	h12 := interpolate(v1.y, v1.h, v2.y, v2.h)

	x02 := interpolate(v0.y, float64(v0.x), v2.y, float64(v2.x))
	h02 := interpolate(v0.y, v0.h, v2.y, v2.h)

	// Since x02 is a long side, we don't have to do anything to it.
	// For the other side, we want to concat x01 and x12.
	// There is a single value of overlap, so drop the last item from x01.
	x012 := append(x01[:len(x01)-1], x12...)
	h012 := append(h01[:len(h01)-1], h12...)

	// Determine which set of x values is the left and which is the right
	mid := len(x02) / 2
	xLeft := x02
	hLeft := h02
	xRight := x012
	hRight := h012
	if x02[mid] >= x012[mid] {
		xLeft = x012
		xRight = x02
		hLeft = h012
		hRight = h02
	}

	// Draw horizontal lines from top to bottom
	drawShadedHorizontalLines(screen, color, v0.y, v2.y, xLeft, xRight, hLeft, hRight)
}

func drawShadedHorizontalLines(screen []byte, color color, yStart int, yEnd int, leftXs, rightXs, leftHs, rightHs []float64) {
	for y := yStart; y <= yEnd; y++ {
		yOffset := y - yStart
		xLeft := int(leftXs[yOffset])
		xRight := int(rightXs[yOffset])
		hVals := interpolate(int(xLeft), leftHs[yOffset], int(xRight), rightHs[yOffset])
		for x := xLeft; x <= xRight; x++ {
			shadedColor := scaleColor(color, hVals[x-xLeft])
			putPixel(screen, shadedColor, int(x), y)
		}
	}
}

func scaleColor(col color, factor float64) color {
	newR := float64(col.r) * factor
	newG := float64(col.g) * factor
	newB := float64(col.b) * factor
	return color{
		r: byte(math.Min(255, newR)),
		g: byte(math.Min(255, newG)),
		b: byte(math.Min(255, newB)),
	}
}

func drawLine(screen []byte, color color, start point, end point) {
	if abs(end.x-start.x) > abs(end.y-start.y) {
		// line is more horizontal than vertical
		if start.x > end.x {
			// Always start at the leftmost point
			start, end = swap(start, end)
		}

		yValues := interpolate(start.x, float64(start.y), end.x, float64(end.y))
		for x := start.x; x <= end.x; x++ {
			putPixel(screen, color, x, int(yValues[x-start.x]))
		}
	} else {
		// line is more vertical than horizontal
		if start.y > end.y {
			// Always start at the bottom point
			start, end = swap(start, end)
		}
		xValues := interpolate(start.y, float64(start.x), end.y, float64(end.x))
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
	x01 := interpolate(p0.y, float64(p0.x), p1.y, float64(p1.x))
	x12 := interpolate(p1.y, float64(p1.x), p2.y, float64(p2.x))
	x02 := interpolate(p0.y, float64(p0.x), p2.y, float64(p2.x))

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

func swapVerts(x vert, y vert) (vert, vert) {
	return y, x
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
func interpolate(iStart int, dStart float64, iEnd int, dEnd float64) []float64 {
	// If there is only one point
	if iStart == iEnd {
		return []float64{dStart}
	}

	values := make([]float64, (iEnd-iStart)+1)

	slope := float64(dEnd-dStart) / float64(iEnd-iStart)
	d := dStart

	for i := iStart; i <= iEnd; i++ {
		values[(i - iStart)] = d
		d += slope
	}

	return values
}
