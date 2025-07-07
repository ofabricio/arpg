package arpg

import "math"

// PingPong returns a value that increments and decrements as t advances.
// It is 0.0 when t is 0.0; 1.0 when t is 0.5; and 0.0 when t is 1.0.
// Multiply or divide t by a factor to change the frequency of the wave.
// For example, PingPong(t*2) will make the wave twice as frequent.
func PingPong(t float32) float32 {
	// This implements the triangle wave formula as in
	// https://en.wikipedia.org/wiki/Triangle_wave.
	// Paste this in www.geogebra.org/classic to see the wave:
	// abs(x - floor(x + 0.5)) * 2
	mod := t - float32(math.Floor(float64(t+0.5)))
	return float32(math.Abs(float64(mod))) * 2
}
