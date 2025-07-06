package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ofabricio/arpg"
)

func main() {

	rl.InitWindow(800, 450, "arpg")
	defer rl.CloseWindow()

	hero := arpg.NewHero()

	rl.SetExitKey(rl.KeyEscape)
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		hero.Update(rl.GetFrameTime())
		hero.Draw()

		rl.EndDrawing()
	}
}
