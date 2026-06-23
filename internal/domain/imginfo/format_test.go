package imginfo_test

import (
	"fmt"
	"testing"

	"github.com/QtaroAXE/image-redactor/internal/domain/imginfo"
)

func TestNewFormat_Valid(t *testing.T) {
	cases := []struct {
		in   string
		want imginfo.Format
	}{
		{"jpeg", imginfo.FormatJPEG},
		{"JPEG", imginfo.FormatJPEG},
		{"png", imginfo.FormatPNG},
		{"  png", imginfo.FormatPNG},
		{"png  ", imginfo.FormatPNG},
		{"PNG", imginfo.FormatPNG},
		{"jpg", imginfo.FormatJPEG},
		{"JPG", imginfo.FormatJPEG},
		{"webp", imginfo.FormatWebP},
		{"WEBP", imginfo.FormatWebP},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := imginfo.NewFormat(tc.in)
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
			got, err := imginfo.NewFormat(in)
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
	if imginfo.FormatJPEG.String() != "jpeg" {
		t.Errorf("FormatJPEG.String() = %q, want %q", imginfo.FormatJPEG.String(), "jpeg")
	}
	var zero imginfo.Format
	if !zero.IsZero() {
		t.Errorf("zero Format should report IsZero() == true")
	}
	if imginfo.FormatPNG.IsZero() {
		t.Errorf("FormatPNG must not be zero")
	}
}
