package main

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMockFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "sample.go")
	content := `package sample

import (
	"io/fs"
	"os"
)

type IExample interface {
	ReadFile(name string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
	NoOp()
	Multi(a int, b string) (bool, error)
	DirEntries(f *os.File, n int) ([]fs.DirEntry, error)
}
`
	if err := os.WriteFile(inputPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	outputPath, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
		Alias:         "Example",
		OutputPath:    filepath.Join(tmpDir, "outdir"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(outputPath, "/example.go") && !strings.HasSuffix(outputPath, "\\example.go") {
		t.Fatalf("expected generated file name to end with example.go, got %s", outputPath)
	}

	generated, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	text := string(generated)
	if !strings.Contains(text, "type MockExample struct") {
		t.Fatal("generated code missing MockExample definition")
	}
	if !strings.Contains(text, "type mockOnCallExample struct") {
		t.Fatal("generated code missing mockOnCallExample definition")
	}
	if !strings.Contains(text, "type mockSetReturnExample struct") {
		t.Fatal("generated code missing mockSetReturnExample definition")
	}
	if !strings.Contains(text, "func NewMockExample() *MockExample") {
		t.Fatal("generated code missing NewMockExample")
	}
	if !strings.Contains(text, "func (m *MockExample) ReadFile(name string) ([]byte, error)") {
		t.Fatal("generated code missing ReadFile method")
	}
	if !strings.Contains(text, "func (r *mockSetReturnExample) Stat(") {
		t.Fatal("generated code missing Stat SetReturn method")
	}
}

func TestParseSourceFile_ParsesImports(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "sample.go")
	content := `package sample

import (
	ioAlias "io"
	"os"
)
`
	if err := os.WriteFile(inputPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	file, imports, err := parseSourceFile(inputPath, []byte(content))
	if err != nil {
		t.Fatal(err)
	}
	if file.Name.Name != "sample" {
		t.Fatalf("expected package name sample, got %s", file.Name.Name)
	}
	if got, ok := imports["ioAlias"]; !ok || got != "io" {
		t.Fatalf("expected import alias ioAlias => io, got %v", imports)
	}
	if got, ok := imports["os"]; !ok || got != "os" {
		t.Fatalf("expected import alias os => os, got %v", imports)
	}
}

func TestParseSourceFile_ParseFailure(t *testing.T) {
	tmpDir := t.TempDir()
	content := []byte("package sample\nimport (unquoted)")
	_, _, err := parseSourceFile(filepath.Join(tmpDir, "bad.go"), content)
	if err == nil {
		t.Fatal("expected parse error for invalid import syntax")
	}
}

func TestParseSourceFile_InvalidRawStringImportPath(t *testing.T) {
	tmpDir := t.TempDir()
	content := []byte("package sample\nimport `io`\n")
	_, _, err := parseSourceFile(filepath.Join(tmpDir, "raw.go"), content)
	if err == nil || !strings.Contains(err.Error(), "invalid import path") {
		t.Fatalf("expected invalid import path error for raw string import, got %v", err)
	}
}

func TestStrconvUnquote_Error(t *testing.T) {
	_, err := strconvUnquote("notquoted")
	if err == nil || !strings.Contains(err.Error(), "invalid import path") {
		t.Fatalf("expected invalid import path error, got %v", err)
	}
}

func TestGenerateMockFile_MissingInputFile(t *testing.T) {
	_, err := GenerateMockFile(GeneratorOptions{
		InterfaceName: "IExample",
	})
	if err == nil || !strings.Contains(err.Error(), "input file is required") {
		t.Fatalf("expected input file required error, got %v", err)
	}
}

func TestGenerateMockFile_MissingInterfaceName(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "sample.go")
	if err := os.WriteFile(inputPath, []byte("package sample\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := GenerateMockFile(GeneratorOptions{
		InputFile: inputPath,
	})
	if err == nil || !strings.Contains(err.Error(), "interface name is required") {
		t.Fatalf("expected interface name required error, got %v", err)
	}
}

