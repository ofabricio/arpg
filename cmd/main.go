package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ofabricio/arpg"
)

func main() {

	rl.InitWindow(1280, 720, "arpg")
	defer rl.CloseWindow()

	var game arpg.Game

	hero := arpg.NewHero()
	hero.G = &game

	enemy1 := arpg.NewEnemy(1280-1280/4, 720/2)
	enemy1.Hero = &hero
	enemy1.AttackSpeed = 1.0

	// enemy2 := arpg.NewEnemy(1280/4, 720/2)
	// enemy2.Hero = &hero
	// enemy2.AttackSpeed = 2.0

	game.Entities = append(game.Entities, &hero, &enemy1) //, &enemy2)

	for range 10 {
		enemy := arpg.NewEnemy(float32(rand.Intn(1280)), float32(rand.Intn(720)))
		enemy.Hero = &hero
		enemy.AttackSpeed = 1.0
		game.Entities = append(game.Entities, &enemy)
	}

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
