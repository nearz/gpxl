package filter

import (
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
