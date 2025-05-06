package pxl

/*
TODO:
- Should I convert everything to RGBA or RGBA64?
*/

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

type ImageFormat string

// Will there be a use for these?
const (
	JPG ImageFormat = "jpeg"
	PNG ImageFormat = "png"
	GIF ImageFormat = "gif"
)

type Pxl struct {
	Image image.Image
}

func Read(path string) (*Pxl, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	pxl := &Pxl{
		Image: img,
	}

	return pxl, nil
}

func (p *Pxl) WritePNG(path string, cl png.CompressionLevel) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := png.Encoder{
		CompressionLevel: cl,
	}
	return enc.Encode(f, p.Image)
}

func (p *Pxl) WriteJPEG(path string, quality int) error {
	fq := math.Min(math.Max(float64(quality), 0.0), 100.0)
	quality = int(fq)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, p.Image, &jpeg.Options{Quality: quality})
}

func (p *Pxl) PrintType() {
	fmt.Printf("Type: %T", p.Image)
}

// Should just be New RGBA. Not format specific, format with happen on Write.
// func Create(w, h int) *Pxl {
// 	// Create a new RGBA image with 100x100 dimensions
// 	img := image.NewRGBA(image.Rect(0, 0, w, h))
//
// 	// Light gray color (RGB: 200, 200, 200)
// 	lightGray := color.RGBA{200, 200, 200, 255}
//
// 	// Fill the entire image with light gray
// 	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
// 		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
// 			img.Set(x, y, lightGray)
// 		}
// 	}
// 	pxl := &Pxl{
// 		Image: img,
// 		// Format: "png",
// 	}
// 	return pxl
// }
