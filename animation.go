package arpg

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type HurtAnimation struct {
	timer Timer
}

func (a *HurtAnimation) Value() rl.Color {
	return rl.ColorLerp(rl.White, rl.Red, PingPong(a.timer.Elapsed*8))
}

func (a *HurtAnimation) Update(dt float32) {
	a.timer.Update(dt)
}

func (a *HurtAnimation) Play() {
	a.timer.Play(1.0 / 8)
}
