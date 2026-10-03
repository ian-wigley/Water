package main

type Vector2 struct {
	x float64
	y float64
}

func (vec2 *Vector2) Construct(x float64, y float64) {
	vec2.x = x
	vec2.y = y
}

func (vec2 *Vector2) Zero() {
	vec2.x = 0
	vec2.y = 0
}
