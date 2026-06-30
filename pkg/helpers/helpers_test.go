package helpers

import (
	"path/filepath"
	"testing"
)

func TestCheckFormInput(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  bool
	}{
		{name: "Positive value is valid", input: 5, want: true},
		{name: "Value of one is valid", input: 1, want: true},
		{name: "Zero is invalid", input: 0, want: false},
		{name: "Negative value is invalid", input: -3, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckFormInput(tt.input)
			if got != tt.want {
				t.Errorf("CheckFormInput(%d) = %t; want %t", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	name := ".env"
	got, err := LoadEnv(name)
	if err != nil {
		t.Fatalf("LoadEnv(%q) returned unexpected error: %v", name, err)
	}

	if !filepath.IsAbs(got) {
		t.Errorf("LoadEnv(%q) = %q; expected an absolute path", name, got)
	}

	if filepath.Base(got) != name {
		t.Errorf("LoadEnv(%q) base = %q; want %q", name, filepath.Base(got), name)
	}
}