func TestGenerateMockFile_ReadFileError(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     filepath.Join(tmpDir, "doesnotexist.go"),
		InterfaceName: "IExample",
	})
	if err == nil || !os.IsNotExist(err) {
		t.Fatalf("expected read file error, got %v", err)
	}
}

func TestGenerateMockFile_ParseError(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "bad.go")
	if err := os.WriteFile(inputPath, []byte("package sample\ntype IExample interface {"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
	})
	if err == nil || !strings.Contains(err.Error(), "expected") {
		t.Fatalf("expected parse error, got %v", err)
	}
}

func TestGenerateMockFile_InterfaceNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "sample.go")
	if err := os.WriteFile(inputPath, []byte("package sample\n\ntype IOther interface { NoOp() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
	})
	if err == nil || !strings.Contains(err.Error(), "interface IExample not found") {
		t.Fatalf("expected interface not found error, got %v", err)
	}
}

func TestFindInterface_TypeIsNotInterface(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "sample.go")
	content := []byte("package sample\n\ntype IExample struct { Name string }\n")
	if err := os.WriteFile(inputPath, content, 0o644); err != nil {
		t.Fatal(err)
	}

	file, _, err := parseSourceFile(inputPath, content)
	if err != nil {
		t.Fatal(err)
	}

	_, err = findInterface(file, "IExample")
	if err == nil || !strings.Contains(err.Error(), "type IExample is not an interface") {
		t.Fatalf("expected type is not an interface error, got %v", err)
	}
}

func TestBuildMethods_UnsupportedInterfaceMember(t *testing.T) {
	iface := &ast.InterfaceType{
		Methods: &ast.FieldList{List: []*ast.Field{{
			Names: []*ast.Ident{{Name: "Method"}},
			Type:  ast.NewIdent("int"),
		}}},
	}

	_, _, err := buildMethods(iface, map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "unsupported interface member Method") {
		t.Fatalf("expected unsupported interface member error, got %v", err)
	}
}

func TestBuildMethods_CollectImportsError(t *testing.T) {
	iface := &ast.InterfaceType{
		Methods: &ast.FieldList{List: []*ast.Field{{
			Names: []*ast.Ident{{Name: "Method"}},
			Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{{
				Type: &ast.SelectorExpr{X: ast.NewIdent("pkg"), Sel: ast.NewIdent("Type")},
			}}}},
		}}},
	}

	_, _, err := buildMethods(iface, map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "unknown selector package pkg") {
		t.Fatalf("expected unknown selector package error, got %v", err)
	}
}

func TestFormatParams_NilList(t *testing.T) {
	params, args := formatParams(nil)
	if params != "" || args != nil {
		t.Fatalf("expected empty params and nil args for nil list, got %q %v", params, args)
	}
}

func TestFormatParams_AnonymousParameter(t *testing.T) {
	list := &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent("string")}}}
	params, args := formatParams(list)
	if params != "string" {
		t.Fatalf("expected params string, got %q", params)
	}
	if len(args) != 1 || args[0] != "arg0" {
		t.Fatalf("expected anonymous arg arg0, got %v", args)
	}
}

func TestFormatResults_SingleResult(t *testing.T) {
	list := &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent("error")}}}
	results, types := formatResults(list)
	if results != "error" {
		t.Fatalf("expected single result string error, got %q", results)
	}
	if len(types) != 1 || types[0] != "error" {
		t.Fatalf("expected result types [error], got %v", types)
	}
}

func TestDeriveSetReturn_SingleResultName(t *testing.T) {
	signature, values := deriveSetReturn([]string{"string"})
	if signature != "result string" || values != "result" {
		t.Fatalf("expected result string / result, got %q / %q", signature, values)
	}

	signature, values = deriveSetReturn([]string{"error"})
	if signature != "err error" || values != "err" {
		t.Fatalf("expected err error / err, got %q / %q", signature, values)
	}
}

