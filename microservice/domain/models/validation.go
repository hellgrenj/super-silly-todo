package models

// Ok represents types capable of validating
// themselves.
type Ok interface {
	OK() error
}

// ErrMissingField is a custom validation error
type ErrMissingField string

func (e ErrMissingField) Error() string {
	return string(e) + " is required"
}
