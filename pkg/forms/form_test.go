package forms

import (
	"net/url"
	"regexp"
	"testing"
)

func TestNewForm(t *testing.T) {
	data := url.Values{}
	data.Set("name", "Bico")
	form := NewForm(data)

	if form.Get("name") != "Bico" {
		t.Errorf("expected Get(%q) to return %q; got %q", "name", "Bico", form.Get("name"))
	}

	if len(form.Errors) != 0 {
		t.Errorf("expected form.Errors to be empty; got length %d", len(form.Errors))
	}
}

func TestFormRequired(t *testing.T) {
	tests := []struct {
		name       string
		formData   url.Values
		fields     []string
		wantValid  bool
		wantErrors map[string]string
	}{
		{
			name: "All fields present",
			formData: url.Values{
				"name":  []string{"Apple"},
				"sugar": []string{"10"},
			},
			fields:    []string{"name", "sugar"},
			wantValid: true,
		},
		{
			name: "Missing required fields",
			formData: url.Values{
				"name": []string{"Apple"},
			},
			fields:    []string{"name", "sugar"},
			wantValid: false,
			wantErrors: map[string]string{
				"sugar": "This field is required",
			},
		},
		{
			name: "Blank required fields",
			formData: url.Values{
				"name":  []string{"   "},
				"sugar": []string{"10"},
			},
			fields:    []string{"name", "sugar"},
			wantValid: false,
			wantErrors: map[string]string{
				"name": "This field is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := NewForm(tt.formData)
			form.Required(tt.fields...)

			if form.Valid() != tt.wantValid {
				t.Errorf("expected form valid: %t; got %t", tt.wantValid, form.Valid())
			}

			for field, errMsg := range tt.wantErrors {
				if form.Errors.Get(field) != errMsg {
					t.Errorf("expected error for %s to be %q; got %q", field, errMsg, form.Errors.Get(field))
				}
			}
		})
	}
}

func TestFormMaxLength(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		maxLen    int
		wantValid bool
	}{
		{
			name:      "Under max length",
			field:     "name",
			value:     "Apple",
			maxLen:    10,
			wantValid: true,
		},
		{
			name:      "Equal to max length",
			field:     "name",
			value:     "Apple",
			maxLen:    5,
			wantValid: true,
		},
		{
			name:      "Over max length",
			field:     "name",
			value:     "Apple Pie",
			maxLen:    5,
			wantValid: false,
		},
		{
			name:      "Empty field is valid (ignored)",
			field:     "name",
			value:     "",
			maxLen:    5,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set(tt.field, tt.value)
			form := NewForm(formData)
			form.MaxLength(tt.field, tt.maxLen)

			if form.Valid() != tt.wantValid {
				t.Errorf("expected form valid: %t; got %t", tt.wantValid, form.Valid())
			}
		})
	}
}

func TestFormMinLength(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		minLen    int
		wantValid bool
	}{
		{
			name:      "Over min length",
			field:     "password",
			value:     "secret123",
			minLen:    6,
			wantValid: true,
		},
		{
			name:      "Equal to min length",
			field:     "password",
			value:     "secret",
			minLen:    6,
			wantValid: true,
		},
		{
			name:      "Under min length",
			field:     "password",
			value:     "sec",
			minLen:    6,
			wantValid: false,
		},
		{
			name:      "Empty field is valid (ignored)",
			field:     "password",
			value:     "",
			minLen:    6,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set(tt.field, tt.value)
			form := NewForm(formData)
			form.MinLength(tt.field, tt.minLen)

			if form.Valid() != tt.wantValid {
				t.Errorf("expected form valid: %t; got %t", tt.wantValid, form.Valid())
			}
		})
	}
}

