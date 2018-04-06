package tween

// Easing is the mapping function used by a Tween
type Easing func(float32) float32

// Linear ease
func Linear(t float32) float32 {
	return t
}

// EaseInOut start slow and end slow
func EaseInOut(t float32) float32 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseIn starts slow
func EaseIn(t float32) float32 {
	return t * t
}

// EaseOut ends slow
func EaseOut(t float32) float32 {
	return t * (2 - t)
}
