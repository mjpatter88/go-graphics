package main

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

func computeLighting(position vec3, normal vec3) float64 {
	intensity := 0.0
	for i := 0; i < len(lights); i++ {
		light := lights[i]
		if light.lightType == "ambinent" {
			intensity += light.intensity
		} else {
			var lightDirection vec3
			if light.lightType == "point" {
				lightDirection = vec3Sub(light.position, position)
			} else {
				lightDirection = light.direction
			}

			normalDotLight := vec3Dot(normal, lightDirection)
			if normalDotLight > 0 {
				intensity += light.intensity * normalDotLight / (vec3Len(normal) * vec3Len(lightDirection))
			}
		}
	}

	return intensity
}
