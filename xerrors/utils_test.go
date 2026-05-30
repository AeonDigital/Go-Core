package xerrors_test

import (
	"errors"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

func TestPrint_Resilience(t *testing.T) {
	// Cenário 1: Garante que passar um erro nil não causa pânico e retorna em silêncio
	t.Run("Should return early without panic when error is nil", func(t *testing.T) {
		xerrors.Print(nil)
	})

	// Cenário 2: Executa o fluxo feliz escrevendo o erro no Stderr
	t.Run("Should write error message string to stderr stream successfully", func(t *testing.T) {
		testErr := xerrors.NewErr("temporary infrastructure tracking failure")
		xerrors.Print(testErr)
	})
}

func TestMsgData_FormattingAndBoundaries(t *testing.T) {
	type customInt int
	var customZero customInt = 0

	tests := []struct {
		name        string
		inputArgs   []any
		expectedStr string
	}{
		{
			name:        "Should return empty string immediately when zero arguments are provided",
			inputArgs:   nil,
			expectedStr: "",
		},
		{
			name: "Should apply default filters ignoring pairs where value is nil, zero, or empty text",
			// Ajustado para conter aspas simples no valor booleano 'true'
			inputArgs:   []any{"user_id", 1042, "tenant", "", "retry", 0, "active", true, "token", nil},
			expectedStr: "user_id=1042, active='true'",
		},
		{
			name:        "Should intercept an odd number of inputs and apply the missing value safeguard block",
			inputArgs:   []any{"timeout", 30, "debug_mode"},
			expectedStr: "timeout=30, debug_mode=(MISSING_VALUE)",
		},
		{
			name:        "Should support custom blacklist overrides in the first argument allowing zeros to pass",
			inputArgs:   []any{[]any{nil}, "user_id", 1042, "retry", 0, "token", nil},
			expectedStr: "user_id=1042, retry=0",
		},
		{
			name:        "Should force deep reflection check by matching a customInt(0) option against a standard int(0) blacklist",
			inputArgs:   []any{[]any{0}, "user_id", 1042, "custom_metric", customZero},
			expectedStr: "user_id=1042",
		},
		{
			name:        "Should return empty string if the remaining arguments list after blacklist resolution is empty",
			inputArgs:   []any{[]any{nil}},
			expectedStr: "",
		},
		{
			name: "Should format literal nil values with quotes function when custom blacklist allows nil to pass",
			// Passamos a blacklist contendo apenas o número 999. O 'nil' passa direto e é formatado!
			inputArgs:   []any{[]any{999}, "user_id", 1042, "token", nil},
			expectedStr: "user_id=1042, token=nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := xerrors.MsgData(tt.inputArgs...)
			if gotResult != tt.expectedStr {
				t.Errorf("[%s] mismatch:\n got:  %q\n want: %q", tt.name, gotResult, tt.expectedStr)
			}
		})
	}
}

func TestMsgOptions_FormattingAndBoundaries(t *testing.T) {
	type dummyStruct struct{ Name string }
	type customInt int
	var customZero customInt = 0
	var customFortyTwo customInt = 42

	tests := []struct {
		name        string
		inputArgs   []any
		expectedStr string
	}{
		{
			name:        "Should return empty string immediately when zero arguments are provided",
			inputArgs:   nil,
			expectedStr: "",
		},
		{
			name:        "Should apply default filters ignoring nil, zero, and empty strings completely",
			inputArgs:   []any{"json", "", 0, "yaml", nil, "text"},
			expectedStr: "'json', 'yaml', 'text'",
		},
		{
			name:        "Should format numbers without single quotes and keep strings wrapped in single quotes",
			inputArgs:   []any{10, 25.5, "max_threshold"},
			expectedStr: "10, 25.5, 'max_threshold'",
		},
		{
			name:        "Should wrap boolean values in single quotes as standard primitives",
			inputArgs:   []any{true, false},
			expectedStr: "'true', 'false'",
		},
		{
			name:        "Should format complex structs using standard stringification layout inside single quotes",
			inputArgs:   []any{dummyStruct{Name: "core"}},
			expectedStr: "'{core}'",
		},
		{
			name:        "Should support custom blacklist overrides in the first argument allowing empty strings to pass",
			inputArgs:   []any{[]any{nil}, "json", "", "yaml"},
			expectedStr: "'json', '', 'yaml'",
		},
		{
			name:        "Should append literal nil string when custom blacklist allows nil to bypass verification",
			inputArgs:   []any{[]any{""}, "json", nil, "yaml"},
			expectedStr: "'json', nil, 'yaml'",
		},
		{
			name:        "Should force deep reflection check by matching a customInt(0) option against a standard int(0) blacklist",
			inputArgs:   []any{[]any{0}, customZero, customFortyTwo, "text"},
			expectedStr: "42, 'text'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := xerrors.MsgOptions(tt.inputArgs...)
			if gotResult != tt.expectedStr {
				t.Errorf("[%s] mismatch:\n got:  %q\n want: %q", tt.name, gotResult, tt.expectedStr)
			}
		})
	}
}

