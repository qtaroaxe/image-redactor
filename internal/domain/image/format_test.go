package image_test

import (
	"fmt"
	"testing"

	"github.com/QtaroAXE/image-redactor/internal/domain/image"
)

func TestNewFormat_Valid(t *testing.T) {
	cases := []struct {
		in   string
		want image.Format
	}{
		{"jpeg", image.FormatJPEG},
		{"JPEG", image.FormatJPEG},
		{"png", image.FormatPNG},
		{"  png", image.FormatPNG},
		{"png  ", image.FormatPNG},
		{"PNG", image.FormatPNG},
		{"jpg", image.FormatJPEG},
		{"JPG", image.FormatJPEG},
		{"webp", image.FormatWebP},
		{"WEBP", image.FormatWebP},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := image.NewFormat(tc.in)
			if err != nil {
				t.Fatalf("NewFormat(%q) returned error: %v", tc.in, err)
			} else {
				fmt.Println("Format succesfuly checked")
			}
			if !got.Equals(tc.want.String()) {
				t.Errorf("NewFormat(%q) = %v, want %v", tc.in, got, tc.want)
			} else {
				fmt.Printf("Format %q is valid\n", tc.in)
			}
		})
	}
}

func TestNewFormat_Invalid(t *testing.T) {
	cases := []string{"", " ", "bmp", "gif", "tuff", "something"}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := image.NewFormat(in)
			if err != nil {
				t.Fatalf("NewFormat(%q) = %v, want error", in, got)
			}
			if !got.IsZero() {
				t.Errorf("NewFormat(%q) returned non-zero Format on error: %v", in, got)
			}
		})
	}
}

func TestFormat_StringAndIsZero(t *testing.T) {
	if image.FormatJPEG.String() != "jpeg" {
		t.Errorf("FormatJPEG.String() = %q, want %q", image.FormatJPEG.String(), "jpeg")
	}
	var zero image.Format
	if !zero.IsZero() {
		t.Errorf("zero Format should report IsZero() == true")
	}
	if image.FormatPNG.IsZero() {
		t.Errorf("FormatPNG must not be zero")
	}
}
