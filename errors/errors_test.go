package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/haleyrc/pkg/errors"
)

func bareErrorMaker() error    { return fmt.Errorf("oops") }
func wrappedErrorMaker() error { return errors.WrapError(bareErrorMaker()) }

func TestWrappedErrorsReportTheirCaller(t *testing.T) {
	err := errors.WrapError(wrappedErrorMaker())
	got := err.Error()
	want := "errors_test.TestWrappedErrorsReportTheirCaller: errors_test.wrappedErrorMaker: oops"
	if got != want {
		t.Errorf("Expected error to be %q, but got %q.", want, got)
	}
}

func TestWrappedErrorsReportTheCodeOfTheirCause(t *testing.T) {
	want := http.StatusTeapot
	err := errors.WrapError(errors.NewRawError(want, "oops"))
	got := errors.ErrorCode(err)
	if got != want {
		t.Errorf("Expected code to be %d, but got %d.", want, got)
	}
}

func TestWrappedErrorsReportTheMessageOfTheirCause(t *testing.T) {
	want := "oops"
	err := errors.WrapError(errors.NewRawError(http.StatusTeapot, want))
	got := errors.ErrorMessage(err)
	if got != want {
		t.Errorf("Expected message to be %q, but got %q.", want, got)
	}
}
