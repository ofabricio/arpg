package arpg

import (
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var DebugEnabled = false

type Game struct {
	Entities []GameObj
}

func (g *Game) FindInside(position rl.Vector2, distance float32) []GameObj {
	var found []GameObj
	for _, entity := range g.Entities {
		if rl.Vector2Distance(position, entity.Position()) <= distance {
			found = append(found, entity)
		}
	}
	return found
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
	Duration float32
	Elapsed  float32
}

func (a *Timer) Play(secs float32) {
	if a.Completed() {
		a.Restart(secs)
	}
}

func (a *Timer) Restart(secs float32) {
	a.Elapsed = 0
	a.Duration = secs
}

func (a *Timer) Update(dt float32) {
	a.Elapsed += dt
	if a.Completed() {
		a.Elapsed = a.Duration
	}
}

func (a *Timer) Completed() bool {
	return a.Elapsed >= a.Duration
}
