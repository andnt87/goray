package main

import (
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1280
	screenHeight = 780
)

var cardValues = map[string]int{
	"cardJoker1":   1,
	"cardSpades2":  2,
	"cardSpades3":  3,
	"cardSpades4":  4,
	"cardSpades5":  5,
	"cardSpades6":  6,
	"cardSpades7":  7,
	"cardSpades8":  8,
	"cardSpades9":  9,
	"cardSpades10": 10,
	"cardSpades11": 11,
	"cardSpades12": 12,
	"cardSpades13": 13,
	"cardSpades14": 14,
	"cardHearts12": 12,
	"cardHearts13": 13,
	"cardHearts14": 14,
}

type Texture struct {
	Name    string
	Texture rl.Texture2D
}

type Card struct {
	ID      int
	Name    string
	Value   int
	Texture rl.Texture2D
}

type Game struct {
	RT   rl.RenderTexture2D
	Font rl.Font

	PlayerTextures []Texture
	TableTextures  []Texture
	BackTextures   []Texture

	TableCards  []Card
	PlayerCards []Card

	Scale   float32
	SrcRect rl.Rectangle
	DstRect rl.Rectangle
}

func main() {
	game := NewGame()
	defer game.Shutdown()
	game.Run()
}

func NewGame() *Game {
	loadTextures := func(imagesPath string) []Texture {
		files, _ := filepath.Glob(imagesPath)
		textures := make([]Texture, 0, len(files))
		for _, file := range files {
			name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
			img := rl.LoadImage(file)
			textures = append(textures, Texture{Name: name, Texture: rl.LoadTextureFromImage(img)})
			rl.UnloadImage(img)
		}
		return textures
	}

	rl.SetTraceLogLevel(rl.LogError)
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "GB Game")

	g := &Game{
		RT:             rl.LoadRenderTexture(1920, 1056),
		Font:           rl.LoadFontEx("res/fonts/Roboto-Regular.ttf", 40, nil, 0),
		PlayerTextures: loadTextures("res/cards/player/*.png"),
		TableTextures:  loadTextures("res/cards/table/*.png"),
		BackTextures:   loadTextures("res/cards/back/*.png"),
	}
	rl.SetTextureFilter(g.Font.Texture, rl.FilterBilinear)

	createDeck := func(textures []Texture, count int) []Card {
		deck := make([]Card, 0, count)
		for i := 0; i < count; i++ {
			rn := rand.Intn(len(textures))
			name := textures[rn].Name
			card := Card{
				ID:      i,
				Name:    name,
				Value:   cardValues[name],
				Texture: textures[rn].Texture,
			}
			deck = append(deck, card)
		}
		return deck
	}
	g.TableCards = createDeck(g.TableTextures, 24)
	g.PlayerCards = createDeck(g.PlayerTextures, 48)

	return g
}

func (g *Game) Shutdown() {
	for _, t := range g.PlayerTextures {
		rl.UnloadTexture(t.Texture)
	}
	for _, t := range g.TableTextures {
		rl.UnloadTexture(t.Texture)
	}
	for _, t := range g.BackTextures {
		rl.UnloadTexture(t.Texture)
	}
	rl.UnloadFont(g.Font)
	rl.UnloadRenderTexture(g.RT)
	rl.CloseWindow()
}

func (g *Game) Run() {
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		// update
		sw := float32(rl.GetScreenWidth())
		sh := float32(rl.GetScreenHeight())
		vw := float32(g.RT.Texture.Width)
		vh := float32(g.RT.Texture.Height)

		g.Scale = float32(math.Min(float64(sw/vw), float64(sh/vh)))
		g.SrcRect = rl.Rectangle{Width: vw, Height: -vh}
		g.DstRect = rl.Rectangle{
			X:      (sw - vw*g.Scale) / 2,
			Y:      (sh - vh*g.Scale) / 2,
			Width:  vw * g.Scale,
			Height: vh * g.Scale,
		}

		// draw virtual texture
		rl.BeginTextureMode(g.RT)
		rl.ClearBackground(rl.SkyBlue)
		g.TheGame()
		rl.EndTextureMode()

		// draw to screen
		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)
		rl.DrawTexturePro(g.RT.Texture, g.SrcRect, g.DstRect, rl.Vector2{}, 0, rl.White)
		rl.EndDrawing()
	}
}

func (g *Game) TheGame() {
	vw := float32(g.RT.Texture.Width)
	vh := float32(g.RT.Texture.Height)
	column := vw / 12
	xLeft := column
	xCenter := vw / 2
	xRight := vw - column
	row := vh / 24
	yTop := row
	yCenter := vh / 2
	yBottom := vh - yTop
	mousePos := rl.GetMousePosition()
	mousePos.X = (mousePos.X - g.DstRect.X) / g.Scale
	mousePos.Y = (mousePos.Y - g.DstRect.Y) / g.Scale

	// -- The Game
	debugCircles := func(xLeft float32, yTop float32, yCenter float32, yBottom float32, xCenter float32, xRight float32, radius float32) {
		if radius <= 0 {
			return
		}
		rl.DrawCircle(int32(xLeft), int32(yTop), radius, rl.Red)
		rl.DrawCircle(int32(xLeft), int32(yCenter), radius, rl.Red)
		rl.DrawCircle(int32(xLeft), int32(yBottom), radius, rl.Red)
		rl.DrawCircle(int32(xCenter), int32(yTop), radius, rl.Red)
		rl.DrawCircle(int32(xCenter), int32(yCenter), radius, rl.Red)
		rl.DrawCircle(int32(xCenter), int32(yBottom), radius, rl.Red)
		rl.DrawCircle(int32(xRight), int32(yTop), radius, rl.Red)
		rl.DrawCircle(int32(xRight), int32(yCenter), radius, rl.Red)
		rl.DrawCircle(int32(xRight), int32(yBottom), radius, rl.Red)
	}

	debugCircles(xLeft, yTop, yCenter, yBottom, xCenter, xRight, 0)

	fpsText := fmt.Sprintf("%d fps, Screen: %.0fx%.0f, Viewport: %.0fx%.0f, Scale: %.2f", rl.GetFPS(),
		float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()), vw, vh, g.Scale)
	rl.DrawTextEx(g.Font, fpsText, rl.Vector2{X: xLeft, Y: yTop}, 40, 0, rl.DarkBlue)

	footerText := fmt.Sprintf("Table Cards: %d, Player Cards: %d", len(g.TableCards), len(g.PlayerCards))
	rl.DrawTextEx(g.Font, footerText, rl.Vector2{X: xLeft, Y: yBottom - 40}, 40, 0, rl.DarkBlue)

	yPadding := float32(50)

	for i, card := range g.TableCards {
		x := column + float32(i%10)*column
		y := yTop + yPadding*2.5 + float32(i/10)*250.0
		if y+250 > yBottom {
			break
		}
		cardRect := rl.Rectangle{X: x, Y: y, Width: float32(card.Texture.Width), Height: float32(card.Texture.Height)}

		rl.DrawTexture(card.Texture, int32(x), int32(y), rl.White)
		rl.DrawTextEx(g.Font, fmt.Sprintf("Card %d", i+1), rl.Vector2{X: x, Y: y - yPadding}, 30, 0, rl.Black)

		if rl.CheckCollisionPointRec(mousePos, cardRect) {
			rl.DrawRectangleLinesEx(cardRect, 2, rl.Red)
			rl.DrawTextEx(g.Font, fmt.Sprintf("Card %d", i+1), rl.Vector2{X: x, Y: y - yPadding}, 30, 0, rl.Red)
		}
	}
}
