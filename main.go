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
	rt   rl.RenderTexture2D
	font rl.Font

	playerTextures []Texture
	tableTextures  []Texture
	backTextures   []Texture

	tableCards  []Card
	playerCards []Card

	scale   float32
	srcRect rl.Rectangle
	dstRect rl.Rectangle
}

func NewGame() *Game {
	rl.SetTraceLogLevel(rl.LogError)
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "GB Game")

	g := &Game{
		rt:             rl.LoadRenderTexture(1920, 1056),
		font:           rl.LoadFontEx("res/fonts/Roboto-Regular.ttf", 40, nil, 0),
		playerTextures: loadTextures("res/cards/player/*.png"),
		tableTextures:  loadTextures("res/cards/table/*.png"),
		backTextures:   loadTextures("res/cards/back/*.png"),
	}
	rl.SetTextureFilter(g.font.Texture, rl.FilterBilinear)

	g.tableCards = createDeck(g.tableTextures, 24)
	g.playerCards = createDeck(g.playerTextures, 48)

	return g
}

func (g *Game) Shutdown() {
	for _, t := range g.playerTextures {
		rl.UnloadTexture(t.Texture)
	}
	for _, t := range g.tableTextures {
		rl.UnloadTexture(t.Texture)
	}
	for _, t := range g.backTextures {
		rl.UnloadTexture(t.Texture)
	}
	rl.UnloadFont(g.font)
	rl.UnloadRenderTexture(g.rt)
	rl.CloseWindow()
}

func (g *Game) Run() {
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()
	}
}

func (g *Game) Update() {
	sw := float32(rl.GetScreenWidth())
	sh := float32(rl.GetScreenHeight())
	vw := float32(g.rt.Texture.Width)
	vh := float32(g.rt.Texture.Height)

	g.scale = float32(math.Min(float64(sw/vw), float64(sh/vh)))
	g.srcRect = rl.Rectangle{Width: vw, Height: -vh}
	g.dstRect = rl.Rectangle{
		X:      (sw - vw*g.scale) / 2,
		Y:      (sh - vh*g.scale) / 2,
		Width:  vw * g.scale,
		Height: vh * g.scale,
	}
}

func (g *Game) Draw() {
	rl.BeginTextureMode(g.rt)
	rl.ClearBackground(rl.SkyBlue)
	g.TheGame()
	rl.EndTextureMode()

	rl.BeginDrawing()
	rl.ClearBackground(rl.SkyBlue)
	rl.DrawTexturePro(g.rt.Texture, g.srcRect, g.dstRect, rl.Vector2{}, 0, rl.White)
	rl.EndDrawing()
}

func (g *Game) GetTransformedMousePos() rl.Vector2 {
	mousePos := rl.GetMousePosition()
	mousePos.X = (mousePos.X - g.dstRect.X) / g.scale
	mousePos.Y = (mousePos.Y - g.dstRect.Y) / g.scale
	return mousePos
}

func (g *Game) TheGame() {
	vw := float32(g.rt.Texture.Width)
	vh := float32(g.rt.Texture.Height)
	column := vw / 12
	xLeft := column
	xCenter := vw / 2
	xRight := vw - column
	row := vh / 24
	yTop := row
	yCenter := vh / 2
	yBottom := vh - yTop

	debugCircles(xLeft, yTop, yCenter, yBottom, xCenter, xRight, 0)

	fpsText := fmt.Sprintf("%d fps, Screen: %.0fx%.0f, Viewport: %.0fx%.0f, Scale: %.2f", rl.GetFPS(),
		float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()), vw, vh, g.scale)
	rl.DrawTextEx(g.font, fpsText, rl.Vector2{X: xLeft, Y: yTop}, 40, 0, rl.DarkBlue)

	footerText := fmt.Sprintf("Table Cards: %d, Player Cards: %d", len(g.tableCards), len(g.playerCards))
	rl.DrawTextEx(g.font, footerText, rl.Vector2{X: xLeft, Y: yBottom - 40}, 40, 0, rl.DarkBlue)

	yPadding := float32(50)
	mousePos := g.GetTransformedMousePos()

	for i, card := range g.tableCards {
		x := column + float32(i%10)*column
		y := yTop + yPadding*2.5 + float32(i/10)*250.0
		if y+250 > yBottom {
			break
		}
		cardRect := rl.Rectangle{X: x, Y: y, Width: float32(card.Texture.Width), Height: float32(card.Texture.Height)}

		rl.DrawTexture(card.Texture, int32(x), int32(y), rl.White)
		rl.DrawTextEx(g.font, fmt.Sprintf("Card %d", i+1), rl.Vector2{X: x, Y: y - yPadding}, 30, 0, rl.Black)

		if rl.CheckCollisionPointRec(mousePos, cardRect) {
			rl.DrawRectangleLinesEx(cardRect, 2, rl.Red)
			rl.DrawTextEx(g.font, fmt.Sprintf("Card %d", i+1), rl.Vector2{X: x, Y: y - yPadding}, 30, 0, rl.Red)
		}
	}
}

func loadTextures(imagesPath string) []Texture {
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

func createDeck(textures []Texture, count int) []Card {
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

func debugCircles(xLeft float32, yTop float32, yCenter float32, yBottom float32, xCenter float32, xRight float32, radius float32) {
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

func main() {
	game := NewGame()
	defer game.Shutdown()
	game.Run()
}
