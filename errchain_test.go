package errchain_test

import (
	"errors"
	"io"
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
				t.Fatalf("wrong error text: expected %q; got %q", c.expectedText, text)
			}
		})
	}
}

func TestIs(t *testing.T) {
	var (
		myErr = errors.New("my err")
		errs  = []error{
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

	err := errchain.New(errs...)

	for _, e := range errs {
		if errors.Is(err, e) {
			continue
		}

		t.Fatalf("unexpected result FALSE for error %s and %s", err, e)
	}

	for _, e := range unexpectedErrs {
		if !errors.Is(err, e) {
			continue
		}

		t.Fatalf("unexpected result TRUE for error %s and %s", err, e)
	}
}
