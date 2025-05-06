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

func NearestNeighbor() *transformer {
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
}

func Bilinear() *transformer {
	return &transformer{
		fn: func(x, y float64, src image.Image) color.Color {
			fx := x - math.Floor(x)
			fy := y - math.Floor(y)
			x0, y0 := int(math.Floor(x)), int(math.Floor(y))
			x0 = clampEdge(x0, src.Bounds().Min.X, src.Bounds().Max.X)
			y0 = clampEdge(y0, src.Bounds().Min.Y, src.Bounds().Max.Y)
			if !inBounds(x0, y0, src.Bounds()) {
				return color.RGBA{0, 0, 0, 0}
			}
			x1, y1 := x0+1, y0+1
			x1 = clampEdge(x1, src.Bounds().Min.X, src.Bounds().Max.X)
			y1 = clampEdge(y1, src.Bounds().Min.Y, src.Bounds().Max.Y)
			c00 := src.At(x0, y0).(color.RGBA)
			c01 := src.At(x0, y1).(color.RGBA)
			c10 := src.At(x1, y0).(color.RGBA)
			c11 := src.At(x1, y1).(color.RGBA)
			r := filterLerp(filterLerp(c00.R, c10.R, fx), filterLerp(c01.R, c11.R, fx), fy)
			g := filterLerp(filterLerp(c00.G, c10.G, fx), filterLerp(c01.G, c11.G, fx), fy)
			b := filterLerp(filterLerp(c00.B, c10.B, fx), filterLerp(c01.B, c11.B, fx), fy)
			a := filterLerp(filterLerp(c00.A, c10.A, fx), filterLerp(c01.A, c11.A, fx), fy)
			c := color.RGBA{r, g, b, a}
			return c
		},
	}
}

func filterLerp(a, b uint8, t float64) uint8 {
	af := float64(a)
	bf := float64(b)
	return uint8(math.Round(af + t*(bf-af)))
}

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
