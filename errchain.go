package errchain

import "errors"

type chain struct {
	error

	next error
}

// New builds a chain of errors.
func New(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	return chain{
		error: errs[0],
		next:  New(errs[1:]...),
	}
}

// Error returns string with concatenated underlying errors
// strings, nested in "(" and ")".
func (c chain) Error() string {
	if c.error == nil && c.next == nil {
		return ""
	}

	if c.error == nil {
		return c.next.Error()
	}

	if c.next == nil {
		return c.error.Error()
	}

	return c.error.Error() + " (" + c.next.Error() + ")"
}

// Is examines a chain for compliance with any error.
func (c chain) Is(target error) bool {
	return errors.Is(c.error, target) || errors.Is(c.next, target)
}

// As finds the first error in a chain that matches target, and
// if one is found, sets target to that error value and returns
// true. Otherwise, it returns false.
func (c chain) As(target any) bool {
	return errors.As(c.error, target) || errors.As(c.next, target)
}
