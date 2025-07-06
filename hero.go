package arpg

import (
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
	h.Speed = 100
	h.stateIdle = &HeroIdle{}
	h.stateMove = &HeroMove{}
	h.stateAttack = &HeroAttack{}
	h.stateDefend = &HeroGuard{}
	h.state = h.stateIdle
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

	Speed float32

	position rl.Vector2
	target   rl.Vector2

	dirMovmt rl.Vector2 // Movement direction.
	dirMouse rl.Vector2 // Mouse direction.
}

func (h *Hero) Update(dt float32) {

	h.dirMouse = rl.Vector2Normalize(rl.Vector2Subtract(rl.GetMousePosition(), h.position))

	next := h.state.Update(h, dt)
	h.state = next
	h.state.Update(h, dt)

	h.animation.Flip = h.dirMovmt.X < 0
	h.animation.Update(dt)
}

func (h *Hero) Draw() {
	p := h.position
	p.X -= h.animationOffset.X
	p.Y -= h.animationOffset.Y
	h.animation.Draw(p)
}

func (h *Hero) wantToMove() bool {
	t := rl.GetMousePosition()
	return rl.IsMouseButtonDown(rl.MouseButtonLeft) && rl.Vector2Distance(h.position, t) > 15
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
	return s
}

type HeroMove struct{}

func (s *HeroMove) Update(h *Hero, dt float32) HeroState {

	h.animation = &h.animationMove

	if h.wantToMove() {
		// Set the new target.
		h.target = rl.GetMousePosition()
		h.dirMovmt = rl.Vector2Subtract(h.target, h.position)
	}

	// Move to target.
	h.position = rl.Vector2MoveTowards(h.position, h.target, h.Speed*dt)

	// Reached the target? Idle.
	if rl.Vector2Distance(h.position, h.target) <= h.Speed*dt {
		h.position = h.target
		return h.stateIdle
	}

	if h.wantToAttack() {
		// Stop moving.
		h.target = h.position
		return h.stateAttack
	}

	if h.wantToDefend() {
		// Stop moving.
		h.target = h.position
		return h.stateDefend
	}

	return s
}

type HeroAttack struct {
	toggle bool
}

func (s *HeroAttack) Update(h *Hero, dt float32) HeroState {
	if s.toggle {
		h.animation = &h.animationAttack2
	} else {
		h.animation = &h.animationAttack1
	}
	if h.animation.Complete() {
		s.toggle = !s.toggle
		h.lookAtMouse()
		if h.wantToAttack() {
			return s
		}
		return h.stateIdle
	}
	return s
}

type HeroGuard struct{}

func (s *HeroGuard) Update(h *Hero, dt float32) HeroState {
	h.animation = &h.animationGuard
	if h.wantToDefend() {
		if h.animation.Complete() {
			h.lookAtMouse()
		}
		return s
	}
	if h.wantToMove() {
		return h.stateMove
	}
	if h.wantToAttack() {
		return h.stateAttack
	}
	return s
}
