package main

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Water struct {
	x                  float64
	y                  float64
	vertices           []Vector2
	waterValues        []WaterValues
	quantityOfVertices int
	pb                 *PrimitiveBatch
	particles          []Particle
}

func (water *Water) Construct(x float64, y float64) {
	water.x = x
	water.y = y
	water.quantityOfVertices = 201
	water.pb = new(PrimitiveBatch)
	water.pb.Construct()
}

func (water *Water) Create() {
	spacing := 0.0
	for i := 0; i < water.quantityOfVertices; i++ {
		spacing += 2.0
		waterVal := new(WaterValues)
		waterVal.Construct(waterSurface, waterDepth, 0)
		water.waterValues = append(water.waterValues, *waterVal)
		point := new(Vector2)
		point.Construct(water.x+spacing+10, water.y)
		water.vertices = append(water.vertices, *point)
	}
}

func (water *Water) GetHeight(x float64) float64 {
	if x < 0 || x > screenWidth {
		return waterSurface
	}
	index := (int)(x / water.GetScale())
	return water.waterValues[index].height
}

func (water *Water) UpdateParticle(particle Particle) {
	var gravity = 0.3
	particle.velocity.y += gravity
	particle.position.x += particle.velocity.x
	particle.position.y += particle.velocity.y
	particle.orientation = water.GetAngle(*particle.velocity)
}

func (water *Water) Clamp(value float64, minValue int, maxValue int) float64 {
	if value > float64(maxValue) {
		return float64(maxValue)
	} else if value < float64(minValue) {
		return float64(minValue)
	} else {
		return value
	}
}

func (water *Water) Splash(xPosition float64, speed float64) {
	var index = int(water.Clamp(xPosition/water.GetScale(), 0, len(water.waterValues)-1))
	for i := max(0, index); i < min(len(water.waterValues)-1, index+1); i++ {
		water.waterValues[index].speed = speed
	}
	water.CreateSplashParticles(xPosition, int(speed))
}

func (water *Water) GetScale() float64 {
	numCols := len(water.waterValues) - 1
	result := screenWidth / numCols
	return float64(result)
}

func (water *Water) CreateSplashParticles(xPosition float64, speed int) {
	var y = water.GetHeight(xPosition)
	if speed > 0 {
		for i := 0; i < speed/8; i++ {
			temp := water.GetRandomVector2(40)
			pos := new(Vector2)
			pos.Construct(xPosition+temp.x, y+temp.y)
			s1 := water.GetRandomFloat(-150.0, -30.0)
			s2 := water.GetRandomFloat(0, 0.5*math.Sqrt(float64(speed)))
			vel := water.FromPolar(water.ToRadians(s1), s2)
			water.CreateParticle(*pos, vel)
		}
	}
}

func (water *Water) ToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func (water *Water) CreateParticle(pos Vector2, velocity Vector2) {
	p := new(Particle)
	p.Construct(pos, velocity, 0)
	water.particles = append(water.particles, *p)
}

func (water *Water) FromPolar(angle float64, magnitude float64) Vector2 {
	v := new(Vector2)
	s1 := (math.Cos(angle)) * magnitude
	s2 := (math.Sin(angle)) * magnitude
	v.Construct(s1, s2)
	return *v
}

func (water *Water) GetRandomFloat(min float64, max float64) float64 {
	return (float64)(rand.Float32())*(max-min) + min
}

func (water *Water) GetRandomVector2(maxLength float64) Vector2 {
	s1 := water.GetRandomFloat(-math.Pi, math.Pi)
	s2 := water.GetRandomFloat(0, maxLength)
	return water.FromPolar(s1, s2)
}

func (water *Water) GetAngle(vector Vector2) float64 {
	return math.Atan2(vector.y, vector.x)
}

func (water *Water) createDeltaSlices() []float64 {
	deltaSlice := make([]float64, 0, len(water.waterValues))
	for i := 0; i < len(water.waterValues); i++ {
		deltaSlice = append(deltaSlice, 0)
	}
	return deltaSlice
}

func (water *Water) Update() {

	for i := 0; i < len(water.waterValues); i++ {
		water.waterValues[i].Update(0.025, 0.025)
	}

	lDeltas := water.createDeltaSlices()
	rDeltas := water.createDeltaSlices()

	spread := 0.25
	// do some passes where columns pull on their neighbours
	for j := 0; j < 8; j++ {
		for i := 0; i < len(water.waterValues); i++ {
			if i > 0 {
				lDeltas[i] = spread * (water.waterValues[i].height - water.waterValues[i-1].height)
				water.waterValues[i-1].speed += lDeltas[i]
			}
			if i < len(water.waterValues)-1 {
				rDeltas[i] = spread * (water.waterValues[i].height - water.waterValues[i+1].height)
				water.waterValues[i+1].speed += rDeltas[i]
			}
		}

		for i := 0; i < len(water.waterValues); i++ {
			if i > 0 {
				water.waterValues[i-1].height += lDeltas[i]
			}
			if i < len(water.waterValues)-1 {
				water.waterValues[i+1].height += rDeltas[i]
			}
		}
	}

	for _, particle := range water.particles {
		water.UpdateParticle(particle)
	}
}

func (water *Water) Draw(screen *ebiten.Image) error {
	water.pb.Begin()
	bottom := float64(screenHeight)
	scale := water.GetScale()

	for i := 1; i < water.quantityOfVertices; i++ {
		p1 := new(Vector2)
		p1.Construct(float64(i-1)*scale, water.waterValues[i-1].height)

		p2 := new(Vector2)
		p2.Construct(float64(i)*scale, water.waterValues[i].height)

		p3 := new(Vector2)
		p3.Construct(p2.x, bottom)

		p4 := new(Vector2)
		p4.Construct(p1.x, bottom)

		water.pb.AddVertex(*p1, screen)
		water.pb.AddVertex(*p2, screen)
		water.pb.AddVertex(*p3, screen)

		water.pb.AddVertex(*p1, screen)
		water.pb.AddVertex(*p3, screen)
		water.pb.AddVertex(*p4, screen)
	}

	water.pb.End(screen)

	for _, party := range water.particles {
		party.Draw(screen, particle)
	}

	return nil
}
