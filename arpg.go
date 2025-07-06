package arpg

import (
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var DebugEnabled = false

type Game struct {
	Entities []GameObj
}

func (g *Game) Update(dt float32) {
	for _, v := range g.Entities {
		v.Update(dt)
	}
	sort.SliceStable(g.Entities, func(a, b int) bool { return g.Entities[a].Position().Y < g.Entities[b].Position().Y })
}

func (g *Game) Draw() {
	for _, v := range g.Entities {
		v.Draw()
	}
}

type GameObj interface {
	Update(float32)
	Draw()
	Position() rl.Vector2
}

type Timer struct {
	Time      float32
	Completed bool
	elapsed   float32
}

func (t *Timer) Update(dt float32) {
	if t.Completed {
		t.elapsed = 0
	}
	t.elapsed += dt
	t.Completed = t.elapsed >= t.Time
}
