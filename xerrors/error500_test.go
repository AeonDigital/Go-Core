package xerrors_test

import (
	"errors"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

func TestError500_InterfaceGettersAndUnwrap(t *testing.T) {
	rootErr := errors.New("database connection refused")
	ctxToken := xerrors.ErrorCode("ERR_INFRA")
	codeToken := xerrors.ErrorCode("E5001")
	messageText := "internal server instability"
	infoText := "host: cluster-01, port: 5432"

	// 1. Initialize standard instance
	err := xerrors.NewError500(ctxToken, codeToken, rootErr, messageText, infoText)

	// 2. Validate all exposed interface contract getters
	if err.CTX() != ctxToken {
		t.Errorf("expected CTX %s, got %s", ctxToken, err.CTX())
	}

	if err.Code() != codeToken {
		t.Errorf("expected Code %s, got %s", codeToken, err.Code())
	}

	if err.Message() != messageText {
		t.Errorf("expected Message %s, got %s", messageText, err.Message())
	}

	if err.Info() != infoText {
		t.Errorf("expected Info %s, got %s", infoText, err.Info())
	}

	// 3. Validate native Unwrap alignment (errors.Is / errors.As support)
	unwrapped := errors.Unwrap(err)
	if unwrapped != rootErr {
		t.Errorf("expected unwrapped error to be %v, got %v", rootErr, unwrapped)
	}
}

func TestError500_WithArgs_EmptyBoundary(t *testing.T) {
	errBase := xerrors.NewError500("CTX", "CODE", nil, "msg", "")

	// 1. Invoking WithArgs with absolutely no parameters should return the exact same instance pointer
	// This branch protects against unnecessary memory allocation/cloning overhead
	errSame := errBase.WithArgs()

	if errBase != errSame {
		t.Error("expected WithArgs with no elements to bypass cloning and return identical instance pointer")
	}
}

func TestError500_ErrorFormattingMethod(t *testing.T) {
	err := xerrors.NewError500("CTX_SERVER", "E5002", nil, "unexpected crash", "")

	// Validates basic native output rendering string mapping behavior
	output := err.Error()
	if output == "" {
		t.Error("expected Error() method string generation to be populated, got empty string")
	}
}
