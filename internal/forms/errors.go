package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form strtuct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks if form field is in post and not empty (only mandatory fields)
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "El campo no puede estar vacío.")
		}
	}
}

// Has checks if form field is in post and not empty (checkboxes or radiobuttons)
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "El campo no puede estar vacío.")
		return false
	}
	return true
}

// MinLength checks for string minimun length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("El campo debe tener una logitud mínima de %d caracteres.", length))
		return false
	}
	return true
}

// Check is IsEmail with GoValidator package
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Correo electrónico no válido.")
	}
}
