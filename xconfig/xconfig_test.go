package xconfig_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xconfig"
)

// MockParser implementa a interface Parser para simular retornos nos testes de motor
type MockParser struct {
	opts    xconfig.Options
	readFn  func() (map[string]any, error)
	errOpts error
}

func (m *MockParser) SetOptions(opts xconfig.Options) error {
	m.opts = opts
	return m.errOpts
}

func (m *MockParser) Read() (map[string]any, error) {
	if m.readFn != nil {
		return m.readFn()
	}
	return make(map[string]any), nil
}

func TestInitAppConfig_AsymmetricSizes(t *testing.T) {
	parsers := []xconfig.Parser{&MockParser{}}
	options := []xconfig.Options{} // Zero size, will fail validation

	cfg, err := xconfig.InitAppConfig(parsers, options)
	if err == nil {
		t.Fatalf("expected error due to asymmetric input sizes, got nil")
	}
	if cfg != nil {
		t.Errorf("expected return config pointer to be nil on critical validation failure")
	}

	// Technical Note: Asserts compliance with the structured corporate layout format
	expectedSubStr := "[CTX: ERR_XCONFIG][ERR: E3008]"
	if !strings.Contains(err.Error(), expectedSubStr) {
		t.Errorf("unexpected corporate error payload shape: %q", err.Error())
	}
}
func TestInitAppConfig_EmptySlices(t *testing.T) {
	var parsers []xconfig.Parser
	var options []xconfig.Options

	cfg, err := xconfig.InitAppConfig(parsers, options)
	if err != nil {
		t.Fatalf("expected no error on empty configuration bootstrapping, got: %v", err)
	}
	if cfg == nil {
		t.Fatalf("expected a valid initialized pointer to Config, got nil")
	}
	if len(cfg.Keys()) != 0 {
		t.Errorf("expected a clean configuration state with 0 keys, got %d", len(cfg.Keys()))
	}
}

func TestInitAppConfig_FullSuccessPipeline(t *testing.T) {
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{"app.port": 8080}, nil
		},
	}

	parsers := []xconfig.Parser{mock}
	options := []xconfig.Options{{Prefix: "TEST"}}

	cfg, err := xconfig.InitAppConfig(parsers, options)
	if err != nil {
		t.Fatalf("expected bootstrap chain to run seamlessly, got error: %v", err)
	}
	if cfg == nil {
		t.Fatalf("returned config instance pointer is nil")
	}

	val, exists := cfg.Get("app.port")
	if !exists || val != 8080 {
		t.Errorf("expected key 'app.port' to be populated with 8080, got %v (exists: %t)", val, exists)
	}
}

func TestInitAppConfig_RegisterFailure(t *testing.T) {
	mockErr := errors.New("mock register error")
	mock := &MockParser{
		errOpts: mockErr,
	}

	parsers := []xconfig.Parser{mock}
	options := []xconfig.Options{{Prefix: "FAIL"}}

	cfg, err := xconfig.InitAppConfig(parsers, options)
	if err == nil {
		t.Fatalf("expected error during registration phase, got nil")
	}
	if cfg != nil {
		t.Errorf("expected config pointer to be nil when registration fails")
	}
	if !errors.Is(err, mockErr) {
		t.Errorf("expected root cause error %v, got %v", mockErr, err)
	}
}

func TestInitAppConfig_LoadFailure(t *testing.T) {
	mockErr := errors.New("mock read error")
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return nil, mockErr
		},
	}

	parsers := []xconfig.Parser{mock}
	options := []xconfig.Options{{Prefix: "OK"}}

	cfg, err := xconfig.InitAppConfig(parsers, options)
	if err == nil {
		t.Fatalf("expected error during loading phase, got nil")
	}
	if cfg != nil {
		t.Errorf("expected config pointer to be nil when loading fails")
	}
	if !errors.Is(err, mockErr) {
		t.Errorf("expected root cause error %v, got %v", mockErr, err)
	}
}

func TestConfig_Register_Success(t *testing.T) {
	cfg := xconfig.NewConfig()
	mock := &MockParser{}
	opts := xconfig.Options{Prefix: "APP"}

	err := cfg.Register(mock, opts)
	if err != nil {
		t.Fatalf("expected register to succeed, got: %v", err)
	}
}

