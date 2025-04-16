package pxl

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"time"
)

type ImageFormat string

const (
	JPG ImageFormat = "jpeg"
	PNG ImageFormat = "png"
	GIF ImageFormat = "gif"
)

type Pxl struct {
	Image image.Image
}

func Read(path string) (*Pxl, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("Decode time %s\n", time.Since(start))
	}()

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

func (p *Pxl) Write(path string, format ImageFormat) error {
	start := time.Now()
	defer func() {
		fmt.Printf("Encode time %s\n", time.Since(start))
	}()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case JPG:
		return jpeg.Encode(f, p.Image, &jpeg.Options{Quality: 95})
	case PNG:
		return png.Encode(f, p.Image)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
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