func TestRenderMock_ImportsWithAlias(t *testing.T) {
	code, err := renderMock("Example", nil, map[string]string{"ioAlias": "io"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(code), "ioAlias \"io\"") {
		t.Fatalf("expected aliased import line, got: %s", string(code))
	}
}

func TestResolveOutputPath_DefaultOutput(t *testing.T) {
	output, err := resolveOutputPath("", "Example")
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join("internal", "pkgxmock", "example.go")
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestResolveOutputPath_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	output, err := resolveOutputPath(tmpDir, "Example")
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(tmpDir, "example.go")
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestResolveOutputPath_StatErrorAfterExists(t *testing.T) {
	originalStat := osStat
	defer func() { osStat = originalStat }()

	callCount := 0
	osStat = func(name string) (os.FileInfo, error) {
		callCount++
		if callCount == 1 {
			return originalStat(name)
		}
		return nil, fmt.Errorf("stat failed")
	}

	tmpDir := t.TempDir()
	_, err := resolveOutputPath(tmpDir, "Example")
	if err == nil || !strings.Contains(err.Error(), "stat failed") {
		t.Fatalf("expected stat failed error, got %v", err)
	}
}

func TestPathExists_ExistingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	exists, err := pathExists(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("expected existing path to return true")
	}
}

func TestGenerateMockFile_EmbeddedInterfaceNotSupported(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "sample.go")
	content := `package sample

import "io"

type IExample interface {
	io.Closer
}
`
	if err := os.WriteFile(inputPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
	})
	if err == nil || !strings.Contains(err.Error(), "embedded interfaces are not supported") {
		t.Fatalf("expected embedded interface error, got %v", err)
	}
}

func TestGenerateMockFile_RenderMockError(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "sample.go")
	content := `package sample

type IExample interface {
	NoOp()
}
`
	if err := os.WriteFile(inputPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
		Alias:         "Bad-Alias",
	})
	if err == nil || !strings.Contains(err.Error(), "expected") {
		t.Fatalf("expected render error, got %v", err)
	}
}

func TestGenerateMockFile_ResolveOutputPathError(t *testing.T) {
	tmpDir := t.TempDir()
	symlinkPath := filepath.Join(tmpDir, "broken")
	if err := os.Symlink(filepath.Join(tmpDir, "missing"), symlinkPath); err != nil {
		t.Fatal(err)
	}

	inputPath := filepath.Join(tmpDir, "sample.go")
	if err := os.WriteFile(inputPath, []byte("package sample\ntype IExample interface { NoOp() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
		OutputPath:    symlinkPath,
	})
	if err == nil {
		t.Fatal("expected resolve output path error")
	}
}

func TestGenerateMockFile_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	badParent := filepath.Join(tmpDir, "badparent")
	if err := os.WriteFile(badParent, []byte("not a dir"), 0o644); err != nil {
		t.Fatal(err)
	}

	inputPath := filepath.Join(tmpDir, "sample.go")
	if err := os.WriteFile(inputPath, []byte("package sample\ntype IExample interface { NoOp() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	outputPath := filepath.Join(badParent, "example.go")
	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
		OutputPath:    outputPath,
	})
	if err == nil {
		t.Fatal("expected mkdirall error")
	}
}

func TestGenerateMockFile_WriteFileError(t *testing.T) {
	tmpDir := t.TempDir()

	inputPath := filepath.Join(tmpDir, "sample.go")
	if err := os.WriteFile(inputPath, []byte("package sample\ntype IExample interface { NoOp() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	outputPath := filepath.Join(tmpDir, "example.go")
	if err := os.WriteFile(outputPath, []byte("existing content"), 0o444); err != nil {
		t.Fatal(err)
	}

	_, err := GenerateMockFile(GeneratorOptions{
		InputFile:     inputPath,
		InterfaceName: "IExample",
		Alias:         "Example",
		OutputPath:    outputPath,
	})
	if err == nil {
		t.Fatal("expected write file error")
	}
}
