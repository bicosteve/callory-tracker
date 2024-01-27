package forms

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Custom Form struct which anonymously embeds  url.Values object
// which holds the form data and an Errors field that hold validation
// errors for form data.

type Form struct {
	url.Values
	Errors errors
}

//	NewForm function that initializes a custom Form struct.
//
// Takes form data as the parameter
func NewForm(formData url.Values) *Form {
	return &Form{formData, errors(map[string][]string{})}
}

// Required method to check that specific fields in the form
// data are present and not blank.
// If any field fails this check, add appropriate message to the form errors
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}

// MaxLength checks that specific field in the form contains maximum number
// of characters. If the check fails, then add the appropriate error message
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (max %d)", d))
	}
}

// MinValue checks if the field input is less than the required value
// If check fails, add appropriate error message
func (f *Form) MinValue(d int, fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if value == "" {
			return
		}

		r, _ := strconv.Atoi(value)
		if r < d {
			f.Errors.Add(field, fmt.Sprintf("This field is too low (min %d)", d))
		}
	}

}

// AllowedValues checks if a specific field in the form matches set of specific
func (f *Form) AllowedValues(field string, options ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, option := range options {
		if value == option {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) ComparePasswords(password, confirmPassword string) {
	valueOne := f.Get(password)
	if valueOne == "" {
		return
	}

	valueTwo := f.Get(confirmPassword)
	if valueTwo == "" {
		return
	}

	if valueOne != valueTwo {
		f.Errors.Add(confirmPassword, "Confirm password does not match password")
	}

}

func (f *Form) ValidateEmail(email string) error {
	return nil
}

// Valid check if the form is valid
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
