package main

import "math"

type light struct {
	lightType string
	intensity float64
	position  vec3
	direction vec3
}

var light1 = light{
	lightType: "ambient",
	intensity: 0.2,
	position:  vec3{0, 0, 0}, // Ambient light doesn't have a position, so set to origin.
	direction: vec3{0, 0, 0}, // Ambient light doesn't have a direction, so set to origin.
}
var light2 = light{
	lightType: "point",
	intensity: 0.6,
	position:  vec3{2, 1, 0},
	direction: vec3{0, 0, 0}, // Point light doesn't have a direction, so set to origin.
}
var light3 = light{
	lightType: "directional",
	intensity: 0.2,
	position:  vec3{0, 0, 0}, // Directional light doesn't have a position, so set to origin.
	direction: vec3{1, 4, 4},
}

var lights = [...]light{light1, light2, light3}

// v is the vector from the point to the camera.
func computeLighting(position vec3, normal vec3, v vec3, s int) float64 {
	intensity := 0.0
	for i := 0; i < len(lights); i++ {
		light := lights[i]
		if light.lightType == "ambient" {
			intensity += light.intensity
		} else {
			var lightDirection vec3
			var tMax float64
			if light.lightType == "point" {
				lightDirection = vec3Sub(light.position, position)
				tMax = 1
			} else {
				lightDirection = light.direction
				tMax = math.Inf(0)
			}

			// Check if this point is in a shadow relative to this light source
			shadow_sphere, _ := closestIntersection(position, lightDirection, 0.001, tMax)
			if shadow_sphere != nil {
				continue
			}

			// Diffuse
			normalDotLight := vec3Dot(normal, lightDirection)
			if normalDotLight > 0 {
				intensity += light.intensity * normalDotLight / (vec3Len(normal) * vec3Len(lightDirection))
			}

			// Specular
			if s != -1 {
				var reflectionDirection vec3
				scale := 2 * normalDotLight
				reflectionDirection = vec3Sub(vec3Scale(normal, scale), lightDirection)
				reflectionDotV := vec3Dot(reflectionDirection, v)
				if reflectionDotV > 0 {
					cosVal := reflectionDotV / (vec3Len(reflectionDirection) * vec3Len(v))
					intensity += light.intensity * math.Pow(cosVal, float64(s))
				}
			}
		}
	}

	return intensity
}
