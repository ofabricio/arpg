package arpg

import (
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ofabricio/arpg/sprite"
)

func NewEnemy(x, y float32) Enemy {
	var e Enemy
	e.animationMove = sprite.NewSheet(rl.LoadTexture("assets/warrior_run.png"), 10, 192, 0, 0, 0, 5)
	e.animationIdle = sprite.NewSheet(rl.LoadTexture("assets/warrior_idle.png"), 10, 192, 0, 0, 0, 7)
	e.animationGuard = sprite.NewSheet(rl.LoadTexture("assets/warrior_guard.png"), 10, 192, 0, 0, 0, 5)
	e.animationAttack1 = sprite.NewSheet(rl.LoadTexture("assets/warrior_attack1.png"), 10, 192, 0, 0, 0, 3)
	e.animationAttack2 = sprite.NewSheet(rl.LoadTexture("assets/warrior_attack2.png"), 10, 192, 0, 0, 0, 3)
	e.animation = &e.animationIdle
	e.animationOffset = rl.Vector2{X: 96, Y: 96 + 33} // Half of the sprite size (192x192).
	e.stateIdle = &EnemyIdle{}
	e.stateChase = &EnemyChase{}
	e.stateAttack = &EnemyAttack{}
	e.state = e.stateIdle
	e.position.X = x
	e.position.Y = y
	e.AttackDistance = 60
	e.AttackSpeed = 1
	e.ChaseDistance = 500
	e.Speed = 80
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

	state       EnemyState
	stateIdle   EnemyState
	stateChase  EnemyState
	stateAttack EnemyState

	Hero  *Hero
	Speed float32

	position       rl.Vector2
	AttackSpeed    float32
	AttackDistance float32
	ChaseDistance  float32

	hurtAnimation HurtAnimation
}

func (e *Enemy) Update(dt float32) {

	next := e.state.Update(e, dt)
	if next != nil {
		if v, ok := next.(interface{ Enter(*Enemy) }); ok {
			v.Enter(e)
		}
		e.state = next
	}

	e.animation.Flip = rl.Vector2Subtract(e.Hero.Position(), e.position).X < 0
	e.animation.Update(dt)

	e.animation.Tint = e.hurtAnimation.Value()
	e.hurtAnimation.Update(dt)
}

func (e *Enemy) Draw() {
	p := e.position
	p.X -= e.animationOffset.X
	p.Y -= e.animationOffset.Y
	e.animation.Draw(p)

	if DebugEnabled {
		rl.DrawCircleLines(int32(e.position.X), int32(e.position.Y), e.AttackDistance, rl.Red)
		rl.DrawCircleLines(int32(e.position.X), int32(e.position.Y), e.ChaseDistance, rl.Gray)
		rl.DrawText(reflect.TypeOf(e.state).Elem().Name(), int32(e.position.X), int32(e.position.Y)-100, 20, rl.DarkGray)
	}
}

func (e *Enemy) Position() rl.Vector2 {
	return e.position
}

func (e *Enemy) Hurt() {
	e.hurtAnimation.Play()
}

func (e *Enemy) inChaseRange() bool {
	return rl.Vector2Distance(e.position, e.Hero.Position()) <= e.ChaseDistance
}

func (e *Enemy) inAttackRange() bool {
	return rl.Vector2Distance(e.position, e.Hero.Position()) <= e.AttackDistance
}

type EnemyState interface {
	Update(*Enemy, float32) EnemyState
}

type EnemyIdle struct{}

func (s *EnemyIdle) Update(e *Enemy, dt float32) EnemyState {

	if e.inAttackRange() {
		return e.stateAttack
	}

	if e.inChaseRange() {
		return e.stateChase
	}

	e.animation = &e.animationIdle

	return nil
}

type EnemyChase struct{}

func (s *EnemyChase) Update(e *Enemy, dt float32) EnemyState {

	if e.inAttackRange() {
		return e.stateIdle
	}

	e.position = rl.Vector2MoveTowards(e.position, e.Hero.Position(), e.Speed*dt)
	e.animation = &e.animationMove

	return nil
}

type EnemyAttack struct {
	time float32
}

func (s *EnemyAttack) Enter(e *Enemy) {
	e.animation = &e.animationIdle
}

func (s *EnemyAttack) Update(e *Enemy, dt float32) EnemyState {
	if e.animation.Completed {
		e.animation = &e.animationIdle
	}
	if e.inAttackRange() {
		s.time += dt
		if s.time >= 1/e.AttackSpeed {
			s.time = 0
			e.Hero.Hurt()
			e.animation = &e.animationAttack1
		}
	} else {
		if e.animation.Completed {
			return e.stateIdle
		}
	}

	return nil
}
