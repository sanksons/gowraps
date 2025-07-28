package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

// Helper function to create a simple test image
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// Fill with red color
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	return img
}

// Helper function to create PNG bytes
func createPNGBytes() []byte {
	img := createTestImage()
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

// Helper function to create JPEG bytes
func createJPEGBytes() []byte {
	img := createTestImage()
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	return buf.Bytes()
}

// Helper function to create GIF bytes (mock)
func createGIFBytes() []byte {
	return []byte("GIF87a\x01\x00\x01\x00\x00\x00\x00!")
}

func TestGetMime(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "JPEG data",
			data:     []byte("\xff\xd8\xff\xe0\x00\x10JFIF"),
			expected: MIME_TYPE_JPEG,
			wantErr:  false,
		},
		{
			name:     "PNG data",
			data:     []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"),
			expected: MIME_TYPE_PNG,
			wantErr:  false,
		},
		{
			name:     "GIF87a data",
			data:     []byte("GIF87a\x01\x00\x01\x00"),
			expected: MIME_TYPE_GIF,
			wantErr:  false,
		},
		{
			name:     "GIF89a data",
			data:     []byte("GIF89a\x01\x00\x01\x00"),
			expected: MIME_TYPE_GIF,
			wantErr:  false,
		},
		{
			name:     "Unknown format",
			data:     []byte("unknown format"),
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty data",
			data:     []byte{},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Real PNG bytes",
			data:     createPNGBytes(),
			expected: MIME_TYPE_PNG,
			wantErr:  false,
		},
		{
			name:     "Real JPEG bytes",
			data:     createJPEGBytes(),
			expected: MIME_TYPE_JPEG,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetMime(tt.data)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetMime() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("GetMime() unexpected error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("GetMime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetExtension4mMime(t *testing.T) {
	tests := []struct {
		name     string
		mime     string
		expected string
		wantErr  bool
	}{
		{
			name:     "JPEG mime type",
			mime:     MIME_TYPE_JPEG,
			expected: EXT_JPEG,
			wantErr:  false,
		},
		{
			name:     "PNG mime type",
			mime:     MIME_TYPE_PNG,
			expected: EXT_PNG,
			wantErr:  false,
		},
		{
			name:     "GIF mime type",
			mime:     MIME_TYPE_GIF,
			expected: MIME_TYPE_GIF,
			wantErr:  false,
		},
		{
			name:     "Invalid mime type",
			mime:     "invalid/mime",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty mime type",
			mime:     "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetExtension4mMime(tt.mime)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetExtension4mMime() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("GetExtension4mMime() unexpected error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("GetExtension4mMime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetMime4mExt(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected string
		wantErr  bool
	}{
		{
			name:     "JPEG extension",
			ext:      EXT_JPEG,
			expected: MIME_TYPE_JPEG,
			wantErr:  false,
		},
		{
			name:     "JPG extension",
			ext:      EXT_JPG,
			expected: MIME_TYPE_JPEG,
			wantErr:  false,
		},
		{
			name:     "PNG extension",
			ext:      EXT_PNG,
			expected: MIME_TYPE_PNG,
			wantErr:  false,
		},
		{
			name:     "GIF extension",
			ext:      EXT_GIF,
			expected: MIME_TYPE_GIF,
			wantErr:  false,
		},
		{
			name:     "Invalid extension",
			ext:      "invalid",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty extension",
			ext:      "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetMime4mExt(tt.ext)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetMime4mExt() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("GetMime4mExt() unexpected error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("GetMime4mExt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetBytes4mImage(t *testing.T) {
	testImg := createTestImage()

	tests := []struct {
		name    string
		img     image.Image
		mime    string
		wantErr bool
	}{
		{
			name:    "Valid image to PNG",
			img:     testImg,
			mime:    MIME_TYPE_PNG,
			wantErr: false,
		},
		{
			name:    "Valid image to JPEG",
			img:     testImg,
			mime:    MIME_TYPE_JPEG,
			wantErr: false,
		},
		{
			name:    "Nil image",
			img:     nil,
			mime:    MIME_TYPE_PNG,
			wantErr: false, // Function returns nil, nil for nil image
		},
		{
			name:    "Unsupported mime type",
			img:     testImg,
			mime:    "unsupported/mime",
			wantErr: true,
		},
		{
			name:    "GIF mime type (unsupported)",
			img:     testImg,
			mime:    MIME_TYPE_GIF,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetBytes4mImage(tt.img, tt.mime)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetBytes4mImage() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("GetBytes4mImage() unexpected error = %v", err)
				return
			}
			if tt.img == nil {
				if result != nil {
					t.Errorf("GetBytes4mImage() with nil image should return nil")
				}
			} else {
				if result == nil {
					t.Errorf("GetBytes4mImage() returned nil bytes for valid image")
				}
				if len(result) == 0 {
					t.Errorf("GetBytes4mImage() returned empty bytes for valid image")
				}
			}
		})
	}
}

func TestGetCoreImage(t *testing.T) {
	pngBytes := createPNGBytes()
	jpegBytes := createJPEGBytes()

	tests := []struct {
		name      string
		dataBytes []byte
		mime      string
		wantErr   bool
	}{
		{
			name:      "Valid PNG bytes",
			dataBytes: pngBytes,
			mime:      MIME_TYPE_PNG,
			wantErr:   false,
		},
		{
			name:      "Valid JPEG bytes",
			dataBytes: jpegBytes,
			mime:      MIME_TYPE_JPEG,
			wantErr:   false,
		},
		{
			name:      "Empty data bytes",
			dataBytes: []byte{},
			mime:      MIME_TYPE_PNG,
			wantErr:   true,
		},
		{
			name:      "Unsupported mime type",
			dataBytes: pngBytes,
			mime:      "unsupported/mime",
			wantErr:   true,
		},
		{
			name:      "GIF mime type (unsupported)",
			dataBytes: createGIFBytes(),
			mime:      MIME_TYPE_GIF,
			wantErr:   true,
		},
		{
			name:      "Invalid PNG data",
			dataBytes: []byte("invalid png data"),
			mime:      MIME_TYPE_PNG,
			wantErr:   true,
		},
		{
			name:      "Invalid JPEG data",
			dataBytes: []byte("invalid jpeg data"),
			mime:      MIME_TYPE_JPEG,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetCoreImage(tt.dataBytes, tt.mime)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCoreImage() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("GetCoreImage() unexpected error = %v", err)
				return
			}
			if result == nil {
				t.Errorf("GetCoreImage() returned nil image for valid data")
			}
		})
	}
}

func TestPngToImage(t *testing.T) {
	validPngBytes := createPNGBytes()

	tests := []struct {
		name      string
		dataBytes []byte
		wantErr   bool
	}{
		{
			name:      "Valid PNG bytes",
			dataBytes: validPngBytes,
			wantErr:   false,
		},
		{
			name:      "Invalid PNG data",
			dataBytes: []byte("invalid png data"),
			wantErr:   true,
		},
		{
			name:      "Empty data",
			dataBytes: []byte{},
			wantErr:   true,
		},
		{
			name:      "JPEG data (wrong format)",
			dataBytes: createJPEGBytes(),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PngToImage(tt.dataBytes)
			if tt.wantErr {
				if err == nil {
					t.Errorf("PngToImage() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("PngToImage() unexpected error = %v", err)
				return
			}
			if result == nil {
				t.Errorf("PngToImage() returned nil image for valid PNG data")
			}
		})
	}
}

func TestImageToBytesPng(t *testing.T) {
	testImg := createTestImage()

	tests := []struct {
		name    string
		img     image.Image
		wantErr bool
	}{
		{
			name:    "Valid image",
			img:     testImg,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ImageToBytesPng(tt.img)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ImageToBytesPng() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("ImageToBytesPng() unexpected error = %v", err)
				return
			}
			if len(result) == 0 {
				t.Errorf("ImageToBytesPng() returned empty bytes for valid image")
			}

			// Verify it's actually PNG data
			if !bytes.HasPrefix(result, []byte("\x89PNG\r\n\x1a\n")) {
				t.Errorf("ImageToBytesPng() did not return valid PNG data")
			}
		})
	}
}

func TestJpegToImage(t *testing.T) {
	validJpegBytes := createJPEGBytes()

	tests := []struct {
		name      string
		dataBytes []byte
		wantErr   bool
	}{
		{
			name:      "Valid JPEG bytes",
			dataBytes: validJpegBytes,
			wantErr:   false,
		},
		{
			name:      "Invalid JPEG data",
			dataBytes: []byte("invalid jpeg data"),
			wantErr:   true,
		},
		{
			name:      "Empty data",
			dataBytes: []byte{},
			wantErr:   true,
		},
		{
			name:      "PNG data (wrong format)",
			dataBytes: createPNGBytes(),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JpegToImage(tt.dataBytes)
			if tt.wantErr {
				if err == nil {
					t.Errorf("JpegToImage() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("JpegToImage() unexpected error = %v", err)
				return
			}
			if result == nil {
				t.Errorf("JpegToImage() returned nil image for valid JPEG data")
			}
		})
	}
}

func TestImageToBytesJpeg(t *testing.T) {
	testImg := createTestImage()

	tests := []struct {
		name    string
		img     image.Image
		wantErr bool
	}{
		{
			name:    "Valid image",
			img:     testImg,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ImageToBytesJpeg(tt.img)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ImageToBytesJpeg() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("ImageToBytesJpeg() unexpected error = %v", err)
				return
			}
			if len(result) == 0 {
				t.Errorf("ImageToBytesJpeg() returned empty bytes for valid image")
			}

			// Verify it's actually JPEG data
			if !bytes.HasPrefix(result, []byte("\xff\xd8\xff")) {
				t.Errorf("ImageToBytesJpeg() did not return valid JPEG data")
			}
		})
	}
}

// Integration tests
func TestImageConversionRoundTrip(t *testing.T) {
	originalImg := createTestImage()

	// Test PNG round trip
	t.Run("PNG round trip", func(t *testing.T) {
		pngBytes, err := ImageToBytesPng(originalImg)
		if err != nil {
			t.Fatalf("Failed to convert image to PNG: %v", err)
		}

		recoveredImg, err := PngToImage(pngBytes)
		if err != nil {
			t.Fatalf("Failed to convert PNG back to image: %v", err)
		}

		if recoveredImg.Bounds() != originalImg.Bounds() {
			t.Errorf("Image bounds changed during PNG round trip")
		}
	})

	// Test JPEG round trip
	t.Run("JPEG round trip", func(t *testing.T) {
		jpegBytes, err := ImageToBytesJpeg(originalImg)
		if err != nil {
			t.Fatalf("Failed to convert image to JPEG: %v", err)
		}

		recoveredImg, err := JpegToImage(jpegBytes)
		if err != nil {
			t.Fatalf("Failed to convert JPEG back to image: %v", err)
		}

		if recoveredImg.Bounds() != originalImg.Bounds() {
			t.Errorf("Image bounds changed during JPEG round trip")
		}
	})
}

// Benchmark tests
func BenchmarkGetMime(b *testing.B) {
	pngData := createPNGBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetMime(pngData)
	}
}

func BenchmarkImageToBytesPng(b *testing.B) {
	img := createTestImage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ImageToBytesPng(img)
	}
}

func BenchmarkImageToBytesJpeg(b *testing.B) {
	img := createTestImage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ImageToBytesJpeg(img)
	}
}

func BenchmarkPngToImage(b *testing.B) {
	pngBytes := createPNGBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PngToImage(pngBytes)
	}
}

func BenchmarkJpegToImage(b *testing.B) {
	jpegBytes := createJPEGBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JpegToImage(jpegBytes)
	}
}
