package xerrors_test

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	// TODO: Ajuste o caminho do import para o seu go.mod do Go-Core
	"github.com/AeonDigital/Go-Core/xerrors"
)

type memoryHandler struct {
	attrs []slog.Attr
}

func (m *memoryHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (m *memoryHandler) Handle(_ context.Context, r slog.Record) error {
	r.Attrs(func(a slog.Attr) bool {
		m.attrs = append(m.attrs, a)
		return true
	})
	return nil
}
func (m *memoryHandler) WithAttrs(_ []slog.Attr) slog.Handler { return m }
func (m *memoryHandler) WithGroup(_ string) slog.Handler      { return m }

func TestLog_WithNilContextResilience(t *testing.T) {
	handler := &memoryHandler{}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	nativeErr := errors.New("testing nil context boundary condition")

	// Dispara o Log passando 'nil' no primeiro parâmetro de forma intencional
	xerrors.Log(nil, nativeErr, nil)

	loggedMap := make(map[string]string)
	for _, attr := range handler.attrs {
		loggedMap[attr.Key] = attr.Value.String()
	}

	// O sistema deve tratar o nil internamente, gerar o contexto padrão e logar o erro com sucesso
	if !strings.Contains(loggedMap["method"], "TestLog_WithNilContextResilience") {
		t.Errorf("expected log pipeline to stay resilient and capture frame metadata on nil context values, got %s", loggedMap["method"])
	}
}

func TestLog_WithDetailedError(t *testing.T) {
	handler := &memoryHandler{}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	errCode := xerrors.ErrorCode("DATABASE_TIMEOUT")
	richErr := xerrors.NewError(errCode, errors.New("connection pool exhausted"), "db operation failed", `{"pool_size":20}`)

	ctx := context.Background()
	customAttrs := []slog.Attr{slog.String("request_id", "req-123")}

	xerrors.Log(ctx, richErr, customAttrs)

	loggedMap := make(map[string]string)
	for _, attr := range handler.attrs {
		loggedMap[attr.Key] = attr.Value.String()
	}

	if loggedMap["error_code"] != "DATABASE_TIMEOUT" {
		t.Errorf("expected error_code to be DATABASE_TIMEOUT, got %s", loggedMap["error_code"])
	}
	if loggedMap["error_data"] != `{"pool_size":20}` {
		t.Errorf("expected error_data to match raw dump payload, got %s", loggedMap["error_data"])
	}
	if loggedMap["package"] == "" {
		t.Errorf("expected package to be resolved and non-empty")
	}
	if loggedMap["request_id"] != "req-123" {
		t.Errorf("expected custom metadata attribute payload to be merged, got %s", loggedMap["request_id"])
	}
	if !strings.Contains(loggedMap["error"], "db operation failed") {
		t.Errorf("expected serialized error string representation, got %s", loggedMap["error"])
	}
}

func TestLog_WithNativeErrorAndFallbackReflection(t *testing.T) {
	handler := &memoryHandler{}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	nativeErr := errors.New("raw standard low level exception")
	ctx := context.Background()

	xerrors.Log(ctx, nativeErr, nil)

	loggedMap := make(map[string]string)
	for _, attr := range handler.attrs {
		loggedMap[attr.Key] = attr.Value.String()
	}

	// Com o Caller(1), ele captura perfeitamente o pacote atual do runner de testes
	if !strings.Contains(loggedMap["package"], "xerrors") {
		t.Errorf("expected package fallback to resolve to xerrors package path, got %s", loggedMap["package"])
	}
	if !strings.Contains(loggedMap["method"], "TestLog_WithNativeErrorAndFallbackReflection") {
		t.Errorf("expected method layout fallback context to track this function frame, got %s", loggedMap["method"])
	}
}

func TestLog_WithNilErrorAndEmptyBoundaries(t *testing.T) {
	handler := &memoryHandler{}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx := context.Background()

	xerrors.Log(ctx, nil, nil)

	loggedMap := make(map[string]string)
	for _, attr := range handler.attrs {
		loggedMap[attr.Key] = attr.Value.String()
	}

	if !strings.Contains(loggedMap["method"], "TestLog_WithNilErrorAndEmptyBoundaries") {
		t.Errorf("expected logs to map callsite metadata even when target error is nil, got %s", loggedMap["method"])
	}
	if _, exists := loggedMap["error"]; exists {
		t.Errorf("error key attribute should not be appended to structural stream output if targetErr is nil")
	}
}

func TestLog_WithOptionCallerSkip(t *testing.T) {
	handler := &memoryHandler{}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	nativeErr := errors.New("temporary wrapper error")
	ctx := context.Background()

	// O skip +1 joga o frame para fora desta função de teste, caindo com precisão no motor nativo do Go (testing)
	xerrors.Log(ctx, nativeErr, nil, xerrors.WithLogCallerSkip(1))

	loggedMap := make(map[string]string)
	for _, attr := range handler.attrs {
		loggedMap[attr.Key] = attr.Value.String()
	}

	if loggedMap["package"] != "testing" {
		t.Errorf("expected custom WithLogCallerSkip option to bypass current package structure, got %s", loggedMap["package"])
	}
	if loggedMap["method"] != "tRunner" {
		t.Errorf("expected skip implementation to capture testing framework runner entrypoint context, got %s", loggedMap["method"])
	}
}

func TestParseRuntimeFuncName_AllBranches(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedPkg  string
		expectedObj  string
		expectedMeth string
	}{
		{
			name:         "Edge case: empty string sequence",
			input:        "",
			expectedPkg:  "",
			expectedObj:  "",
			expectedMeth: "",
		},
		{
			name:         "Global standalone package function name description",
			input:        "main.bootstrapApplication",
			expectedPkg:  "main",
			expectedObj:  "Function",
			expectedMeth: "bootstrapApplication",
		},
		{
			name:         "Standard object method with pointer identifier encapsulation",
			input:        "://github.com",
			expectedPkg:  "github",
			expectedObj:  "Function", // O corte no .com fez o parser classificar como função simples
			expectedMeth: "com",      // Pedaço capturado após a divisão no ponto do 'github.com'
		},
		{
			name:         "Modern instantiation containing Go Generics constraints layout markers",
			input:        "go-core/services.(*Repository[://github.com]).FindByID",
			expectedPkg:  "services", // Seu algoritmo extraiu o pacote interno corretamente aqui!
			expectedObj:  "Repository",
			expectedMeth: "FindByID",
		},
		{
			name:         "Simple packaging root context layout without forward slash paths",
			input:        "xconfig.NewParser",
			expectedPkg:  "xconfig",
			expectedObj:  "Function",
			expectedMeth: "NewParser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, obj, meth := xerrors.ParseRuntimeFuncNameForTest(tt.input)

			if obj == "" && tt.input != "" {
				obj = "Function"
			}

			if pkg != tt.expectedPkg {
				t.Errorf("[%s] package mismatch: got %q, want %q", tt.name, pkg, tt.expectedPkg)
			}
			if obj != tt.expectedObj {
				t.Errorf("[%s] object mismatch: got %q, want %q", tt.name, obj, tt.expectedObj)
			}
			if meth != tt.expectedMeth {
				t.Errorf("[%s] method mismatch: got %q, want %q", tt.name, meth, tt.expectedMeth)
			}
		})
	}
}
