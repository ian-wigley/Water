package main

import "github.com/hajimehoshi/ebiten/v2"

type Particle struct {
	position    *Vector2
	velocity    *Vector2
	orientation float64
}

func (particle *Particle) Construct(Position Vector2, Velocity Vector2, Orientation float64) {
	particle.position = &Position
	particle.velocity = &Velocity
	particle.orientation = Orientation
}

func (particle *Particle) Draw(screen *ebiten.Image, image *ebiten.Image) {
	particlePos := &ebiten.DrawImageOptions{}
	particlePos.GeoM.Translate(particle.position.x, particle.position.y)
	screen.DrawImage(image, particlePos)
}
