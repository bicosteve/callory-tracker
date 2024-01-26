package forms

// This file contains errors type which will be used to
// hold validation error messages for forms
// the name of the form field will be used as key in the map

type errors map[string][]string

// Add adds error messages of the given field to the map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves first error message of a given field
func (e errors) Get(field string) string {
	errorString := e[field]
	if len(errorString) == 0 {
		return ""
	}
	return errorString[0]
}
