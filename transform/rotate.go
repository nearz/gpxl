package transform

import (
	"image"
	"image/color"
	"math"

	"github.com/nearz/gpxl/concurrent"
	"github.com/nearz/gpxl/pxl"
)

/*
Todo:
- Degree check, 0, 90, 180 270
  - Add some ifs because not everything is required?
- why getting extra dim on 270?
- Sampling and Interpolation
  - NN: Finished?
  - Bilinear: Finished?
  - Bicubic: Todo
- How does clamp work?
- How to handle different color models
- Is it more effiecient to resize larger since operation is seperable for rows and columns?
*/

// Rename to make more sense.
// I dont think there will be a general transformer struct
// for all transformation funcs.
type transformer struct {
	fn func(float64, float64, image.Image) color.Color
}

func (t *transformer) Rotate(p *pxl.Pxl, angle float64) error {
	radians := angle * math.Pi / 180.0
	cos, sin := math.Cos(radians), math.Sin(radians)

	src := p.Image
	sbds := src.Bounds()
	w, h := sbds.Dx(), sbds.Dy()
	ixc := float64(w-1) / 2
	iyc := float64(h-1) / 2

	wo := int(math.Ceil(math.Abs(float64(w)*cos) + math.Abs(float64(h)*sin)))
	ho := int(math.Ceil(math.Abs(float64(w)*sin) + math.Abs(float64(h)*cos)))
	oxc := float64(wo-1) / 2
	oyc := float64(ho-1) / 2

	dst := image.NewRGBA(image.Rect(0, 0, wo, ho))
	dbds := dst.Bounds()

	rows := dbds.Max.Y - dbds.Min.Y
	concurrent.Rows(rows, dbds.Min.Y, func(start, end int) {
		for y := start; y < end; y++ {
			for x := dbds.Min.X; x < dbds.Max.X; x++ {
				dx := float64(x) - oxc
				dy := float64(y) - oyc
				srcx := cos*dx + sin*dy + ixc
				srcy := -sin*dx + cos*dy + iyc
				resColor := t.fn(srcx, srcy, src)
				dst.Set(x, y, resColor)
			}
		}
	})

	p.Image = dst
	return nil
}

// Rotater?
func Rotater(k Kernel) *transformer {
	if k.Radius <= 0 {
		return &transformer{
			fn: func(x, y float64, src image.Image) color.Color {
				xi := int(x + 0.5)
				yi := int(y + 0.5)
				xi = clampEdge(xi, src.Bounds().Min.X, src.Bounds().Max.X)
				yi = clampEdge(yi, src.Bounds().Min.Y, src.Bounds().Max.Y)
				if !inBounds(xi, yi, src.Bounds()) {
					return color.RGBA{0, 0, 0, 0}
				}
				c := src.At(xi, yi)
				return c
			},
		}
	} else {
		// Bounds check for after calmpEdge?
		return &transformer{
			fn: func(x, y float64, src image.Image) color.Color {
				x0 := int(math.Floor(x)) - (k.Radius - 1)
				y0 := int(math.Floor(y)) - (k.Radius - 1)
				fx := x - math.Floor(x)
				fy := y - math.Floor(y)
				wx := k.Weights(fx)
				wy := k.Weights(fy)

				taps := 2 * k.Radius
				var r, g, b, a float64
				for j := 0; j < taps; j++ {
					sy := clampEdge(y0+j, src.Bounds().Min.Y, src.Bounds().Max.Y)
					for i := 0; i < taps; i++ {
						sx := clampEdge(x0+i, src.Bounds().Min.X, src.Bounds().Max.X)
						c := src.At(sx, sy).(color.RGBA)
						wxy := wx[i] * wy[j]
						r += wxy * float64(c.R)
						g += wxy * float64(c.G)
						b += wxy * float64(c.B)
						a += wxy * float64(c.A)
					}
				}
				return color.RGBA{uint8(r + 0.5), uint8(g + 0.5), uint8(b + 0.5), uint8(a + 0.5)}
			},
		}
	}
}

// Maybe have an out of range flag? So that I can return a blank pizel if so? For angle not of 90 degree increments.
func clampEdge(a, min, max int) int {
	if a == min-1 {
		return min
	} else if a == max {
		return max - 1
	}
	return a
}

func inBounds(x, y int, r image.Rectangle) bool {
	return x >= r.Min.X && x <= r.Max.X-1 && y >= r.Min.Y && y <= r.Max.Y-1
}
