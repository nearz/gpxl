package pxl

import (
	"image/png"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Valid PNG", "../test_images/small.png", false},
		{"Valid JPEG", "../test_images/small.jpg", false},
		{"Non-existent file", "../test_images/nonexistent.png", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pxl, err := Read(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && pxl == nil {
				t.Error("Read() returned nil Pxl without error")
			}
		})
	}
}

func TestWritePNG(t *testing.T) {
	pxl, err := Read("../test_images/small.png")
	if err != nil {
		t.Fatalf("Failed to read test image: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		level   png.CompressionLevel
		wantErr bool
	}{
		{"Write PNG Default", "../test_images/write_tests/write_test.png", png.DefaultCompression, false},
		{"Write PNG NoCompression", "../test_images/write_tests/write_test.png", png.NoCompression, false},
		{"Write PNG BestSpeed", "../test_images/write_tests/write_test.png", png.BestSpeed, false},
		{"Write PNG BestCompression", "../test_images/write_tests/write_test.png", png.BestCompression, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pxl.WritePNG(tt.path, tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("WritePNG() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				os.Remove(tt.path)
			}
		})
	}
}

func TestWriteJPEG(t *testing.T) {
	pxl, err := Read("../test_images/small.png")
	if err != nil {
		t.Fatalf("Failed to read test image: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		quality int
		wantErr bool
	}{
		{"Write JPEG Quality 100", "../test_images/write_tests/write_test.jpg", 100, false},
		{"Write JPEG Quality 75", "../test_images/write_tests/write_test.jpg", 75, false},
		{"Write JPEG Quality 50", "../test_images/write_tests/write_test.jpg", 50, false},
		{"Write JPEG Quality 0", "../test_images/write_tests/write_test.jpg", 0, false},
		{"Write JPEG Quality -1", "../test_images/write_tests/write_test.jpg", -1, false},
		{"Write JPEG Quality 101", "../test_images/write_tests/write_test.jpg", 101, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pxl.WriteJPEG(tt.path, tt.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteJPEG() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				os.Remove(tt.path)
			}
		})
	}
}

func TestPrintType(t *testing.T) {
	pxl, err := Read("../test_images/elephant.jpeg")
	if err != nil {
		t.Fatalf("Failed to read image: %v", err)
	}
	pxl.PrintType()
}

func TestSandbox(t *testing.T) {
	pxl, err := Read("../test_images/elephant.png")
	if err != nil {
		t.Fatalf("Failed to read image: %v\n", err)
	}
	err = pxl.WritePNG("../test_images/write_tests/ElWriteTest.png", png.DefaultCompression)
	if err != nil {
		t.Errorf("Failed to write image: %v\n", err)
	}
}
