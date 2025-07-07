package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ofabricio/arpg"
)

func main() {

	rl.InitWindow(1280, 720, "arpg")
	defer rl.CloseWindow()

	var game arpg.Game

	hero := arpg.NewHero()
	hero.G = &game

	enemy := arpg.NewEnemy()
	enemy.Hero = &hero

	game.Entities = append(game.Entities, &hero, &enemy)

	rl.SetExitKey(rl.KeyEscape)
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("arpg", 190, 200, 20, rl.LightGray)
		rl.DrawText("[left click] move", 190, 225, 20, rl.LightGray)
		rl.DrawText("[right click] attack", 190, 250, 20, rl.LightGray)
		rl.DrawText("[F1] debug", 190, 275, 20, rl.LightGray)

		if rl.IsKeyPressed(rl.KeyF1) {
			arpg.DebugEnabled = !arpg.DebugEnabled
		}

		dt := rl.GetFrameTime()
		game.Update(dt)

		game.Draw()

		rl.EndDrawing()
	}
}
