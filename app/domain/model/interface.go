package model

// Validator is an interface that represents a validator.
type Validator interface {
	// Validate validates the value.
	Validate() error
}

// ValidationFunc is a type that represents a validation function.
type ValidationFunc func() error