func TestMsgArraySize_FormattingAndBoundaries(t *testing.T) {
	tests := []struct {
		name        string
		inputArgs   []any
		expectedStr string
	}{
		{
			name:        "Should return empty string immediately when zero arguments are provided",
			inputArgs:   nil,
			expectedStr: "",
		},
		{
			name:        "Should build unified collection dimensions preserving numeric zero states safely",
			inputArgs:   []any{"parsers", 1, "options", 0, "plugins", 5},
			expectedStr: "parsers(1), options(0), plugins(5)",
		},
		{
			name:        "Should intercept an odd number of inputs and apply the missing size safeguard block",
			inputArgs:   []any{"parsers", 2, "options"},
			expectedStr: "parsers(2), options((MISSING_SIZE))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := xerrors.MsgArraySize(tt.inputArgs...)
			if gotResult != tt.expectedStr {
				t.Errorf("[%s] mismatch:\n got:  %q\n want: %q", tt.name, gotResult, tt.expectedStr)
			}
		})
	}
}

func TestNewErr_UniversalMaskProgressiveMatrix(t *testing.T) {
	// Inicializa e registra o token complexo de teste no mapa mestre de produção
	xerrors.RegisterTestTokenInRegistry()

	mask := xerrors.XERR_TEST_COMPLEX

	// Testamos o fallback tradicional do NewErr primeiro (mensagens normais sem token)
	t.Run("Standard Fallback: plain static message", func(t *testing.T) {
		got := xerrors.NewErr("plain generic system failure")
		if got.Error() != "plain generic system failure" {
			t.Errorf("unexpected output: %s", got.Error())
		}
	})

	t.Run("Standard Fallback: native formatting verbs", func(t *testing.T) {
		got := xerrors.NewErr("retry %d on %s", 2, "localhost")
		if got.Error() != "retry 2 on localhost" {
			t.Errorf("unexpected output: %s", got.Error())
		}
	})

	// Matriz Progressiva: Estressa todas as ramificações de loops e tamanhos de slices do motor mestre
	t.Run("Matrix 0: zero arguments supplied (triggers Early Return)", func(t *testing.T) {
		got := xerrors.NewErr(mask)
		if got.Error() != mask {
			t.Errorf("expected %q, got %q", mask, got.Error())
		}
	})

	t.Run("Matrix 1: only CTX provided (forces ø fallback on MSG and all extra tags)", func(t *testing.T) {
		got := xerrors.NewErr(mask, "INFRA")
		expected := "[CTX: INFRA][MSG: ø][TAG1: ø][TAG2: ø][TAG3: ø]::[ERR: ø]"
		if got.Error() != expected {
			t.Errorf("expected %q, got %q", expected, got.Error())
		}
	})

	t.Run("Matrix 2: CTX and empty string MSG provided (triggers defaultMsg normalization and loops tags to ø)", func(t *testing.T) {
		got := xerrors.NewErr(mask, "INFRA", "")
		expected := "[CTX: INFRA][MSG: test default message layout][TAG1: ø][TAG2: ø][TAG3: ø]::[ERR: ø]"
		if got.Error() != expected {
			t.Errorf("expected %q, got %q", expected, got.Error())
		}
	})

	t.Run("Matrix 3: CTX, MSG, and TAG1 provided (TAG2 and TAG3 fall back to ø)", func(t *testing.T) {
		got := xerrors.NewErr(mask, "INFRA", "custom operation", "val1")
		expected := "[CTX: INFRA][MSG: custom operation][TAG1: val1][TAG2: ø][TAG3: ø]::[ERR: ø]"
		if got.Error() != expected {
			t.Errorf("expected %q, got %q", expected, got.Error())
		}
	})

	t.Run("Matrix 4: CTX, MSG, TAG1, and TAG2 provided (TAG3 automatically intercepts defaultRule)", func(t *testing.T) {
		got := xerrors.NewErr(mask, "INFRA", "custom operation", "val1", "val2")
		expected := "[CTX: INFRA][MSG: custom operation][TAG1: val1][TAG2: val2][TAG3: test default rule requirement]::[ERR: ø]"
		if got.Error() != expected {
			t.Errorf("expected %q, got %q", expected, got.Error())
		}
	})

	t.Run("Matrix 5 (with native error): all fields filled + root cause (triggers %w wrapping branch)", func(t *testing.T) {
		rootErr := errors.New("timeout downstream network connection")
		got := xerrors.NewErr(mask, "INFRA", "custom operation", "val1", "val2", "custom_rule", rootErr)
		expected := "[CTX: INFRA][MSG: custom operation][TAG1: val1][TAG2: val2][TAG3: custom_rule]::[ERR: timeout downstream network connection]"
		if got.Error() != expected {
			t.Errorf("expected %q, got %q", expected, got.Error())
		}
	})

	t.Run("Matrix 5 (with empty rule): all fields filled but rule string is empty (triggers defaultRule normalization + primitive %v error fallback)", func(t *testing.T) {
		got := xerrors.NewErr(mask, "INFRA", "custom operation", "val1", "val2", "", "primitive raw data text exception context")
		expected := "[CTX: INFRA][MSG: custom operation][TAG1: val1][TAG2: val2][TAG3: test default rule requirement]::[ERR: primitive raw data text exception context]"
		if got.Error() != expected {
			t.Errorf("expected %q, got %q", expected, got.Error())
		}
	})
}
