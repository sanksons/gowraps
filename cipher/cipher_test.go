package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"strings"
	"testing"
)

// Helper function to create a test key
func createTestKey() [32]byte {
	var key [32]byte
	copy(key[:], "this-is-a-32-byte-key-for-test!!")
	return key
}

// Helper function to create a random key
func createRandomKey() [32]byte {
	var key [32]byte
	rand.Read(key[:])
	return key
}

// Helper function for min (Go < 1.21 compatibility)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "Hello, World!",
			expected: "SGVsbG8sIFdvcmxkIQ==",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Special characters",
			input:    "Hello@#$%^&*()World!",
			expected: "SGVsbG9AIyQlXiYqKClXb3JsZCE=",
		},
		{
			name:     "Unicode characters",
			input:    "Hello 世界",
			expected: "SGVsbG8g5LiW55WM",
		},
		{
			name:     "Long text",
			input:    strings.Repeat("A", 50),
			expected: "QUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUE=",
		},
		{
			name:     "Newlines and spaces",
			input:    "Line 1\nLine 2\t\rLine 3",
			expected: "TGluZSAxCkxpbmUgMgkNTGluZSAz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Base64Encode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple text",
			input:    "SGVsbG8sIFdvcmxkIQ==",
			expected: "Hello, World!",
			wantErr:  false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "Special characters",
			input:    "SGVsbG9AIyQlXiYqKClXb3JsZCE=",
			expected: "Hello@#$%^&*()World!",
			wantErr:  false,
		},
		{
			name:     "Unicode characters",
			input:    "SGVsbG8g5LiW55WM",
			expected: "Hello 世界",
			wantErr:  false,
		},
		{
			name:     "Invalid base64",
			input:    "Invalid@Base64!",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Incomplete base64",
			input:    "SGVsbG9",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Base64 with spaces (invalid)",
			input:    "SGVs bG8g 5LiW 55WM",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Base64Decode(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Base64Decode() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Base64Decode() unexpected error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Base64Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBase64EncodeDecodeRoundTrip(t *testing.T) {
	testCases := []string{
		"Hello, World!",
		"",
		"Special chars: @#$%^&*()",
		"Unicode: 你好世界",
		strings.Repeat("Long text test ", 100),
		"Newlines\nand\ttabs\rand\rcarriage\rreturns",
	}

	for _, tc := range testCases {
		t.Run("RoundTrip_"+tc[:min(len(tc), 20)], func(t *testing.T) {
			encoded := Base64Encode(tc)
			decoded, err := Base64Decode(encoded)
			if err != nil {
				t.Errorf("Round trip failed on decode: %v", err)
			}
			if decoded != tc {
				t.Errorf("Round trip failed: got %v, want %v", decoded, tc)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	key := createTestKey()

	tests := []struct {
		name    string
		key     [32]byte
		text    []byte
		wantErr bool
	}{
		{
			name:    "Simple text",
			key:     key,
			text:    []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "Empty text",
			key:     key,
			text:    []byte(""),
			wantErr: false,
		},
		{
			name:    "Large text",
			key:     key,
			text:    bytes.Repeat([]byte("A"), 10000),
			wantErr: false,
		},
		{
			name:    "Binary data",
			key:     key,
			text:    []byte{0, 1, 2, 3, 255, 254, 253},
			wantErr: false,
		},
		{
			name:    "Unicode text",
			key:     key,
			text:    []byte("Hello 世界 🌍"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encrypt(tt.key, tt.text)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Encrypt() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Encrypt() unexpected error = %v", err)
				return
			}

			// Verify the result is not empty (unless input was empty)
			if len(tt.text) > 0 && len(result) == 0 {
				t.Errorf("Encrypt() returned empty result for non-empty input")
			}

			// Verify the result is longer than AES block size for non-empty input
			if len(tt.text) > 0 && len(result) <= aes.BlockSize {
				t.Errorf("Encrypt() result too short: got %d bytes, expected > %d", len(result), aes.BlockSize)
			}

			// Verify the result doesn't equal the input (unless empty)
			if len(tt.text) > 0 && bytes.Equal(result, tt.text) {
				t.Errorf("Encrypt() result equals input (no encryption occurred)")
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	key := createTestKey()

	// First, create some valid encrypted data
	plaintext := []byte("Hello, World!")
	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Failed to create test encrypted data: %v", err)
	}

	tests := []struct {
		name     string
		key      [32]byte
		text     []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Valid encrypted data",
			key:      key,
			text:     encrypted,
			expected: plaintext,
			wantErr:  false,
		},
		{
			name:     "Empty ciphertext",
			key:      key,
			text:     []byte{},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Ciphertext too short",
			key:      key,
			text:     []byte{1, 2, 3},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid ciphertext (wrong size)",
			key:      key,
			text:     make([]byte, aes.BlockSize+1),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Wrong key",
			key:      createRandomKey(),
			text:     encrypted,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decrypt(tt.key, tt.text)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Decrypt() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Decrypt() unexpected error = %v", err)
				return
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Decrypt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := createTestKey()

	testCases := [][]byte{
		[]byte("Hello, World!"),
		[]byte(""),
		[]byte("Special chars: @#$%^&*()"),
		[]byte("Unicode: 你好世界 🌍"),
		bytes.Repeat([]byte("Long text test "), 1000),
		[]byte{0, 1, 2, 3, 255, 254, 253}, // Binary data
		[]byte("Newlines\nand\ttabs\rand\rcarriage\rreturns"),
	}

	for i, tc := range testCases {
		t.Run("RoundTrip_"+string(rune('A'+i)), func(t *testing.T) {
			// Encrypt
			encrypted, err := Encrypt(key, tc)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			// Decrypt
			decrypted, err := Decrypt(key, encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			// Compare
			if !bytes.Equal(decrypted, tc) {
				t.Errorf("Round trip failed: got %v, want %v", decrypted, tc)
			}
		})
	}
}

func TestEncryptWithDifferentKeys(t *testing.T) {
	key1 := createTestKey()
	key2 := createRandomKey()
	plaintext := []byte("Secret message")

	// Encrypt with key1
	encrypted1, err := Encrypt(key1, plaintext)
	if err != nil {
		t.Fatalf("Encrypt with key1 failed: %v", err)
	}

	// Encrypt with key2
	encrypted2, err := Encrypt(key2, plaintext)
	if err != nil {
		t.Fatalf("Encrypt with key2 failed: %v", err)
	}

	// Results should be different
	if bytes.Equal(encrypted1, encrypted2) {
		t.Errorf("Encryption with different keys produced identical results")
	}

	// Each key should only decrypt its own data
	decrypted1, err := Decrypt(key1, encrypted1)
	if err != nil {
		t.Errorf("Decrypt with correct key1 failed: %v", err)
	}
	if !bytes.Equal(decrypted1, plaintext) {
		t.Errorf("Decryption with key1 failed: got %v, want %v", decrypted1, plaintext)
	}

	// Wrong key should fail
	_, err = Decrypt(key2, encrypted1)
	if err == nil {
		t.Errorf("Decrypt with wrong key should have failed")
	}
}

// Benchmark tests
func BenchmarkBase64Encode(b *testing.B) {
	data := "Hello, World! This is a benchmark test for Base64 encoding."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Base64Encode(data)
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	data := Base64Encode("Hello, World! This is a benchmark test for Base64 decoding.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Base64Decode(data)
	}
}

func BenchmarkEncrypt(b *testing.B) {
	key := createTestKey()
	data := []byte("Hello, World! This is a benchmark test for encryption.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encrypt(key, data)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	key := createTestKey()
	data := []byte("Hello, World! This is a benchmark test for decryption.")
	encrypted, _ := Encrypt(key, data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decrypt(key, encrypted)
	}
}

func BenchmarkEncryptDecryptRoundTrip(b *testing.B) {
	key := createTestKey()
	data := []byte("Hello, World! This is a benchmark test for round trip encryption.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, _ := Encrypt(key, data)
		Decrypt(key, encrypted)
	}
}
