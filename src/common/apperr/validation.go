package apperr

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// FromValidation converts an ozzo-validation error into *Error.
// A nil input returns nil.
func FromValidation(err error) *Error {
	if err == nil {
		return nil
	}

	details := map[string]string{}
	if validationErrors, ok := err.(validation.Errors); ok {
		for field, fieldErr := range validationErrors {
			details[field] = fieldErr.Error()
		}
		return Validation(details)
	}

	return Internal(err)
}
