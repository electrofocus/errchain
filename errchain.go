package errchain

import "errors"

type chain struct {
	error

	next error
}

// New builds errors chain from errs.
// You can use errors.Is to check the chain for compliance with any error.
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

// Error returns string with concatenated underlying errors strings, separated by "; ".
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

// Is allows to check the chain for compliance with any error.
func (c chain) Is(target error) bool {
	return errors.Is(c.error, target) || errors.Is(c.next, target)
}
