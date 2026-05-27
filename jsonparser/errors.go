package jsonparser

import "errors"

var (
	ErrInvalidNumber      = errors.New("invalid number")
	ErrInvalidString      = errors.New("invalid string")
	ErrInvalidIdentifier  = errors.New("invalid identifier")
	ErrTrailingComma      = errors.New("illegal trailing comma")
	ErrUnknown            = errors.New("unknown error")
	ErrMissingKey         = errors.New("missing string key in object")
	ErrInvalidObjectValue = errors.New("invalid value in object")
	ErrMissingColon       = errors.New("missing colon after key in object")
	ErrInvalidArray       = errors.New("invalid array")
	ErrInvalidBool        = errors.New("invalid bool")
	ErrUnexpectedToken    = errors.New("unexpected token")
)
