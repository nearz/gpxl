package pxl

import (
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

func TestWrite(t *testing.T) {
	pxl, err := Read("../test_images/small.png")
	if err != nil {
		t.Fatalf("Failed to read test image: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		format  ImageFormat
		wantErr bool
	}{
		{"Write PNG", "../test_images/write_tests/write_test.png", PNG, false},
		{"Write JPEG", "../test_images/write_tests/write_test.jpg", JPG, false},
		{"Write JPG", "../test_images/write_tests/write_test.jpg", JPG, false},
		{"Invalid format", "../test_images/write_tests/write_test.tiff", "tiff", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pxl.Write(tt.path, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
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
	err = pxl.Write("../test_images/write_tests/ElWriteTest.png", PNG)
	if err != nil {
		t.Errorf("Failed to write image: %v\n", err)
	}
}
