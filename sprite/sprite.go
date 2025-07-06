package sprite

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewSheet(t rl.Texture2D, speed, spriteSize float32, row1, col1, row2, col2 int) Sheet {
	var s Sheet
	s.Texture = t
	s.Speed = speed
	for c := col1; c <= col2; c++ {
		for r := row1; r <= row2; r++ {
			s.Frames = append(s.Frames, rl.Rectangle{
				X:      float32(c) * spriteSize,
				Y:      float32(r) * spriteSize,
				Width:  spriteSize,
				Height: spriteSize,
			})
		}
	}
	return s
}

type Sheet struct {
	Texture   rl.Texture2D
	Frames    []rl.Rectangle
	Speed     float32
	Flip      bool
	frameTime float32
	frame     int
	complete  bool
}

func (s *Sheet) Update(dt float32) {
	s.frameTime += dt * s.Speed
	frame := int(s.frameTime) % len(s.Frames)
	s.complete = frame < s.frame
	s.frame = frame
}

func (s *Sheet) Draw(p rl.Vector2) {
	f := s.Frames[s.frame]
	if s.Flip {
		f.Width = -f.Width
	}
	rl.DrawTextureRec(s.Texture, f, p, rl.White)
}

func (s *Sheet) Complete() bool {
	return s.complete
}
