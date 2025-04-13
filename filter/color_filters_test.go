package filter

import (
	"image/color"
	"testing"

	"github.com/nearz/gpxl/pxl"
)

func TestGray(t *testing.T) {
	i, err := pxl.Read("../test_images/penguin.png")
	if err != nil {
		t.Error(err)
		return
	}
	f := Grayscale()
	f.Render(i)
	err = i.Write("../test_images/write_tests/GrayTest.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

func TestGrayRows(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Error(err)
		return
	}
	f := Grayscale()
	f.RenderRows(i)
	err = i.Write("../test_images/write_tests/GrayRowsTest.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

func TestGraySeq(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Error(err)
		return
	}
	f := Grayscale()
	f.RenderSequential(i)
	err = i.Write("../test_images/write_tests/GrayTestSeq.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

func TestGrayscale16(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant16.png")
	if err != nil {
		t.Error(err)
		return
	}
	f := Grayscale16()
	f.Render(i)
	err = i.Write("../test_images/write_tests/GrayTest16.png", pxl.PNG)
}

func TestSepia(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Error(err)
		return
	}
	f := Sepia(1.0)
	f.Render(i)
	err = i.Write("../test_images/write_tests/SepiaTest.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

// func TestSepiaSeq(t *testing.T) {
// 	i, err := pxl.Read("../test_images/elephant.png")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	f := Sepia(1.0)
// 	f.(*colorFilter).RenderSequential(i)
// 	err = i.Write("../test_images/write_tests/SepiaTestSeq.png", pxl.PNG)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestSepia16(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant16.png")
	if err != nil {
		t.Error(err)
		return
	}
	f := Sepia16(1.0)
	f.Render(i)
	err = i.Write("../test_images/write_tests/SepiaTest16.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

func TestDuotone(t *testing.T) {
	i, err := pxl.Read("../test_images/deer.png")
	if err != nil {
		t.Error(err)
		return
	}
	h := color.RGBA{R: 255, G: 0, B: 40, A: 255}
	s := color.RGBA{R: 122, G: 0, B: 79, A: 255}
	f := Duotone(h, s)
	f.Render(i)
	err = i.Write("../test_images/write_tests/DuotoneTest.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

func TestDuotoneBlock(t *testing.T) {
	i, err := pxl.Read("../test_images/deer.png")
	if err != nil {
		t.Error(err)
		return
	}
	h := color.RGBA{R: 255, G: 0, B: 40, A: 255}
	s := color.RGBA{R: 122, G: 0, B: 79, A: 255}
	f := Duotone(h, s)
	f.(*colorFilter).RenderBlock(i)
	err = i.Write("../test_images/write_tests/DuotoneTestBlock.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

func TestDuotoneSeq(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Error(err)
		return
	}
	h := color.RGBA{R: 255, G: 0, B: 40, A: 255}
	s := color.RGBA{R: 122, G: 0, B: 79, A: 255}
	f := Duotone(h, s)
	f.(*colorFilter).RenderSequential(i)
	err = i.Write("../test_images/write_tests/DuotoneTestSeq.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}

func TestDuotoneRows(t *testing.T) {
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Error(err)
		return
	}
	h := color.RGBA{R: 255, G: 0, B: 40, A: 255}
	s := color.RGBA{R: 122, G: 0, B: 79, A: 255}
	f := Duotone(h, s)
	f.(*colorFilter).RenderRows(i)
	err = i.Write("../test_images/write_tests/DuotoneTestRows.png", pxl.PNG)
	if err != nil {
		t.Error(err)
	}
}
