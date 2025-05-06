package transform

import (
	"fmt"
	"image/png"
	"testing"
	"time"

	"github.com/nearz/gpxl/pxl"
)

func TestRotate(t *testing.T) {
	s := time.Now()
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Read image failed: %v", err)
	}
	fmt.Printf("Read time: %s\n", time.Since(s))

	s = time.Now()
	optimus := Bilinear()
	optimus.Rotate(i, 45.0)
	fmt.Printf("Rotate time: %s\n", time.Since(s))

	s = time.Now()
	err = i.WritePNG("../test_images/write_tests/e45BL.png", png.BestSpeed)
	if err != nil {
		t.Fatalf("Write image failed: %v", err)
	}
	fmt.Printf("Write time: %s\n", time.Since(s))
}

func TestRotater(t *testing.T) {
	s := time.Now()
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Read image failed: %v", err)
	}
	fmt.Printf("Read time: %s\n", time.Since(s))

	s = time.Now()
	r := Rotater(CatmullRom)
	r.Rotate(i, 90.0)
	fmt.Printf("Rotate time: %s\n", time.Since(s))

	s = time.Now()
	err = i.WritePNG("../test_images/write_tests/e90CR.png", png.BestSpeed)
	if err != nil {
		t.Fatalf("Write image failed: %v", err)
	}
	fmt.Printf("Write time: %s\n", time.Since(s))
}

func TestReflectH(t *testing.T) {
	s := time.Now()
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Read image failed: %v", err)
	}
	fmt.Printf("Read time: %s\n", time.Since(s))

	s = time.Now()
	ReflectH(i)
	fmt.Printf("Flip Horz. time:%s\n", time.Since(s))

	s = time.Now()
	err = i.WritePNG("../test_images/write_tests/eFlipH.png", png.BestSpeed)
	if err != nil {
		t.Fatalf("Write image failed: %v", err)
	}
	fmt.Printf("Write time: %s\n", time.Since(s))
}

func TestReflectV(t *testing.T) {
	s := time.Now()
	i, err := pxl.Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Read image failed: %v", err)
	}
	fmt.Printf("Read time: %s\n", time.Since(s))

	s = time.Now()
	ReflectV(i)
	fmt.Printf("Flip Horz. time:%s\n", time.Since(s))

	s = time.Now()
	err = i.WritePNG("../test_images/write_tests/eFlipV.png", png.BestSpeed)
	if err != nil {
		t.Fatalf("Write image failed: %v", err)
	}
	fmt.Printf("Write time: %s\n", time.Since(s))
}
