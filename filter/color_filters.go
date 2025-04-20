package filter

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/nearz/gpxl/pxl"
	"github.com/nearz/gpxl/utils"
)

const (
	shft8   = 8
	shft16  = 0
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
	start := time.Now()
	defer func() {
		fmt.Printf("Render time %s\n", time.Since(start))
	}()

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
			gs := luminance(float64(r>>shft8), float64(g>>shft8), float64(b>>shft8), clamp8)
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
			gs := luminance(float64(r), float64(g), float64(b), clamp16)
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

			dst.(*image.NRGBA).Set(x, y, color.NRGBA{uint8(fnlR), uint8(fnlG), uint8(fnlB), uint8(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewNRGBA(r)
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

			dst.(*image.NRGBA64).Set(x, y, color.RGBA64{uint16(fnlR), uint16(fnlG), uint16(fnlB), uint16(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewNRGBA64(r)
		},
	}
}

func Duotone(h, s color.RGBA) Filter {
	hr, hg, hb, _ := h.RGBA()
	sr, sg, sb, _ := s.RGBA()
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, a := src.At(x, y).RGBA()
			r8 := float64(r >> shft8)
			g8 := float64(g >> shft8)
			b8 := float64(b >> shft8)
			fnlR, fnlG, fnlB := duotoneCalc(r8, g8, b8, clamp8, shft8, hr, hg, hb, sr, sg, sb)

			dst.(*image.NRGBA).Set(x, y, color.NRGBA{uint8(fnlR), uint8(fnlG), uint8(fnlB), uint8(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewNRGBA(r)
		},
	}
}

func Duotone16(h, s color.RGBA) Filter {
	hr, hg, hb, _ := h.RGBA()
	sr, sg, sb, _ := s.RGBA()
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, a := src.At(x, y).RGBA()
			r16 := float64(r)
			g16 := float64(g)
			b16 := float64(b)
			fnlR, fnlG, fnlB := duotoneCalc(r16, g16, b16, clamp16, shft16, hr, hg, hb, sr, sg, sb)

			dst.(*image.NRGBA64).Set(x, y, color.NRGBA64{uint16(fnlR), uint16(fnlG), uint16(fnlB), uint16(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewNRGBA64(r)
		},
	}
}

func Cool(intensity float64) Filter {
	intensity = math.Min(math.Max(intensity, 0.0), 1.0)
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, a := src.At(x, y).RGBA()
			r8 := float64(r >> shft8)
			fnlG := float64(g >> shft8)
			b8 := float64(b >> shft8)
			fnlR, fnlB := coolCalc(r8, b8, clamp8, intensity)

			dst.(*image.NRGBA).Set(x, y, color.NRGBA{uint8(fnlR), uint8(fnlG), uint8(fnlB), uint8(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewNRGBA(r)
		},
	}
}

func Warm(intensity float64) Filter {
	intensity = math.Min(math.Max(intensity, 0.0), 1.0)
	return &colorFilter{
		fn: func(src, dst image.Image, x, y int) {
			r, g, b, a := src.At(x, y).RGBA()
			r8 := float64(r >> shft8)
			fnlG := float64(g >> shft8)
			b8 := float64(b >> shft8)
			fnlR, fnlB := coolCalc(r8, b8, clamp8, intensity)

			dst.(*image.NRGBA).Set(x, y, color.NRGBA{uint8(fnlR), uint8(fnlG), uint8(fnlB), uint8(a)})
		},
		dst: func(r image.Rectangle) image.Image {
			return image.NewNRGBA(r)
		},
	}
}

const (
	tempLo = 0.15
	tempHi = 0.25
)

func coolCalc(r, b, c, i float64) (fr, fb float64) {
	rt := tempLo * i
	bt := tempHi * i
	fr = utils.Lerp(r, c, rt)
	fb = utils.Lerp(b, 0, bt)

	fr = utils.Clamp(fr, c)
	fb = utils.Clamp(fb, c)
	return
}

func warmCalc(r, b, c, i float64) (fr, fb float64) {
	rt := tempHi * i
	bt := tempLo * i
	fr = utils.Lerp(r, 0, rt)
	fb = utils.Lerp(b, c, bt)

	fr = utils.Clamp(fr, c)
	fb = utils.Clamp(fb, c)
	return
}

const (
	rl = 0.299
	gl = 0.587
	bl = 0.114
)

func luminance(r, g, b, c float64) (l float64) {
	l = rl*r + gl*g + bl*b
	l = utils.Clamp(l, c)
	return
}

func duotoneCalc(r, g, b, c float64, shft int, hr, hg, hb, sr, sg, sb uint32) (fr, fg, fb float64) {
	l := luminance(r, g, b, c)
	t := l / c
	fr = utils.Lerp(float64(sr>>shft), float64(hr>>shft), t)
	fg = utils.Lerp(float64(sg>>shft), float64(hg>>shft), t)
	fb = utils.Lerp(float64(sb>>shft), float64(hb>>shft), t)

	fr = utils.Clamp(fr, c)
	fg = utils.Clamp(fg, c)
	fb = utils.Clamp(fb, c)
	return
}

const (
	s1 = 0.393
	s2 = 0.769
	s3 = 0.189
	s4 = 0.349
	s5 = 0.686
	s6 = 0.168
	s7 = 0.272
	s8 = 0.534
	s9 = 0.131
)

func sepiaClac(r, g, b, i, c float64) (fr, fg, fb float64) {
	sr := s1*r + s2*g + s3*b
	sg := s4*r + s5*g + s6*b
	sb := s7*r + s8*g + s9*b
	fr = utils.Lerp(r, sr, i)
	fg = utils.Lerp(g, sg, i)
	fb = utils.Lerp(b, sb, i)

	fr = utils.Clamp(fr, c)
	fg = utils.Clamp(fg, c)
	fb = utils.Clamp(fb, c)
	return
}
