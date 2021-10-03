package main

type color struct {
	r byte
	g byte
	b byte
}

const windowWidth int = 1000
const windowHeight int = 1000

var backgroundColor = color{0x00, 0x00, 0x00}

// X and Y are canvas coordinates
// (0,0 in middle, -(Width/2), -(Hight/2) in bottom left).
func putPixel(screen *[windowWidth * windowHeight * 4]byte, color color, x int, y int) {
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
			putPixel(screen, backgroundColor, x, y)
		}
	}

}
