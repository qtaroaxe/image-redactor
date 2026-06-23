package compression

import (
	"testing"
)

func TestNewQuality(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantVal int
		wantErr bool
	}{
		{"Zero input returns default", 0, 85, false},
		{"Valid min boundary", 1, 1, false},
		{"Valid mid value", 50, 50, false},
		{"Valid max boundary", 100, 100, false},
		{"Invalid negative value", -1, 0, true},
		{"Invalid too high value", 101, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuality(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQuality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.wantVal {
				t.Errorf("NewQuality() got = %v, want %v", got.Value(), tt.wantVal)
			}
		})
	}
}

func TestNewCompressionLevel(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantVal int
		wantErr bool
	}{
		{"Zero input returns default", 0, 6, false},
		{"Valid min boundary", 1, 1, false},
		{"Valid mid value", 5, 5, false},
		{"Valid max boundary", 10, 10, false},
		{"Invalid negative value", -5, 0, true},
		{"Invalid too high value", 11, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCompressionLevel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCompressionLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.wantVal {
				t.Errorf("NewCompressionLevel() got = %v, want %v", got.Value(), tt.wantVal)
			}
		})
	}
}
