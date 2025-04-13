package utils

func Lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}
