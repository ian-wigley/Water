package main

import (
	// "log"
	"slices"
	"syscall/js"
	"time"
)

var (
	rocks     []Rock
	water     *Water
	mousePosX int
)

const (
	screenWidth  = 800
	screenHeight = 480
	waterSurface = 240
	waterDepth   = 340
)

type Game struct {
	ctx          js.Value
	sky          js.Value
	rock         js.Value
	ground       js.Value
	particle     js.Value
	mouseX       float64
	mouseY       float64
	keys         map[string]bool
	mouseButtons map[string]bool
}

func Initialise() {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "gameCanvas")
	g := &Game{
		ctx:          canvas.Call("getContext", "2d"),
		sky:          doc.Call("getElementById", "sky"),
		rock:         doc.Call("getElementById", "rock"),
		ground:       doc.Call("getElementById", "ground"),
		particle:     doc.Call("getElementById", "particle"),
		mouseX:      100,
		mouseY:      100,
		keys:         make(map[string]bool),
		mouseButtons: make(map[string]bool),
	}

	water = new(Water)
	water.Construct(0, 400)
	water.Create()

	// Wait for image assets to load
	for !g.sky.Get("complete").Bool() {
		time.Sleep(10 * time.Millisecond)
	}

	g.SetupMouse()

	var runLoop js.Func
	var lastTime float64

	runLoop = js.FuncOf(func(this js.Value, args []js.Value) any {
		// requestAnimationFrame passes a DOMHighResTimeStamp (float64 ms)
		currentTime := args[0].Float()

		if lastTime == 0 {
			lastTime = currentTime
		}

		// Calculate delta time in seconds
		dt := (currentTime - lastTime) / 1000.0
		lastTime = currentTime

		// Pass dt into Update so movements scale with time, not frames
		g.Update(dt)
		g.Draw()

		js.Global().Call("requestAnimationFrame", runLoop)
		return nil
	})

	defer runLoop.Release()
	js.Global().Call("requestAnimationFrame", runLoop)
	select {}
}

func (g *Game) SetupMouse() {
	doc := js.Global().Get("document")
	win := js.Global()

	mouseDown := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		event.Call("preventDefault")

		button := event.Get("button").Int()

		switch button {
		case 0:
			g.mouseButtons["left"] = true
			if g.mouseX > 40 && g.mouseX < screenWidth-80 {
				newRock := new(Rock)
				newRock.Construct(Vector2{float64(g.mouseX), 100}, Vector2{0, 0})
				newRock.Update(water)
				rocks = append(rocks, *newRock)
			}
		}

		return nil
	})

	// Mouse button released.
	mouseUp := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		button := event.Get("button").Int()

		switch button {
		case 0:
			g.mouseButtons["left"] = false
		}

		return nil
	})

	// Track cursor position.
	mouseMove := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]

		g.mouseX = event.Get("clientX").Float()
		// g.playerX = g.mouseX
		return nil
	})

	// Prevent the right-click context menu.
	contextMenu := js.FuncOf(func(this js.Value, args []js.Value) any {
		args[0].Call("preventDefault")
		return nil
	})

	// Clear held buttons if the browser loses focus.
	blur := js.FuncOf(func(this js.Value, args []js.Value) any {
		for button := range g.mouseButtons {
			g.mouseButtons[button] = false
		}
		return nil
	})

	doc.Call("addEventListener", "mousedown", mouseDown)
	win.Call("addEventListener", "mouseup", mouseUp)
	doc.Call("addEventListener", "mousemove", mouseMove)
	doc.Call("addEventListener", "contextmenu", contextMenu)
	win.Call("addEventListener", "blur", blur)
}

func (g *Game) Update(time float64) {
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
}

func (g *Game) Draw() {
	g.ctx.Call("clearRect", 0, 0, 800, 600)
	g.ctx.Call("drawImage", g.sky, 0, 0)
	g.ctx.Call("drawImage", g.rock, g.mouseX, g.mouseY)

	for _, rocket := range rocks {
		rocket.Draw(g)
	}
	water.Draw(g)
	g.ctx.Call("drawImage", g.ground, -10, 250)
	g.ctx.Call("drawImage", g.ground, 750, 250)
}

// func LoadAssets(fileName string) *ebiten.Image {
// 	var err error
// 	var loadImage *ebiten.Image
// 	loadImage, _, err = ebitenutil.NewImageFromFile(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return loadImage
// }



//func (g *Game) SetUpKeyboard() {

// 	doc := js.Global().Get("document")

// 	keyDown := js.FuncOf(func(this js.Value, args []js.Value) any {
// 		event := args[0]
// 		code := event.Get("code").String()
// 		event.Call("preventDefault")
// 		if event.Get("repeat").Bool() {
// 			return nil
// 		}
// 		g.keys[code] = true
// 		switch code {
// 		// // 	mousePosX, _ = ebiten.CursorPosition()
// 		// // 	dropRock := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
// 		case "KeyA":
// 			if g.playerX > 40 && g.playerX < screenWidth-80 {
// 				// 	if dropRock && mousePosX > 40 && mousePosX < screenWidth - 80 {
// 				// log.Print("Dropping a rock !")
// 				newRock := new(Rock)
// 				newRock.Construct(Vector2{float64(g.playerX), 100}, Vector2{0, 0})
// 				newRock.Update(water)
// 				rocks = append(rocks, *newRock)
// 			}
// 		}

// 		return nil
// 	})

// 	keyUp := js.FuncOf(func(this js.Value, args []js.Value) any {
// 		event := args[0]
// 		code := event.Get("code").String()
// 		event.Call("preventDefault")
// 		g.keys[code] = false

// 		return nil
// 	})

// 	doc.Call("addEventListener", "keydown", keyDown)
// 	doc.Call("addEventListener", "keyup", keyUp)
// }
