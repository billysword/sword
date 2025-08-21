package engine

import "math"

// EaseValue maps a normalized progress t in [0,1] through a named easing curve.
// Supported names: "linear", "easeIn", "easeOut", "easeInOut", "power" (uses power argument).
// Returns value in [0,1].
func EaseValue(name string, t, power float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	switch name {
	case "easeIn":
		// cubic ease-in
		return t * t * t
	case "easeOut":
		// cubic ease-out
		u := 1 - t
		return 1 - u*u*u
	case "easeInOut":
		// smoothstep-like
		return t * t * (3 - 2*t)
	case "power":
		p := power
		if p <= 0 {
			p = 2
		}
		return math.Pow(t, p)
	case "linear":
		fallthrough
	default:
		return t
	}
}
