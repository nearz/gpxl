package transform

import "math"

// Is this a kernel? I have seen this called filter too??
// Look up an rename maybe?
type Kernel struct {
	Radius int
	Fn     func(t float64) float64
}

func (k *Kernel) Weights(f float64) []float64 {
	n := 2 * k.Radius
	w := make([]float64, n)

	for i := 0; i < n; i++ {
		offset := float64(i - (k.Radius - 1))
		w[i] = k.Fn(offset - f)
	}
	return w
}

var (
	Nearest    Kernel
	Linear     Kernel
	CatmullRom Kernel
)

func init() {
	Nearest = Kernel{
		Radius: 0,
		Fn:     nil,
	}
	Linear = Kernel{
		Radius: 1,
		Fn: func(t float64) float64 {
			a := math.Abs(t)
			if a < 1 {
				return 1 - a
			}
			return 0
		},
	}
	// Bug - Color artifacts that need fix. e.g. pixel full Red.
	CatmullRom = Kernel{
		Radius: 2,
		Fn: func(t float64) float64 {
			a := -0.5
			at := math.Abs(t)
			if at < 1 {
				return (a+2)*at*at*at - (a+3)*at*at + 1
			}
			if at < 2 {
				return a*at*at*at - 5*a*at*at + 8*a*at - 4*a
			}
			return 0
		},
	}
}