func TestConfig_Register_ParserFailure(t *testing.T) {
	mockErr := errors.New("invalid configuration injection")
	mock := &MockParser{
		errOpts: mockErr,
	}
	opts := xconfig.Options{Prefix: "INVALID"}

	cfg := xconfig.NewConfig()
	err := cfg.Register(mock, opts)
	if err == nil {
		t.Fatalf("expected registration error from parser validation, got nil")
	}
	if !errors.Is(err, mockErr) {
		t.Errorf("expected error %v, got %v", mockErr, err)
	}
}

func TestConfig_Load_SuccessAndOverridePriority(t *testing.T) {
	cfg := xconfig.NewConfig()

	// Primeiro parser com chaves em caixa alta e espaços extras
	mock1 := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{
				"  DATABASE_HOST  ": "localhost",
				"database_port":     3306,
			}, nil
		},
	}

	// Segundo parser que deve sobrescrever a porta do banco de dados (prioridade linear)
	mock2 := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{
				"DATABASE_PORT": 5432,
				"app.env":       "production",
			}, nil
		},
	}

	_ = cfg.Register(mock1, xconfig.Options{})
	_ = cfg.Register(mock2, xconfig.Options{})

	err := cfg.Load()
	if err != nil {
		t.Fatalf("expected load to clear state and populate data cleanly, got: %v", err)
	}

	// Valida se limpou os espaços e jogou tudo para minusculo
	host, hostExists := cfg.Get("database_host")
	if !hostExists || host != "localhost" {
		t.Errorf("expected clean key 'database_host' to be 'localhost', got %v", host)
	}

	// Valida se o segundo parser ganhou do primeiro de forma linear
	port, portExists := cfg.Get("database_port")
	if !portExists || port != 5432 {
		t.Errorf("expected 'database_port' to be overridden by 5432, got %v", port)
	}

	env, envExists := cfg.Get("app.env")
	if !envExists || env != "production" {
		t.Errorf("expected 'app.env' to be 'production', got %v", env)
	}
}

func TestConfig_Load_ReadFailure(t *testing.T) {
	cfg := xconfig.NewConfig()
	mockErr := errors.New("disk read error failure")

	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return nil, mockErr
		},
	}

	_ = cfg.Register(mock, xconfig.Options{})

	err := cfg.Load()
	if err == nil {
		t.Fatalf("expected loading pipeline to fail when a parser returns an error, got nil")
	}
	if !errors.Is(err, mockErr) {
		t.Errorf("expected root cause error %v, got %v", mockErr, err)
	}
}

func TestConfig_Reload_Success(t *testing.T) {
	cfg := xconfig.NewConfig()
	counter := 0

	// Mock dinâmico para simular uma mudança de valor entre execuções
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			counter++
			return map[string]any{"execution.count": counter}, nil
		},
	}

	_ = cfg.Register(mock, xconfig.Options{})

	// Primeiro carregamento (Load)
	if err := cfg.Load(); err != nil {
		t.Fatalf("initial load failed: %v", err)
	}

	// Executa o recarregamento (Reload)
	if err := cfg.Reload(); err != nil {
		t.Fatalf("expected reload to re-execute parsers cleanly, got: %v", err)
	}

	val, exists := cfg.Get("execution.count")
	if !exists || val != 2 {
		t.Errorf("expected 'execution.count' to be re-read and incremented to 2, got %v", val)
	}
}

func TestConfig_Keys_Success(t *testing.T) {
	cfg := xconfig.NewConfig()
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{
				"app.name": "go-core",
				"app.env":  "production",
			}, nil
		},
	}

	_ = cfg.Register(mock, xconfig.Options{})
	_ = cfg.Load()

	keys := cfg.Keys()

	// Valida se o tamanho do slice de chaves coincide com o total de dados carregados
	if len(keys) != 2 {
		t.Fatalf("expected exactly 2 keys to be returned, got %d", len(keys))
	}

	// Cria um mapa auxiliar para garantir a presença dos itens sem depender de ordem fixa do loop
	keysMap := make(map[string]bool)
	for _, k := range keys {
		keysMap[k] = true
	}

	if !keysMap["app.name"] {
		t.Errorf("expected keys list to contain 'app.name'")
	}
	if !keysMap["app.env"] {
		t.Errorf("expected keys list to contain 'app.env'")
	}
}

