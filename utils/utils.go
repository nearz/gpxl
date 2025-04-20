package utils

import "math"

func Lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func Clamp(v, c float64) float64 {
	return math.Max(0, math.Min(c, v))
}
