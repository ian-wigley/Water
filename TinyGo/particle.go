package main

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

func (particle *Particle) Draw(g *Game) {
	g.ctx.Call("drawImage", g.particle, particle.position.x, particle.position.y)
}
