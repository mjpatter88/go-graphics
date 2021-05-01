package main

import (
	"math"
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

func rayTraceFrame(screen *[windowWidth * windowHeight * 4]byte) {

	var origin = vec3{0, 0, 0}

	for x := -(windowWidth / 2); x < (windowWidth / 2); x++ {
		for y := -(windowHeight / 2); y < (windowHeight / 2); y++ {
			direction := canvasToViewport(x, y)
			color := traceRay(origin, direction, 1, math.Inf(0))
			putPixel(screen, color, x, y)
		}
	}

}