func TestConfig_Has_Behavior(t *testing.T) {
	cfg := xconfig.NewConfig()
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{
				"database.name": "production_db",
			}, nil
		},
	}

	_ = cfg.Register(mock, xconfig.Options{})
	_ = cfg.Load()

	// Cenário 1: Sucesso - Chave existe exatamente como foi salva
	if !cfg.Has("database.name") {
		t.Errorf("expected Has to return true for exact key 'database.name'")
	}

	// Cenário 2: Sucesso - Chave informada com caixa alta e espaços (deve normalizar e achar)
	if !cfg.Has("  DATABASE.NAME  ") {
		t.Errorf("expected Has to normalize search parameter and return true for messy key input")
	}

	// Cenário 3: Falha - Chave não cadastrada no sistema
	if cfg.Has("non.existent.key") {
		t.Errorf("expected Has to return false for missing configuration keys")
	}
}

func TestConfig_Get_Behavior(t *testing.T) {
	cfg := xconfig.NewConfig()
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{
				"server.timeout": 30,
			}, nil
		},
	}

	_ = cfg.Register(mock, xconfig.Options{})
	_ = cfg.Load()

	// Cenário 1: Sucesso - Busca exata da chave cadastrada
	val, exists := cfg.Get("server.timeout")
	if !exists {
		t.Fatalf("expected key 'server.timeout' to exist")
	}
	if val != 30 {
		t.Errorf("expected value to be 30, got %v", val)
	}

	// Cenário 2: Sucesso - Busca utilizando chaves com espaçamento e letras maiúsculas
	valMessy, existsMessy := cfg.Get("  SERVER.TIMEOUT  ")
	if !existsMessy {
		t.Fatalf("expected key to be found after lookup normalization parameters")
	}
	if valMessy != 30 {
		t.Errorf("expected value from normalized lookup to be 30, got %v", valMessy)
	}

	// Cenário 3: Falha - Busca por uma chave inexistente no mapa global
	_, existsMissing := cfg.Get("server.host")
	if existsMissing {
		t.Errorf("expected exists boolean to be false for missing keys")
	}
}

func TestConfig_Populate_ValidationErrors(t *testing.T) {
	cfg := xconfig.NewConfig()

	type DummyStruct struct {
		Port int `json:"port"`
	}

	// Cenário 1: Falha - Passando a struct por valor (não é ponteiro)
	var instanceValue DummyStruct
	err := cfg.Populate(instanceValue)
	if err == nil {
		t.Errorf("expected error when target is not a pointer, got nil")
	}

	// Cenário 2: Falha - Passando um ponteiro nulo (nil)
	var instanceNil *DummyStruct = nil
	err = cfg.Populate(instanceNil)
	if err == nil {
		t.Errorf("expected error when target is a nil pointer, got nil")
	}

	// Cenário 3: Falha - Passando um ponteiro para um tipo básico (não é struct)
	var targetString string
	err = cfg.Populate(&targetString)
	if err == nil {
		t.Errorf("expected error when target does not point to a struct, got nil")
	}
}

func TestConfig_Populate_Success(t *testing.T) {
	cfg := xconfig.NewConfig()
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{
				"app_name": "api-core",
				"port":     9000,
				"debug":    true,
			}, nil
		},
	}

	_ = cfg.Register(mock, xconfig.Options{})
	_ = cfg.Load()

	// Struct de destino que receberá os dados via tags json
	type AppConfig struct {
		AppName string `json:"app_name"`
		Port    int    `json:"port"`
		Debug   bool   `json:"debug"`
	}

	var result AppConfig
	err := cfg.Populate(&result)
	if err != nil {
		t.Fatalf("expected populate to map fields cleanly, got error: %v", err)
	}

	// Valida se as propriedades foram preenchidas e convertidas corretamente pelo encoder
	if result.AppName != "api-core" {
		t.Errorf("expected AppName to be 'api-core', got %q", result.AppName)
	}
	if result.Port != 9000 {
		t.Errorf("expected Port to be 9000, got %d", result.Port)
	}
	if !result.Debug {
		t.Errorf("expected Debug to be true, got %t", result.Debug)
	}
}

func TestConfig_Populate_MarshalFailure(t *testing.T) {
	cfg := xconfig.NewConfig()

	// Funções (func) não podem ser serializadas em JSON e forçam o Marshal a falhar
	mock := &MockParser{
		readFn: func() (map[string]any, error) {
			return map[string]any{
				"invalid_json_field": func() {},
			}, nil
		},
	}

	_ = cfg.Register(mock, xconfig.Options{})
	_ = cfg.Load()

	type DummyStruct struct {
		Field string `json:"invalid_json_field"`
	}

	var result DummyStruct
	err := cfg.Populate(&result)
	if err == nil {
		t.Fatalf("expected json marshal error due to unsupported type (func), got nil")
	}
}
