package xerrors_test

import (
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

func TestError400_PolymorphicConstructor_Case1(t *testing.T) {
	// Case 1: Testing extended domain signature (PKGCTX + Code)
	// Using TestCtx and TestCode registered in our init block
	err := xerrors.NewError400(TestCtx, TestCode, "param1", "param2")

	if err.Code() != TestCode {
		t.Errorf("expected code to be %s, got %s", TestCode, err.Code())
	}

	// Dynamic tags formatting check
	output := err.Error()
	if !strings.Contains(output, "[TAG1: param1]") || !strings.Contains(output, "[TAG2: param2]") {
		t.Errorf("expected output to contain extracted dynamic tags, got:\n%s", output)
	}
}

func TestError400_PolymorphicConstructor_Case2(t *testing.T) {
	// Case 2: Testing core framework signature (only first code token provided)
	// XERR_FIELD_REQUIRED is a native constant inside pkg_xerrors.go
	err := xerrors.NewError400(xerrors.XERR_FIELD_REQUIRED, "email_field")

	if err.Code() != xerrors.XERR_FIELD_REQUIRED {
		t.Errorf("expected code to be %s, got %s", xerrors.XERR_FIELD_REQUIRED, err.Code())
	}

	output := err.Error()
	if !strings.Contains(output, "[FIELD: email_field]") {
		t.Errorf("expected output to contain core metadata FIELD tag, got:\n%s", output)
	}
}

func TestError400_PolymorphicConstructor_Case3(t *testing.T) {
	// Case 3: Fallback to standard text plain or Sprintf formatting string
	errSimple := xerrors.NewError400("simple unmapped plain error")
	if errSimple.Code() != xerrors.XERR_NONE {
		t.Errorf("expected code XERR_NONE for plain text, got %s", errSimple.Code())
	}

	errFormatted := xerrors.NewError400("invalid parameter %s with value %d", "age", 15)
	if !strings.Contains(errFormatted.Error(), "[MSG: invalid parameter age with value 15]") {
		t.Errorf("expected formatted string output, got:\n%s", errFormatted.Error())
	}
}

func TestError400_PolymorphicConstructor_EmptyAndSafeguards(t *testing.T) {
	// Scenario A: Absolutely no parameters passed to the variadic block
	errEmpty := xerrors.NewError400()
	if errEmpty == nil {
		t.Fatal("expected instance to be returned even with zero arguments")
	}

	// Scenario B: Unexpected data type sent into the first positioning slot
	// Triggers the fmt.Sprintf("%v", args) ultimate structural safeguard branch
	badSlice := []int{1, 2, 3}
	errSafeguard := xerrors.NewError400(badSlice)

	if !strings.Contains(errSafeguard.Error(), "[1 2 3]") {
		t.Errorf("expected unexpected type to be safely dumped as string, got:\n%s", errSafeguard.Error())
	}
}

func TestError400_WithArgs_EmptyBoundary(t *testing.T) {
	errBase := xerrors.NewError400("baseline validation message")

	// Invoking WithArgs with absolutely no parameters should return the exact same instance pointer
	// This branch protects against unnecessary memory allocation/cloning overhead
	errSame := errBase.WithArgs()

	if errBase != errSame {
		t.Error("expected WithArgs with no elements to bypass cloning and return identical instance pointer")
	}
}

func TestError400_WithArgs_ValidPayload(t *testing.T) {
	// Start with a basic plain text validation error
	errBase := xerrors.NewError400("initial validation failure")

	// Chain WithArgs to inject structured fields later in the process
	errWithFields := errBase.WithArgs("injected_param_val")

	// Verify that a new instance was generated and it formatted correctly
	output := errWithFields.Error()
	if output == "" {
		t.Error("expected populated output from the cloned instance")
	}
}
