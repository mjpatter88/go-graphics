package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type vec3 struct {
	x float64
	y float64
	z float64
}

func vec3Sub(first vec3, second vec3) vec3 {
	return vec3{first.x - second.x, first.y - second.y, first.z - second.z}
}

func vec3Dot(first vec3, second vec3) float64 {
	return (first.x * second.x) + (first.y * second.y) + (first.z * second.z)
}

type color struct {
	r byte
	g byte
	b byte
}

type sphere struct {
	center vec3
	radius float64
	color  color
}

type solutions struct {
	first  float64
	second float64
}

// Try to render at roughly 60 fps
const msPerFrame = 16

const windowWidth int = 1000
const windowHeight int = 1000
const viewportWidth int = 1
const viewportHeight int = 1
const distCameraToProjectionPlane float64 = 1

var backgroundColor = color{0xFF, 0xFF, 0xFF}

var sphere1 = sphere{vec3{0, -1, 3}, 1, color{255, 0, 0}}
var sphere2 = sphere{vec3{2, 0, 4}, 1, color{0, 0, 255}}
var sphere3 = sphere{vec3{-2, 0, 4}, 1, color{0, 255, 0}}

var shapes = [...]sphere{sphere1, sphere2, sphere3}

func clearFrame(renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
}

func canvasToViewport(x int, y int) vec3 {
	vx := (float64(x) * float64(viewportWidth) / float64(windowWidth))
	vy := (float64(y) * float64(viewportHeight) / float64(windowHeight))
	vz := distCameraToProjectionPlane
	return vec3{vx, vy, vz}
}
func intersectRaySphere(origin vec3, direction vec3, sp sphere) solutions {
	r := sp.radius
	co := vec3Sub(origin, sp.center)

	a := vec3Dot(direction, direction)
	b := 2 * vec3Dot(co, direction)
	c := vec3Dot(co, co) - r*r

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return solutions{math.Inf(0), math.Inf(0)}
	}

	t1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b - math.Sqrt(discriminant)) / (2 * a)

	return solutions{t1, t2}
}

func traceRay(origin vec3, direction vec3, tMin float64, tMax float64) color {
	closestT := math.Inf(0)
	var closestSphere *sphere = nil

	for i := 0; i < 3; i++ {
		sp := shapes[i]
		sols := intersectRaySphere(origin, direction, sp)
		t1 := sols.first
		t2 := sols.second
		if t1 > tMin && t1 < tMax && t1 < closestT {
			closestT = t1
			closestSphere = &sp
		}
		if t2 > tMin && t2 < tMax && t2 < closestT {
			closestT = t2
			closestSphere = &sp
		}

	}
	if closestSphere == nil {
		return backgroundColor
	}
	return closestSphere.color

}

// X and Y are canvas coordinates
// (0,0 in middle, -(Width/2), -(Hight/2) in bottom left).
func putPixel(screen *[windowWidth * windowHeight * 4]byte, r byte, g byte, b byte, x int, y int) {
	screenX := (windowWidth / 2) + x
	screenY := (windowHeight / 2) - y - 1
	base := (screenY*windowWidth + screenX) * 4
	screen[base] = r
	screen[base+1] = g
	screen[base+2] = b
	screen[base+3] = 0xFF
	screen[0] = 0xFF
}

func rayTraceFrame(tex *sdl.Texture) {
	var screen = [windowWidth * windowHeight * 4]byte{}

	var origin = vec3{0, 0, 0}

	for x := -(windowWidth / 2); x < (windowWidth / 2); x++ {
		for y := -(windowHeight / 2); y < (windowHeight / 2); y++ {
			direction := canvasToViewport(x, y)
			color := traceRay(origin, direction, 1, math.Inf(0))
			putPixel(&screen, color.r, color.g, color.b, x, y)
		}
	}

	bytes, _, err := tex.Lock(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(windowWidth*windowHeight*4); i++ {
		bytes[i] = screen[i]
	}
	tex.Unlock()
}

func drawFrame(renderer *sdl.Renderer, tex *sdl.Texture) {
	rect := sdl.Rect{X: 0, Y: 0, W: int32(windowWidth), H: int32(windowHeight)}
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
