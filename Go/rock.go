package main

import "github.com/hajimehoshi/ebiten/v2"

type Rock struct {
	position *Vector2
	velocity *Vector2
	gravity  *Vector2
}

func (rock *Rock) Construct(Position Vector2, Velocity Vector2) {
	rock.position = &Position
	rock.velocity = &Velocity
	rock.gravity = &Vector2{0, 0.98}
}

func (rock *Rock) Update(water *Water) {
	if rock.position.y > water.GetHeight(rock.position.x) {
		rock.velocity.x *= 0.84
		rock.velocity.y *= 0.84
	}
	rock.position.x += rock.velocity.x
	rock.position.y += rock.velocity.y
	rock.velocity.x += rock.gravity.x
	rock.velocity.y += rock.gravity.y
}

func (rock *Rock) Draw(screen *ebiten.Image, image *ebiten.Image) {
	rockPos := &ebiten.DrawImageOptions{}
	rockPos.GeoM.Translate(rock.position.x, rock.position.y)
	screen.DrawImage(image, rockPos)
}
