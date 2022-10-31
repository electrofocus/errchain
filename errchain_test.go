package errchain_test

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"testing"

	"github.com/electrofocus/errchain"
)

func TestNewAndError(t *testing.T) {
	cases := []struct {
		name         string
		errs         []error
		expectNil    bool
		expectedText string
	}{
		{
			name:         "no errs",
			errs:         nil,
			expectNil:    true,
			expectedText: "",
		},
		{
			name:         "empty errs slice",
			errs:         make([]error, 0),
			expectNil:    true,
			expectedText: "",
		},
		{
			name: "some errs",
			errs: []error{
				errors.New("1"),
				errors.New("2"),
				errors.New("3"),
			},
			expectNil:    false,
			expectedText: "1 (2 (3))",
		},
		{
			name:         "one err",
			errs:         []error{errors.New("1")},
			expectNil:    false,
			expectedText: "1",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := errchain.New(c.errs...)

			if c.expectNil && err == nil {
				return
			}

			if err == nil {
				t.Fatal("error is nil, but must not be nil")
			}

			if text := err.Error(); text != c.expectedText {
				t.Fatalf("wrong error text: expected %q, but got %q", c.expectedText, text)
			}
		})
	}
}

func TestIs(t *testing.T) {
	var (
		myErr        = errors.New("my err")
		expectedErrs = []error{
			errchain.New(errchain.New(errors.New("nested one"))),
			myErr,
			io.EOF,
			os.ErrClosed,
			http.ErrHijacked,
		}
		unexpectedErrs = []error{
			io.ErrClosedPipe,
			os.ErrDeadlineExceeded,
			http.ErrAbortHandler,
			nil,
		}
	)

	chain := errchain.New(expectedErrs...)

	for _, e := range expectedErrs {
		if errors.Is(chain, e) {
			continue
		}

		t.Fatalf("unexpected result FALSE for error chain (%s) and (%s)", chain, e)
	}

	for _, e := range unexpectedErrs {
		if !errors.Is(chain, e) {
			continue
		}

		t.Fatalf("unexpected result TRUE for error chain (%s) and (%s)", chain, e)
	}

	if !errors.Is(chain, io.EOF) {
		t.Fatalf("unexpected result FALSE for error chain (%s) and (%s)", chain, io.EOF)
	}

	if errors.Is(chain, io.ErrClosedPipe) {
		t.Fatalf("unexpected result TRUE for error chain (%s) and (%s)", chain, io.ErrClosedPipe)
	}
}

type customErr struct {
	text string
}

func (e customErr) Error() string {
	return e.text
}

func TestAs(t *testing.T) {

	err := errchain.New(
		customErr{"first error text"},
		&fs.PathError{Op: "readdir", Path: "home", Err: errors.New("not implemented")},
		&json.SyntaxError{},
	)

	var (
		err1 customErr
		err2 *fs.PathError
		err3 *json.SyntaxError
		err4 *json.MarshalerError
	)

	if !errors.As(err, &err1) {
		t.Fatalf("unexpected result FALSE for error (%s) and error type %T", err, err1)
	}

	if !errors.As(err, &err2) {
		t.Fatalf("unexpected result FALSE for error (%s) and error type %T", err, err2)
	}

	if !errors.As(err, &err3) {
		t.Fatalf("unexpected result FALSE for error (%s) and error type %T", err, err3)
	}

	if errors.As(err, &err4) {
		t.Fatalf("unexpected result TRUE for error (%s) and error type %T", err, err4)
	}
}
