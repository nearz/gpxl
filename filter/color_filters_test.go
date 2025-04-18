package filter

import (
	"image/color"
	"image/png"
	"testing"

	"github.com/nearz/gpxl/pxl"
)

func TestGreenTint(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Image read failed: %v", err)
	}
	f := GreenTint(1.0)
	f.Render(i)
	err = i.WritePNG("../test_images/write_tests/GreenTint.png", png.BestSpeed)
	if err != nil {
		t.Fatalf("Image write failed: %v", err)
	}
}

func TestDuotone16(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Image read failed: %v", err)
	}
	h := color.RGBA{R: 16, G: 197, B: 248, A: 255}
	s := color.RGBA{R: 103, G: 54, B: 221, A: 255}
	f := Duotone16(h, s)
	f.Render(i)
	err = i.WritePNG("../test_images/write_tests/Duo16.png", png.BestSpeed)
	if err != nil {
		t.Fatalf("Image write failed: %v", err)
	}
}

func TestDuotone(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Image read failed: %v", err)
	}
	h := color.RGBA{R: 16, G: 197, B: 248, A: 255}
	s := color.RGBA{R: 103, G: 54, B: 221, A: 255}
	f := Duotone(h, s)
	f.Render(i)
	err = i.WritePNG("../test_images/write_tests/Duo.png", png.BestSpeed)
	if err != nil {
		t.Fatalf("Image write failed: %v", err)
	}
}