func TestFormMinValue(t *testing.T) {
	tests := []struct {
		name      string
		fields    []string
		formData  url.Values
		minVal    int
		wantValid bool
	}{
		{
			name:   "Value above min",
			fields: []string{"protein"},
			formData: url.Values{
				"protein": []string{"15"},
			},
			minVal:    1,
			wantValid: true,
		},
		{
			name:   "Value equal to min",
			fields: []string{"protein"},
			formData: url.Values{
				"protein": []string{"1"},
			},
			minVal:    1,
			wantValid: true,
		},
		{
			name:   "Value below min",
			fields: []string{"protein"},
			formData: url.Values{
				"protein": []string{"0"},
			},
			minVal:    1,
			wantValid: false,
		},
		{
			name:   "Empty value is ignored",
			fields: []string{"protein"},
			formData: url.Values{
				"protein": []string{""},
			},
			minVal:    1,
			wantValid: true,
		},
		{
			name:   "Invalid integer is treated as 0",
			fields: []string{"protein"},
			formData: url.Values{
				"protein": []string{"invalid"},
			},
			minVal:    1,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := NewForm(tt.formData)
			form.MinValue(tt.minVal, tt.fields...)

			if form.Valid() != tt.wantValid {
				t.Errorf("expected form valid: %t; got %t", tt.wantValid, form.Valid())
			}
		})
	}
}

func TestFormAllowedValues(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		allowed   []string
		wantValid bool
	}{
		{
			name:      "Value in allowed list",
			field:     "meal",
			value:     "breakfast",
			allowed:   []string{"breakfast", "lunch", "dinner"},
			wantValid: true,
		},
		{
			name:      "Value not in allowed list",
			field:     "meal",
			value:     "snack",
			allowed:   []string{"breakfast", "lunch", "dinner"},
			wantValid: false,
		},
		{
			name:      "Empty value is ignored",
			field:     "meal",
			value:     "",
			allowed:   []string{"breakfast", "lunch", "dinner"},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set(tt.field, tt.value)
			form := NewForm(formData)
			form.AllowedValues(tt.field, tt.allowed...)

			if form.Valid() != tt.wantValid {
				t.Errorf("expected form valid: %t; got %t", tt.wantValid, form.Valid())
			}
		})
	}
}

func TestFormComparePasswords(t *testing.T) {
	tests := []struct {
		name      string
		pass      string
		confirm   string
		wantValid bool
	}{
		{
			name:      "Passwords match",
			pass:      "secret123",
			confirm:   "secret123",
			wantValid: true,
		},
		{
			name:      "Passwords do not match",
			pass:      "secret123",
			confirm:   "different",
			wantValid: false,
		},
		{
			name:      "Empty password is ignored",
			pass:      "",
			confirm:   "secret123",
			wantValid: true,
		},
		{
			name:      "Empty confirm is ignored",
			pass:      "secret123",
			confirm:   "",
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("password", tt.pass)
			formData.Set("confirm", tt.confirm)
			form := NewForm(formData)
			form.ComparePasswords("password", "confirm")

			if form.Valid() != tt.wantValid {
				t.Errorf("expected form valid: %t; got %t", tt.wantValid, form.Valid())
			}
		})
	}
}

func TestFormValidateEmail(t *testing.T) {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	tests := []struct {
		name      string
		field     string
		value     string
		pattern   *regexp.Regexp
		wantValid bool
	}{
		{
			name:      "Valid email",
			field:     "email",
			value:     "test@example.com",
			pattern:   re,
			wantValid: true,
		},
		{
			name:      "Invalid email",
			field:     "email",
			value:     "invalid-email",
			pattern:   re,
			wantValid: false,
		},
		{
			name:      "Empty email is ignored",
			field:     "email",
			value:     "",
			pattern:   re,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set(tt.field, tt.value)
			form := NewForm(formData)
			form.ValidateEmail(tt.field, tt.pattern)

			if form.Valid() != tt.wantValid {
				t.Errorf("expected form valid: %t; got %t", tt.wantValid, form.Valid())
			}
		})
	}
}
