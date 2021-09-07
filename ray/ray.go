package main

import (
	"math"
)

type vec3 struct {
	x float64
	y float64
	z float64
}

func vec3Scale(vec vec3, factor float64) vec3 {
	return vec3{
		x: vec.x * factor,
		y: vec.y * factor,
		z: vec.z * factor,
	}
}

func vec3Len(vec vec3) float64 {
	return math.Sqrt((vec.x*vec.x + vec.y*vec.y + vec.z*vec.z))
}

func vec3Add(first vec3, second vec3) vec3 {
	return vec3{first.x + second.x, first.y + second.y, first.z + second.z}
}

func vec3Sub(first vec3, second vec3) vec3 {
	return vec3{first.x - second.x, first.y - second.y, first.z - second.z}
}

func vec3Dot(first vec3, second vec3) float64 {
	return (first.x * second.x) + (first.y * second.y) + (first.z * second.z)
}

func vec3Neg(vec vec3) vec3 {
	return vec3{-vec.x, -vec.y, -vec.z}
}

func directionMatrixMultiplication(mat [3][3]float64, dir vec3) vec3 {
	var result vec3

	result.x = (mat[0][0] * dir.x) + (mat[0][1] * dir.y) + (mat[0][2] * dir.z)
	result.y = (mat[1][0] * dir.x) + (mat[1][1] * dir.y) + (mat[1][2] * dir.z)
	result.z = (mat[2][0] * dir.x) + (mat[2][1] * dir.y) + (mat[2][2] * dir.z)

	return result
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

func sumColors(col1 color, col2 color) color {
	newR := float64(col1.r) + float64(col2.r)
	newG := float64(col1.g) + float64(col2.g)
	newB := float64(col1.b) + float64(col2.b)
	return color{
		r: byte(math.Min(255, newR)),
		g: byte(math.Min(255, newG)),
		b: byte(math.Min(255, newB)),
	}
}

type color struct {
	r byte
	g byte
	b byte
}

type sphere struct {
	center     vec3
	radius     float64
	color      color
	specular   int     // -1 represents matte object
	reflective float64 // 0 = not reflective, 1 = completely reflective
}

type solutions struct {
	first  float64
	second float64
}

// Cap rendering at roughly 60 fps
const msPerFrame = 16

const windowWidth int = 1000
const windowHeight int = 1000
const viewportWidth int = 1
const viewportHeight int = 1
const distCameraToProjectionPlane float64 = 1

var backgroundColor = color{0x00, 0x00, 0x00}

var sphere1 = sphere{vec3{0, -1, 3}, 1, color{255, 0, 0}, 500, 0.2}
var sphere2 = sphere{vec3{2, 0, 4}, 1, color{0, 0, 255}, 500, 0.3}
var sphere3 = sphere{vec3{-2, 0, 4}, 1, color{0, 255, 0}, 10, 0.4}
var sphere4 = sphere{vec3{0, -5001, 0}, 5000, color{255, 255, 0}, 1000, 0.5}

var shapes = [...]sphere{sphere1, sphere2, sphere3, sphere4}

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

func closestIntersection(origin vec3, direction vec3, tMin float64, tMax float64) (*sphere, float64) {
	closestT := math.Inf(0)
	var closestSphere *sphere = nil

	for i := 0; i < len(shapes); i++ {
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
	return closestSphere, closestT
}

func calculateReflectedRay(ray vec3, normal vec3) vec3 {
	dotProduct := vec3Dot(normal, ray)
	scaledNormal := vec3Scale(normal, 2*dotProduct)
	return vec3Sub(scaledNormal, ray)
}

func traceRay(origin vec3, direction vec3, tMin float64, tMax float64, recursionDepth int) color {
	closestSphere, closestT := closestIntersection(origin, direction, tMin, tMax)

	if closestSphere == nil {
		return backgroundColor
	}
	position := vec3Add(origin, vec3Scale(direction, closestT))
	normal := vec3Sub(position, closestSphere.center)
	normal = vec3Scale(normal, 1/vec3Len(normal))
	localColor := scaleColor(closestSphere.color, computeLighting(position, normal, vec3Neg(direction), closestSphere.specular))

	// Handle reflections
	reflectiveness := closestSphere.reflective
	if reflectiveness <= 0 || recursionDepth <= 0 {
		return localColor
	}

	reflectedRay := calculateReflectedRay(vec3Neg(direction), normal)
	reflectedColor := traceRay(position, reflectedRay, 0.001, math.Inf(0), recursionDepth-1)

	localColorContribution := scaleColor(localColor, 1-reflectiveness)
	reflectedColorContribution := scaleColor(reflectedColor, reflectiveness)
	return sumColors(localColorContribution, reflectedColorContribution)
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

func rayTraceFrame(screen *[windowWidth * windowHeight * 4]byte, recursionDepth int) {

	// var origin = vec3{0, 0, 0}
	var cameraPosition = vec3{3, 0, 1}
	var cameraOrientation = [3][3]float64{
		{0.7071, 0, -0.7071},
		{0, 1, 0},
		{0.7071, 0, 0.7071},
	}

	for x := -(windowWidth / 2); x < (windowWidth / 2); x++ {
		for y := -(windowHeight / 2); y < (windowHeight / 2); y++ {
			direction := directionMatrixMultiplication(cameraOrientation, canvasToViewport(x, y))
			color := traceRay(cameraPosition, direction, 1, math.Inf(0), recursionDepth)
			putPixel(screen, color, x, y)
		}
	}

}
