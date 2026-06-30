package forms

import "testing"

func TestErrorsAdd(t *testing.T) {
	e := errors{}

	e.Add("email", "This field is required")
	e.Add("email", "Invalid email address")

	got := e["email"]
	if len(got) != 2 {
		t.Errorf("expected 2 messages for field 'email'; got %d", len(got))
	}

	if got[0] != "This field is required" {
		t.Errorf("expected first message %q; got %q", "This field is required", got[0])
	}

	if got[1] != "Invalid email address" {
		t.Errorf("expected second message %q; got %q", "Invalid email address", got[1])
	}
}

func TestErrorsGet(t *testing.T) {
	tests := []struct {
		name   string
		errors errors
		field  string
		want   string
	}{
		{
			name:   "Field with a single error",
			errors: errors{"email": []string{"This field is required"}},
			field:  "email",
			want:   "This field is required",
		},
		{
			name:   "Field with multiple errors returns the first",
			errors: errors{"email": []string{"first error", "second error"}},
			field:  "email",
			want:   "first error",
		},
		{
			name:   "Field with no errors returns empty string",
			errors: errors{},
			field:  "email",
			want:   "",
		},
		{
			name:   "Unknown field returns empty string",
			errors: errors{"name": []string{"some error"}},
			field:  "email",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.errors.Get(tt.field)
			if got != tt.want {
				t.Errorf("Get(%q) = %q; want %q", tt.field, got, tt.want)
			}
		})
	}
}
