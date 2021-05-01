package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Render at roughly 60 fps
const msPerFrame = 16

const windowWidth int32 = 1000
const windowHeight int32 = 1000

func clearFrame(renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
}

// X and Y are canvas coordinates
// (0,0 in middle, -(Width/2), -(Hight/2) in bottom left).
func putPixel(screen *[windowWidth * windowHeight * 4]byte, r byte, g byte, b byte, x int32, y int32) {
	screenX := (windowWidth / 2) + x
	screenY := (windowHeight / 2) + y
	base := (screenY*windowWidth + screenX) * 4
	screen[base] = r
	screen[base+1] = g
	screen[base+2] = b
	screen[base+3] = 0xFF
	screen[0] = 0xFF
}

func rayTraceFrame(tex *sdl.Texture) {
	w := windowWidth
	h := windowHeight
	var r byte = 0xFF
	var b byte = 0xFF
	var g byte = 0xFF

	var screen = [windowWidth * windowHeight * 4]byte{}

	for x := -(windowWidth / 2); x < (windowWidth / 2); x++ {
		for y := -(windowHeight / 2); y < (windowHeight / 2); y++ {
			putPixel(&screen, r, g, b, x, y)
		}
	}

	bytes, _, err := tex.Lock(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(w*h*4); i++ {
		bytes[i] = screen[i]
	}
	tex.Unlock()
}

func drawFrame(renderer *sdl.Renderer, tex *sdl.Texture) {
	rect := sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}
	err := renderer.Copy(tex, nil, &rect)
	if err != nil {
		panic(err)
	}

	renderer.Present()
}

func updateOjects() {
}

func initialize() (*sdl.Window, *sdl.Renderer) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer(
		windowWidth,
		windowHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	window.SetTitle("Go Ray Tracing")
	if err != nil {
		panic(err)
	}
	return window, renderer
}

func handleInput(state []uint8) {
	return
}

func main() {
	window, renderer := initialize()
	tex, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		windowWidth,
		windowHeight,
	)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()
	defer window.Destroy()
	defer renderer.Destroy()

	var frameCount int = 0
	framesProcessed := 0
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				frames := frameCount - framesProcessed
				framesProcessed = frameCount
				fmt.Println("fps: ", frames)
			}
		}
	}()

	running := true
	for running {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				running = false
				done <- true
			}
		}
		handleInput(sdl.GetKeyboardState())
		clearFrame(renderer)
		rayTraceFrame(tex)
		drawFrame(renderer, tex)
		frameCount++

		elapsed := time.Since(frameStart).Milliseconds()
		if elapsed < msPerFrame {
			delay := msPerFrame - elapsed
			if delay < 0 {
				fmt.Println(delay)
				panic(delay)
			}
			sdl.Delay(uint32(delay))
		}
	}
}
