package transform

/*
TODO:
- How to handle different color modles?
- Should I make it a method of transformer?
*/

import (
	"image"

	"github.com/nearz/gpxl/concurrent"
	"github.com/nearz/gpxl/pxl"
)

func ReflectH(p *pxl.Pxl) {
	bds := p.Image.Bounds()
	dst := image.NewRGBA(bds)
	concurrent.Rows(bds.Dy(), bds.Min.Y, func(start, end int) {
		for y := start; y < end; y++ {
			for x := bds.Min.X; x < bds.Max.X; x++ {
				dx := bds.Max.X - 1 - x
				c := p.Image.At(x, y)
				dst.Set(dx, y, c)
			}
		}
	})
	p.Image = dst
}

func ReflectV(p *pxl.Pxl) {
	bds := p.Image.Bounds()
	dst := image.NewRGBA(bds)
	concurrent.Rows(bds.Dy(), bds.Min.Y, func(start, end int) {
		for y := start; y < end; y++ {
			for x := bds.Min.X; x < bds.Max.X; x++ {
				dy := bds.Max.Y - 1 - y
				c := p.Image.At(x, y)
				dst.Set(x, dy, c)
			}
		}
	})
	p.Image = dst
}
