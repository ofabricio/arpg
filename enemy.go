package arpg

import (
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ofabricio/arpg/sprite"
)

func NewEnemy() Enemy {
	var e Enemy
	e.animationMove = sprite.NewSheet(rl.LoadTexture("assets/warrior_run.png"), 10, 192, 0, 0, 0, 5)
	e.animationIdle = sprite.NewSheet(rl.LoadTexture("assets/warrior_idle.png"), 10, 192, 0, 0, 0, 7)
	e.animationGuard = sprite.NewSheet(rl.LoadTexture("assets/warrior_guard.png"), 10, 192, 0, 0, 0, 5)
	e.animationAttack1 = sprite.NewSheet(rl.LoadTexture("assets/warrior_attack1.png"), 10, 192, 0, 0, 0, 3)
	e.animationAttack2 = sprite.NewSheet(rl.LoadTexture("assets/warrior_attack2.png"), 10, 192, 0, 0, 0, 3)
	e.animation = &e.animationIdle
	e.animationOffset = rl.Vector2{X: 96, Y: 96 + 33} // Half of the sprite size (192x192).
	e.Speed = 90
	e.stateIdle = &EnemyIdle{}
	e.stateChase = &EnemyChase{}
	e.state = e.stateIdle
	e.position.X = 1280 / 2
	e.position.Y = 720 / 2
	e.attackDistance = 192 / 4
	e.minChaseDistance = e.attackDistance * 1.35
	e.awarenessDistance = e.attackDistance * 7
	return e
}

type Enemy struct {
	animationIdle    sprite.Sheet
	animationMove    sprite.Sheet
	animationAttack1 sprite.Sheet
	animationAttack2 sprite.Sheet
	animationGuard   sprite.Sheet
	animation        *sprite.Sheet
	animationOffset  rl.Vector2 // Offset to center the animation in the proper place.

	state      EnemyState
	stateIdle  EnemyState
	stateChase EnemyState

	Hero  *Hero
	Speed float32

	position          rl.Vector2
	attackDistance    float32 // Distance the enemy can attack the hero.
	minChaseDistance  float32 // Minimum distance to start chasing the hero.
	awarenessDistance float32 // Distance the enemy can see the hero.
}

func (e *Enemy) Update(dt float32) {

	next := e.state.Update(e, dt)
	e.state = next
	e.state.Update(e, dt)

	e.animation.Flip = rl.Vector2Subtract(e.Hero.Position(), e.position).X < 0
	e.animation.Update(dt)
}

func (e *Enemy) Draw() {
	p := e.position
	p.X -= e.animationOffset.X
	p.Y -= e.animationOffset.Y
	e.animation.Draw(p)

	if DebugEnabled {
		rl.DrawCircleLines(int32(e.position.X), int32(e.position.Y), e.attackDistance, rl.Red)
		rl.DrawCircleLines(int32(e.position.X), int32(e.position.Y), e.minChaseDistance, rl.Orange)
		rl.DrawCircleLines(int32(e.position.X), int32(e.position.Y), e.awarenessDistance, rl.Gray)
		rl.DrawText(reflect.TypeOf(e.state).Elem().Name(), int32(e.position.X), int32(e.position.Y)-100, 20, rl.DarkGray)
	}
}

func (e *Enemy) ZIndex() float32 {
	return e.position.Y
}

func (e *Enemy) Position() rl.Vector2 {
	return e.position
}

func (e *Enemy) inAwarenessRange() bool {
	return rl.Vector2Distance(e.position, e.Hero.Position()) <= e.awarenessDistance
}

func (e *Enemy) inAttackRange() bool {
	return rl.Vector2Distance(e.position, e.Hero.Position()) <= e.attackDistance
}

func (e *Enemy) inChaseRange() bool {
	return rl.Vector2Distance(e.position, e.Hero.Position()) >= e.minChaseDistance &&
		e.inAwarenessRange()
}

type EnemyState interface {
	Update(*Enemy, float32) EnemyState
}

type EnemyIdle struct {
	reasoning Timer
}

func (s *EnemyIdle) Update(e *Enemy, dt float32) EnemyState {

	e.animation = &e.animationIdle

	if e.inAttackRange() {
		return s
	}

	if e.inChaseRange() {
		return e.stateChase
	}

	return s
}

type EnemyChase struct{}

func (s *EnemyChase) Update(e *Enemy, dt float32) EnemyState {

	e.animation = &e.animationMove

	if e.inAttackRange() {
		return e.stateIdle
	}

	e.position = rl.Vector2MoveTowards(e.position, e.Hero.Position(), e.Speed*dt)

	return s
}
