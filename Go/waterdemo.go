package main

import (
	"errors"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	sky        *ebiten.Image
	rock       *ebiten.Image
	ground     *ebiten.Image
	particle   *ebiten.Image
	rocks      []Rock
	water      *Water
	mousePosX  int
	terminated = errors.New("terminated")
)

const (
	screenWidth  = 800
	screenHeight = 480
	waterSurface = 240
	waterDepth   = 340
)

func init() {
	sky = LoadAssets("assets/sky.png")
	rock = LoadAssets("assets/rock.png")
	ground = LoadAssets("assets/ground.png")
	particle = LoadAssets("assets/particle.png")
	water = new(Water)
	water.Construct(0, 400)
	water.Create()
}

func LoadAssets(fileName string) *ebiten.Image {
	var err error
	var loadImage *ebiten.Image
	loadImage, _, err = ebitenutil.NewImageFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return loadImage
}

type Game struct{}

func (g *Game) Update() error {

	mousePosX, _ = ebiten.CursorPosition()

	dropRock := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if dropRock && mousePosX > 40 && mousePosX < screenWidth - 80 {
		newRock := new(Rock)
		newRock.Construct(Vector2{float64(mousePosX), 100}, Vector2{0, 0})
		newRock.Update(water)
		rocks = append(rocks, *newRock)
	}

	var indices []int
	for i, rock := range rocks {
		if rock.position.y < waterSurface && rock.position.y+rock.velocity.y >= waterSurface {
			water.Splash(rock.position.x, rock.velocity.y*rock.velocity.y*5)
		}
		rock.Update(water)

		if rock.position.y < waterDepth {
			rock.position.y += 5
		} else {
			indices = append(indices, i)
		}
	}

	water.Update()

	for i := len(indices) - 1; i >= 0; i-- {
		rocks = slices.Delete(rocks, indices[i], indices[i]+1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return terminated
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	skyPos := &ebiten.DrawImageOptions{}
	skyPos.GeoM.Scale(1.0, 1.0)
	screen.DrawImage(sky, skyPos)

	rockPos := &ebiten.DrawImageOptions{}
	rockPos.GeoM.Translate(float64(mousePosX), 100)
	screen.DrawImage(rock, rockPos)
	for _, rocket := range rocks {
		rocket.Draw(screen, rock)
	}
	err := water.Draw(screen)
	if err != nil {
		return
	}

	groundPos := &ebiten.DrawImageOptions{}
	groundPos.GeoM.Translate(-10, 250)
	screen.DrawImage(ground, groundPos)

	groundPosTwo := &ebiten.DrawImageOptions{}
	groundPosTwo.GeoM.Translate(750, 250)
	screen.DrawImage(ground, groundPosTwo)
}

func (g *Game) Layout(int, int) (screenWidth int, screenHeight int) {
	return 800, 480
}
