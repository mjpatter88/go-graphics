package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	window, renderer := initialize()
	tex, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		int32(windowWidth),
		int32(windowHeight),
	)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()
	defer window.Destroy()
	defer renderer.Destroy()

	var screen = [windowWidth * windowHeight * 4]byte{}

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
		rayTraceFrame(&screen)
		drawFrame(renderer, tex, &screen)
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

func drawFrame(renderer *sdl.Renderer, tex *sdl.Texture, screen *[windowWidth * windowHeight * 4]byte) {
	bytes, _, err := tex.Lock(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(windowWidth*windowHeight*4); i++ {
		bytes[i] = screen[i]
	}
	tex.Unlock()
	rect := sdl.Rect{X: 0, Y: 0, W: int32(windowWidth), H: int32(windowHeight)}
	err = renderer.Copy(tex, nil, &rect)
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
		int32(windowWidth),
		int32(windowHeight),
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

func clearFrame(renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
}
