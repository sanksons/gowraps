package timer

import (
	"strings"
	"testing"
	"time"
)

func TestGetHttpTime(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "Unix epoch time",
			input:    time.Unix(0, 0).UTC(),
			expected: "Thu, 01 Jan 1970 00:00:00 UTC",
		},
		{
			name:     "Specific date and time",
			input:    time.Date(2023, 7, 15, 14, 30, 45, 0, time.UTC),
			expected: "Sat, 15 Jul 2023 14:30:45 UTC",
		},
		{
			name:     "New Year's Day 2024",
			input:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "Mon, 01 Jan 2024 00:00:00 UTC",
		},
		{
			name:     "Leap year date",
			input:    time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: "Thu, 29 Feb 2024 12:00:00 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetHttpTime(tt.input)
			if result != tt.expected {
				t.Errorf("GetHttpTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCurrentTime(t *testing.T) {
	tests := []struct {
		name    string
		zone    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "UTC timezone",
			zone:    "UTC",
			wantErr: false,
		},
		{
			name:    "Empty string defaults to UTC",
			zone:    "",
			wantErr: false,
		},
		{
			name:    "New York timezone",
			zone:    "America/New_York",
			wantErr: false,
		},
		{
			name:    "Tokyo timezone",
			zone:    "Asia/Tokyo",
			wantErr: false,
		},
		{
			name:    "London timezone",
			zone:    "Europe/London",
			wantErr: false,
		},
		{
			name:    "Invalid timezone",
			zone:    "Invalid/Timezone",
			wantErr: true,
			errMsg:  "unknown time zone",
		},
		{
			name:    "Another invalid timezone",
			zone:    "BadZone",
			wantErr: true,
			errMsg:  "unknown time zone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetCurrentTime(tt.zone)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCurrentTime() expected error but got none")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("GetCurrentTime() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("GetCurrentTime() unexpected error = %v", err)
				return
			}

			// Verify the time is recent (within last minute)
			now := time.Now()
			if result.After(now.Add(time.Minute)) || result.Before(now.Add(-time.Minute)) {
				t.Errorf("GetCurrentTime() returned time %v is not recent (current: %v)", result, now)
			}

			// Verify timezone is correct
			expectedZone := tt.zone
			if expectedZone == "" {
				expectedZone = "UTC"
			}

			location, _ := time.LoadLocation(expectedZone)
			expectedTime := now.In(location)

			if result.Location().String() != expectedTime.Location().String() {
				t.Errorf("GetCurrentTime() timezone = %v, want %v", result.Location(), expectedTime.Location())
			}
		})
	}
}

func TestGetCurrentTimeTimezone(t *testing.T) {
	// Test that different timezones return different times
	utcTime, err := GetCurrentTime("UTC")
	if err != nil {
		t.Fatalf("Failed to get UTC time: %v", err)
	}

	nyTime, err := GetCurrentTime("America/New_York")
	if err != nil {
		t.Fatalf("Failed to get New York time: %v", err)
	}

	// The times should be the same instant but in different locations
	if !utcTime.Equal(nyTime) {
		t.Errorf("UTC and NY times should represent the same instant: UTC=%v, NY=%v", utcTime, nyTime)
	}

	// But their string representations should be different (different timezones)
	if utcTime.String() == nyTime.String() {
		t.Errorf("UTC and NY time strings should be different due to timezone")
	}
}

func TestGetCurrentTimeDefaultZone(t *testing.T) {
	// Test that empty string defaults to UTC
	utcTime, err := GetCurrentTime("UTC")
	if err != nil {
		t.Fatalf("Failed to get UTC time: %v", err)
	}

	defaultTime, err := GetCurrentTime("")
	if err != nil {
		t.Fatalf("Failed to get default time: %v", err)
	}

	// Both should be in UTC timezone
	if utcTime.Location().String() != defaultTime.Location().String() {
		t.Errorf("Default timezone should be UTC: got %v, want %v",
			defaultTime.Location(), utcTime.Location())
	}
}

// Benchmark tests
func BenchmarkGetHttpTime(b *testing.B) {
	testTime := time.Date(2023, 7, 15, 14, 30, 45, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetHttpTime(testTime)
	}
}

func BenchmarkGetCurrentTime(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCurrentTime("UTC")
	}
}

func BenchmarkGetCurrentTimeWithTimezone(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCurrentTime("America/New_York")
	}
}
