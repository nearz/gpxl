package filter

// Do all the numerical type conversions cause performance issues?

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"

	"github.com/nearz/gpxl/pxl"
	"github.com/nearz/gpxl/utils"
)

const (
	shft8  = 8
	shft16 = 16
)

// For grayscale calculations
const (
	rl = 0.299
	gl = 0.587
	bl = 0.114
)

// For sepia calculations
const (
	s1      = 0.393
	s2      = 0.769
	s3      = 0.189
	s4      = 0.349
	s5      = 0.686
	s6      = 0.168
	s7      = 0.272
	s8      = 0.534
	s9      = 0.131
	clamp8  = 255
	clamp16 = 65535
)

type colorFilter struct {
	fn  func(image.Image, image.Image, int, int)
	dst func(image.Rectangle) image.Image
}

func (c *colorFilter) Render(p *pxl.Pxl) {
	// What checks can I make??
	// How to test or benchmark goroutines

	bds := p.Image.Bounds()
	dst := c.dst(bds)

	ncpus := runtime.NumCPU()
	trows := bds.Max.Y - bds.Min.Y
	chunkSize := (trows + ncpus - 1) / ncpus
	var wg sync.WaitGroup
	for sy := bds.Min.Y; sy < bds.Max.Y; sy += chunkSize {
		wg.Add(1)
		ey := sy + chunkSize
		if ey > bds.Max.Y {
			ey = bds.Max.Y
		}
		go func(start, end int) {
			defer wg.Done()
			for y := start; y < end; y++ {
				procPixRow(p.Image, dst, c.fn, bds.Min.X, bds.Max.X, y)
			}
		}(sy, ey)

	}
	wg.Wait()
	p.Image = dst
}

func procPixRow(src image.Image, dst image.Image, f func(image.Image, image.Image, int, int), minX, maxX, y int) {
	for x := minX; x < maxX; x++ {
		f(src, dst, x, y)
	}
}

func Grayscale() Filter {
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, _ := src.At(x, y).RGBA()
			gs := luminance(float64(r>>8), float64(g>>8), float64(b>>8))
			dst.(*image.Gray).Set(x, y, color.Gray{Y: uint8(gs)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewGray(r)
		},
	}
}

func Grayscale16() Filter {
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, _ := src.At(x, y).RGBA()
			gs := luminance(float64(r), float64(g), float64(b))
			dst.(*image.Gray16).Set(x, y, color.Gray16{Y: uint16(gs)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewGray16(r)
		},
	}
}

func Sepia(intensity float64) Filter {
	intensity = math.Min(math.Max(intensity, 0.0), 1.0)
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, a := src.At(x, y).RGBA()
			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)
			fnlR, fnlG, fnlB := sepiaClac(r8, g8, b8, intensity, clamp8)

			dst.(*image.RGBA).Set(x, y, color.RGBA{uint8(fnlR), uint8(fnlG), uint8(fnlB), uint8(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewRGBA(r)
		},
	}
}

func Sepia16(intensity float64) Filter {
	intensity = math.Min(math.Max(intensity, 0.0), 1.0)
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, a := src.At(x, y).RGBA()
			r16 := float64(r)
			g16 := float64(g)
			b16 := float64(b)
			fnlR, fnlG, fnlB := sepiaClac(r16, g16, b16, intensity, clamp16)

			dst.(*image.RGBA64).Set(x, y, color.RGBA64{uint16(fnlR), uint16(fnlG), uint16(fnlB), uint16(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewRGBA64(r)
		},
	}
}

func sepiaClac(r, g, b, i, c float64) (fr, fg, fb float64) {
	sr := s1*r + s2*g + s3*b
	sg := s4*r + s5*g + s6*b
	sb := s7*r + s8*g + s9*b
	fr = utils.Lerp(r, sr, i)
	fg = utils.Lerp(g, sg, i)
	fb = utils.Lerp(b, sb, i)
	if fr > c {
		fr = c
	}
	if fg > c {
		fg = c
	}
	if fb > c {
		fb = c
	}
	return
}

func Duotone(h, s color.RGBA) Filter {
	hr, hg, hb, _ := h.RGBA()
	sr, sg, sb, _ := s.RGBA()
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, a := src.At(x, y).RGBA()
			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)
			fnlR, fnlG, fnlB := duotoneCalc(r8, g8, b8, clamp8, shft8, hr, hg, hb, sr, sg, sb)

			dst.(*image.RGBA).Set(x, y, color.RGBA{uint8(fnlR), uint8(fnlG), uint8(fnlB), uint8(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewRGBA(r)
		},
	}
}

func duotoneCalc(r, g, b, c float64, shft int, hr, hg, hb, sr, sg, sb uint32) (rd, gd, bd float64) {
	l := luminance(r, g, b)
	t := l / c
	rd = utils.Lerp(float64(sr>>shft), float64(hr>>shft), t)
	gd = utils.Lerp(float64(sg>>shft), float64(hg>>shft), t)
	bd = utils.Lerp(float64(sb>>shft), float64(hb>>shft), t)
	if rd > c {
		rd = c
	}
	if gd > c {
		gd = c
	}
	if bd > c {
		bd = c
	}
	return
}

func luminance(r, g, b float64) float64 {
	return rl*r + gl*g + bl*b
}
