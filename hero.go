package arpg

import (
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ofabricio/arpg/sprite"
)

func NewHero() Hero {
	var h Hero
	h.animationMove = sprite.NewSheet(rl.LoadTexture("assets/warrior_run.png"), 10, 192, 0, 0, 0, 5)
	h.animationIdle = sprite.NewSheet(rl.LoadTexture("assets/warrior_idle.png"), 10, 192, 0, 0, 0, 7)
	h.animationGuard = sprite.NewSheet(rl.LoadTexture("assets/warrior_guard.png"), 10, 192, 0, 0, 0, 5)
	h.animationAttack1 = sprite.NewSheet(rl.LoadTexture("assets/warrior_attack1.png"), 10, 192, 0, 0, 0, 3)
	h.animationAttack2 = sprite.NewSheet(rl.LoadTexture("assets/warrior_attack2.png"), 10, 192, 0, 0, 0, 3)
	h.animation = &h.animationIdle
	h.animationOffset = rl.Vector2{X: 96, Y: 96 + 33} // Half of the sprite size (192x192).
	h.stateIdle = &HeroIdle{}
	h.stateMove = &HeroMove{}
	h.stateAttack = &HeroAttack{}
	h.stateDefend = &HeroGuard{}
	h.state = h.stateIdle
	h.position.X = 100
	h.position.Y = 100
	h.attackDistance = 70
	h.Speed = 100
	return h
}

type Hero struct {
	animationIdle    sprite.Sheet
	animationMove    sprite.Sheet
	animationAttack1 sprite.Sheet
	animationAttack2 sprite.Sheet
	animationGuard   sprite.Sheet
	animation        *sprite.Sheet
	animationOffset  rl.Vector2 // Offset to center the animation in the proper place.

	state       HeroState
	stateIdle   HeroState
	stateMove   HeroState
	stateAttack HeroState
	stateDefend HeroState

	G *Game

	Speed float32

	position rl.Vector2
	target   rl.Vector2

	attackDistance float32

	dirMovmt rl.Vector2 // Movement direction.
	dirMouse rl.Vector2 // Mouse direction.
}

func (h *Hero) attack() {
	for _, entity := range h.G.FindInside(h.position, h.attackDistance) {
		if entity != h {
			if e, ok := entity.(interface{ Hurt() }); ok {
				e.Hurt()
			}
		}
	}
}

func (h *Hero) Update(dt float32) {

	h.dirMouse = rl.Vector2Normalize(rl.Vector2Subtract(rl.GetMousePosition(), h.position))

	next := h.state.Update(h, dt)
	if next != nil {
		if v, ok := next.(interface{ Enter(*Hero) }); ok {
			v.Enter(h)
		}
		h.state = next
	}

	h.animation.Flip = h.dirMovmt.X < 0
	h.animation.Update(dt)
}

func (h *Hero) Draw() {
	p := h.position
	p.X -= h.animationOffset.X
	p.Y -= h.animationOffset.Y
	h.animation.Draw(p)

	if DebugEnabled {
		rl.DrawCircleLines(int32(h.position.X), int32(h.position.Y), h.attackDistance, rl.Red)
		rl.DrawText(reflect.TypeOf(h.state).Elem().Name(), int32(h.position.X), int32(h.position.Y)-100, 20, rl.DarkGray)
	}
}

func (h *Hero) Position() rl.Vector2 {
	return h.position
}

func (h *Hero) wantToMove() bool {
	t := rl.GetMousePosition()
	return rl.IsMouseButtonDown(rl.MouseButtonLeft) && rl.Vector2Distance(h.position, t) >= 15
}

func (h *Hero) wantToAttack() bool {
	return rl.IsMouseButtonDown(rl.MouseButtonRight)
}

func (h *Hero) wantToDefend() bool {
	return rl.IsKeyDown(rl.KeySpace)
}

func (h *Hero) lookAtMouse() {
	h.dirMovmt = h.dirMouse
}

func (h *Hero) stopMoving() {
	h.target = h.position
}

func (h *Hero) setTarget() {
	h.target = rl.GetMousePosition()
	h.dirMovmt = rl.Vector2Subtract(h.target, h.position)
}

type HeroState interface {
	Update(*Hero, float32) HeroState
}

type HeroIdle struct{}

func (s *HeroIdle) Update(h *Hero, dt float32) HeroState {
	h.animation = &h.animationIdle
	if h.wantToMove() {
		return h.stateMove
	}
	if h.wantToAttack() {
		return h.stateAttack
	}
	if h.wantToDefend() {
		return h.stateDefend
	}
	return nil
}

type HeroMove struct{}

func (s *HeroMove) Update(h *Hero, dt float32) HeroState {

	h.animation = &h.animationMove

	if h.wantToAttack() {
		return h.stateAttack
	}

	if h.wantToDefend() {
		return h.stateDefend
	}

	if h.wantToMove() {
		h.setTarget()
	}

	// Move to target.
	h.position = rl.Vector2MoveTowards(h.position, h.target, h.Speed*dt)

	// Reached the target?
	if rl.Vector2Distance(h.position, h.target) <= h.Speed*dt {
		// Snap.
		h.position = h.target
		return h.stateIdle
	}

	return nil
}

type HeroAttack struct{}

func (s *HeroAttack) Enter(h *Hero) {
	h.stopMoving()
	h.lookAtMouse()

	if h.animation == &h.animationAttack2 {
		h.animation = &h.animationAttack1
	} else {
		h.animation = &h.animationAttack2
	}
}

func (s *HeroAttack) Update(h *Hero, dt float32) HeroState {
	if h.animation.Completed {
		h.attack()
		if h.wantToAttack() {
			return s
		}
		return h.stateIdle
	}
	return nil
}

type HeroGuard struct{}

func (s *HeroGuard) Update(h *Hero, dt float32) HeroState {

	h.animation = &h.animationGuard
	h.stopMoving()
	h.lookAtMouse()

	if h.wantToDefend() {
		return nil
	}
	if h.wantToMove() {
		return h.stateMove
	}
	if h.wantToAttack() {
		return h.stateAttack
	}

	return nil
}
