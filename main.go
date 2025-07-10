package main

import (
	"fmt"
	"math"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 780
)

type Game struct {
	rt    rl.RenderTexture2D
	font  rl.Font
	cards []rl.Texture2D

	// screen adjustment
	sw      float32
	sh      float32
	vw      float32
	vh      float32
	scale   float32
	srcRect rl.Rectangle
	dstRect rl.Rectangle

	// layout
	column  float32
	xCenter float32
	xLeft   float32
	xRight  float32
	yCenter float32
	row     float32
	yTop    float32
	yBottom float32
}

func NewGame() *Game {
	g := &Game{}

	rl.SetTraceLogLevel(rl.LogError)
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowResizable)
	rl.InitWindow(ScreenWidth, ScreenHeight, "GB Game")

	g.rt = rl.LoadRenderTexture(1980, 1080)
	g.font = rl.LoadFontEx("res/fonts/Roboto-Regular.ttf", 40, nil, 0)
	rl.SetTextureFilter(g.font.Texture, rl.FilterBilinear)

	cardFiles, _ := filepath.Glob("res/player/*.png")
	g.cards = make([]rl.Texture2D, len(cardFiles))
	for i, cardFile := range cardFiles {
		img := rl.LoadImage(cardFile)
		g.cards[i] = rl.LoadTextureFromImage(img)
		rl.UnloadImage(img)
	}

	return g
}

func (g *Game) Shutdown() {
	for _, card := range g.cards {
		rl.UnloadTexture(card)
	}
	rl.UnloadFont(g.font)
	rl.UnloadRenderTexture(g.rt)
	rl.CloseWindow()
}

func (g *Game) Update() {
	g.sw = float32(rl.GetScreenWidth())
	g.sh = float32(rl.GetScreenHeight())
	g.vw = float32(g.rt.Texture.Width)
	g.vh = float32(g.rt.Texture.Height)
	g.scale = float32(math.Min(float64(g.sw/g.vw), float64(g.sh/g.vh)))
	g.srcRect = rl.Rectangle{Width: g.vw, Height: -g.vh}
	g.dstRect = rl.Rectangle{
		X:      (g.sw - g.vw*g.scale) / 2,
		Y:      (g.sh - g.vh*g.scale) / 2,
		Width:  g.vw * g.scale,
		Height: g.vh * g.scale,
	}

	g.column = g.vw / 12
	g.xCenter = g.vw / 2
	g.xLeft = g.column
	g.xRight = g.vw - g.column
	g.yCenter = g.vh / 2
	g.row = g.vh / 24
	g.yTop = g.row
	g.yBottom = g.vh - g.yTop
}

func (g *Game) Draw() {
	rl.BeginTextureMode(g.rt)
	rl.ClearBackground(rl.SkyBlue)

	// Call the extracted app method
	g.app()

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

func (g *Game) app() {
	// -- debugging circles
	rl.DrawCircle(int32(g.xLeft), int32(g.yTop), 5, rl.Red)
	rl.DrawCircle(int32(g.xLeft), int32(g.yCenter), 5, rl.Red)
	rl.DrawCircle(int32(g.xLeft), int32(g.yBottom), 5, rl.Red)
	rl.DrawCircle(int32(g.xCenter), int32(g.yTop), 5, rl.Red)
	rl.DrawCircle(int32(g.xCenter), int32(g.yCenter), 5, rl.Red)
	rl.DrawCircle(int32(g.xCenter), int32(g.yBottom), 5, rl.Red)
	rl.DrawCircle(int32(g.xRight), int32(g.yTop), 5, rl.Red)
	rl.DrawCircle(int32(g.xRight), int32(g.yCenter), 5, rl.Red)
	rl.DrawCircle(int32(g.xRight), int32(g.yBottom), 5, rl.Red)

	// -- fps
	fps := fmt.Sprintf("%d fps. Press ESC to exit", rl.GetFPS())
	rl.DrawTextEx(g.font, fps, rl.Vector2{X: g.xLeft, Y: g.yTop}, 40, 0, rl.DarkBlue)

	// -- cards
	yPadding := float32(50)
	for i, card := range g.cards {
		x := g.column + float32(i%10)*g.column
		y := g.yTop + yPadding*2 + float32(i/10)*250.0
		cardRect := rl.Rectangle{X: x, Y: y, Width: float32(card.Width), Height: float32(card.Height)}

		rl.DrawTexturePro(card, rl.Rectangle{Width: float32(card.Width), Height: float32(card.Height)}, cardRect, rl.Vector2{}, 0, rl.White)
		rl.DrawTextEx(g.font, fmt.Sprintf("Card %d", i+1), rl.Vector2{X: x, Y: y - yPadding}, 30, 0, rl.Black)

		if rl.CheckCollisionPointRec(g.GetTransformedMousePos(), cardRect) {
			rl.DrawRectangleLinesEx(cardRect, 2, rl.Red)
			rl.DrawTextEx(g.font, fmt.Sprintf("Card %d", i+1), rl.Vector2{X: x, Y: y - yPadding}, 30, 0, rl.Red)
		}
	}
}

func main() {
	g := NewGame()
	defer g.Shutdown()

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()
	}
}
