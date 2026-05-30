package xdb_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xdb"
)

// 1. Criamos tipos que implementam as interfaces básicas do driver nativo do Go
type mockDriver struct{}

func (m *mockDriver) Open(name string) (driver.Conn, error) {
	return &mockConn{}, nil
}

type mockConn struct{}

func (m *mockConn) Prepare(query string) (driver.Stmt, error) { return &mockStmt{}, nil }
func (m *mockConn) Close() error                              { return nil }
func (m *mockConn) Begin() (driver.Tx, error)                 { return nil, nil }

// Execer é a interface que o db.Exec chama internamente
func (m *mockConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return &mockResult{}, nil // Sucesso absoluto em todos os PRAGMAs!
}

// Ping é a interface que o db.PingContext chama internamente.
// Aqui nós forçamos o erro de propósito!
func (m *mockConn) Ping(ctx context.Context) error {
	return fmt.Errorf("erro forcado de ping para cobertura")
}

type mockStmt struct{}

func (m *mockStmt) Close() error                                    { return nil }
func (m *mockStmt) NumInput() int                                   { return 0 }
func (m *mockStmt) Exec(args []driver.Value) (driver.Result, error) { return &mockResult{}, nil }
func (m *mockStmt) Query(args []driver.Value) (driver.Rows, error)  { return nil, nil }

type mockResult struct{}

func (m *mockResult) LastInsertId() (int64, error) { return 0, nil }
func (m *mockResult) RowsAffected() (int64, error) { return 0, nil }

// 2. Registramos o nosso driver mock com um nome customizado no init do arquivo de testes
func init() {
	sql.Register("sqlite_mock_ping_error", &mockDriver{})
}

// TestNewDBConfig valida a inicialização e os valores default.
func TestNewDBConfig(t *testing.T) {
	driver := "sqlite"
	dsn := ""
	migrationsDir := "./migrations"

	cfg := xdb.NewDBConfig(driver, dsn, migrationsDir, "", "", "", "")

	if cfg.Driver != "sqlite" {
		t.Errorf("esperava driver 'sqlite', recebeu '%s'", cfg.Driver)
	}
	if cfg.MaxOpenConnections != 1 {
		t.Errorf("esperava MaxOpenConnections 1, recebeu %d", cfg.MaxOpenConnections)
	}
	if cfg.SQLite.Pragma["journal_mode"] != "WAL" {
		t.Errorf("esperava pragma journal_mode 'WAL', recebeu '%s'", cfg.SQLite.Pragma["journal_mode"])
	}
}

// TestCheckConfiguration cobre os diferentes cenários de criação da string de conexão (antigo buildDSN).
func TestCheckConfiguration(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		inputCfg    xdb.DBConfig
		expectedDSN string
	}{
		{
			name: "DSN preexistente nao deve ser alterada",
			inputCfg: xdb.DBConfig{
				DSN: "file:custom.db?cache=shared",
				SQLite: xdb.SQLiteConfig{
					Mode: "file:",
				},
			},
			expectedDSN: "file:custom.db?cache=shared",
		},
		{
			name: "Memoria pura",
			inputCfg: xdb.DBConfig{
				SQLite: xdb.SQLiteConfig{Mode: ":memory:"},
			},
			expectedDSN: ":memory:",
		},
		{
			name: "Memoria com query string",
			inputCfg: xdb.DBConfig{
				SQLite: xdb.SQLiteConfig{
					Mode:        ":memory:",
					QueryString: "cache=shared",
				},
			},
			expectedDSN: "file::memory:?cache=shared",
		},
		{
			name: "Arquivo em disco valido",
			inputCfg: xdb.DBConfig{
				SQLite: xdb.SQLiteConfig{
					Dir:      tmpDir,
					FileName: "local.db",
				},
			},
			expectedDSN: "file:" + filepath.ToSlash(filepath.Join(tmpDir, "local.db")),
		},
		{
			name: "Arquivo customizado com query string",
			inputCfg: xdb.DBConfig{
				SQLite: xdb.SQLiteConfig{
					Dir:         tmpDir,
					FileName:    "test.db",
					QueryString: "mode=ro",
				},
			},
			expectedDSN: "file:" + filepath.ToSlash(filepath.Join(tmpDir, "test.db")) + "?mode=ro",
		},
		{
			name: "Memoria com file::memory: e query string",
			inputCfg: xdb.DBConfig{
				SQLite: xdb.SQLiteConfig{
					Mode:        "file::memory:",
					QueryString: "cache=shared",
				},
			},
			expectedDSN: "file::memory:?cache=shared",
		},
		{
			name: "Campos vazios nao geram DSN",
			inputCfg: xdb.DBConfig{
				SQLite: xdb.SQLiteConfig{
					Dir:      "",
					FileName: "",
				},
			},
			expectedDSN: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.inputCfg.CheckConfiguration()
			if err != nil {
				t.Fatalf("erro inesperado em CheckConfiguration: %v", err)
			}

			if tt.inputCfg.DSN != tt.expectedDSN {
				t.Errorf("\nesperava: %s\nrecebeu:  %s", tt.expectedDSN, tt.inputCfg.DSN)
			}

			// Se a DSN antiga foi informada, garante que limpou os campos do SQLite
			if strings.Contains(tt.name, "preexistente") && tt.inputCfg.SQLite.Mode != "" {
				t.Error("esperava que os campos SQLite fossem limpos se DSN já existisse")
			}
		})
	}
}

// TestInitDataBaseConnection testa o fluxo de sucesso e as ramificações de erro.
func TestInitDataBaseConnection(t *testing.T) {
	t.Run("Erro no sql.Open se o driver for invalido", func(t *testing.T) {
		cfg := xdb.DBConfig{
			Driver: "driver_invalido_para_forcar_erro",
			DSN:    ":memory:",
		}

		err := cfg.InitDataBaseConnection()
		if err == nil {
			t.Error("esperava erro ao abrir conexao com driver invalido")
		}
	})

	t.Run("Inicializar Pragma caso ele seja nil", func(t *testing.T) {
		cfg := xdb.DBConfig{
			Driver: "sqlite",
			DSN:    ":memory:",
			SQLite: xdb.SQLiteConfig{
				Pragma: nil, // Força a ramificacao 'if o.SQLite.Pragma == nil'
			},
		}

		err := cfg.InitDataBaseConnection()
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		defer cfg.DB.Close()

		if cfg.SQLite.Pragma == nil {
			t.Error("esperava que o mapa Pragma tivesse sido inicializado")
		}
		if cfg.SQLite.Pragma["journal_mode"] != "WAL" {
			t.Errorf("esperava defaults aplicados, recebeu %v", cfg.SQLite.Pragma["journal_mode"])
		}
	})

	t.Run("Erro ao executar query de Pragma invalido", func(t *testing.T) {
		cfg := xdb.DBConfig{
			Driver: "sqlite",
			DSN:    ":memory:",
			SQLite: xdb.SQLiteConfig{
				Pragma: map[string]string{
					// Injeta uma quebra de sintaxe SQL para forçar o erro no db.Exec
					"journal_mode": "WAL; CREATE TABLE );",
				},
			},
		}

		err := cfg.InitDataBaseConnection()
		if err == nil {
			if cfg.DB != nil {
				cfg.DB.Close()
			}
			t.Error("esperava erro de execucao do pragma invalido")
		}
	})

	t.Run("Erro no PingContext simulado por driver mock", func(t *testing.T) {
		// Passamos o nosso driver customizado que aceita todos os PRAGMAs,
		// mas falha miseravelmente apenas no handshake do Ping!
		cfg := xdb.DBConfig{
			Driver: "sqlite_mock_ping_error",
			DSN:    ":memory:",
			SQLite: xdb.SQLiteConfig{
				Pragma: map[string]string{
					"journal_mode": "WAL",
				},
			},
		}

		err := cfg.InitDataBaseConnection()
		if err == nil {
			if cfg.DB != nil {
				cfg.DB.Close()
			}
			t.Error("esperava erro no ping do banco")
		}

		if err != nil && !strings.Contains(err.Error(), "erro forcado de ping para cobertura") {
			t.Errorf("recebeu um erro diferente do esperado: %v", err)
		}
	})
}

// TestRunMigrations valida o fluxo de migrações em ordem alfabética e erros de diretório.
// TestRunMigrations valida todas as ramificações de erro e fluxos lógicos de migração.
func TestRunMigrations(t *testing.T) {
	// Base de sucesso comum para reaproveitar nos subtestes que precisam de um banco ativo
	newMockDB := func(t *testing.T) *sql.DB {
		db, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			t.Fatalf("falha ao abrir banco em memoria para teste: %v", err)
		}
		return db
	}

	t.Run("Erro se DB for nil", func(t *testing.T) {
		cfg := xdb.DBConfig{MigrationsDirPath: "./migrations"}
		err := cfg.RunMigrations()
		if err == nil || !strings.Contains(err.Error(), "database instance is nil") {
			t.Errorf("esperava erro de banco nil, recebeu: %v", err)
		}
	})

	t.Run("Erro se MigrationsDirPath for vazio", func(t *testing.T) {
		db := newMockDB(t)
		defer db.Close()

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: "   ", // Força o trim a resultar em string vazia
		}

		err := cfg.RunMigrations()
		if err == nil || !strings.Contains(err.Error(), "migrations directory path is empty") {
			t.Errorf("esperava erro de diretorio vazio, recebeu: %v", err)
		}
	})

	t.Run("Erro se o caminho fornecido for um arquivo e nao um diretorio", func(t *testing.T) {
		db := newMockDB(t)
		defer db.Close()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "dummy.txt")
		if err := os.WriteFile(filePath, []byte("conteudo"), 0644); err != nil {
			t.Fatal(err)
		}

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: filePath, // Passa um arquivo no lugar do diretorio
		}

		err := cfg.RunMigrations()
		if err == nil || !strings.Contains(err.Error(), "provided path is a file, not a directory") {
			t.Errorf("esperava erro de caminho sendo arquivo, recebeu: %v", err)
		}
	})

	t.Run("Erro se falhar ao ler o diretorio por falta de permissao", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Ignorando teste de permissão POSIX 0000 no Windows")
		}
		db := newMockDB(t)
		defer db.Close()

		tmpDir := t.TempDir()
		permDir := filepath.Join(tmpDir, "no_read_dir")
		if err := os.Mkdir(permDir, 0000); err != nil { // Sem nenhuma permissao de leitura (0000)
			t.Fatal(err)
		}
		defer func() { _ = os.Chmod(permDir, 0755) }() // Restaura para permitir limpeza automática

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: permDir,
		}

		err := cfg.RunMigrations()
		if err == nil || !strings.Contains(err.Error(), "failed to read migrations directory") {
			t.Errorf("esperava erro ao ler diretorio protegido, recebeu: %v", err)
		}
	})

	t.Run("Retorna nil se nao houver nenhum arquivo de migracao valido", func(t *testing.T) {
		db := newMockDB(t)
		defer db.Close()

		tmpDir := t.TempDir() // Diretorio completamente vazio, sem arquivos .sql

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: tmpDir,
		}

		err := cfg.RunMigrations()
		if err != nil {
			t.Errorf("nao esperava erro quando len(migrationFiles) == 0, recebeu: %v", err)
		}
	})

	t.Run("Erro se um arquivo SQL sumir ou falhar a leitura no meio do processo", func(t *testing.T) {
		db := newMockDB(t)
		defer db.Close()

		// O modernc.org/sqlite ou os.ReadDir não travam a leitura subsequente se removermos a permissão
		// do arquivo individual após a listagem em alguns ambientes. Para forçar com 100% de certeza
		// o erro de os.ReadFile sem depender de concorrência impura, podemos usar um truque de nome de arquivo
		// que passe na validação do ReadDir, mas falhe na leitura real do sistema de arquivos, ou simplesmente
		// simular um link simbólico quebrado.
		tmpDir := t.TempDir()

		// Criando um link simbólico quebrado que termina com .sql
		// os.ReadDir lê a entrada perfeitamente, mas os.ReadFile tenta seguir o link inexistente e estoura o erro!
		targetInexistente := filepath.Join(tmpDir, "nao_existo.sql")
		symlinkPath := filepath.Join(tmpDir, "01_broken_link.sql")
		if err := os.Symlink(targetInexistente, symlinkPath); err != nil {
			// Fallback alternativo caso o sistema operacional (ex: Windows sem modo desenvolvedor) proiba symlinks
			// Removemos a permissao total do arquivo
			if err := os.WriteFile(symlinkPath, []byte("SELECT 1;"), 0000); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.Chmod(symlinkPath, 0644) }()
		}

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: tmpDir,
		}

		err := cfg.RunMigrations()
		if err == nil || !strings.Contains(err.Error(), "failed to read migration file") {
			t.Errorf("esperava erro de leitura do arquivo individual, recebeu: %v", err)
		}
	})

	t.Run("Pular arquivo de migracao se ele estiver vazio ou apenas com espacos", func(t *testing.T) {
		db := newMockDB(t)
		defer db.Close()

		tmpDir := t.TempDir()

		// 01 está totalmente vazio, deve ativar o branch 'continue'
		if err := os.WriteFile(filepath.Join(tmpDir, "01_empty.sql"), []byte("   \n   "), 0644); err != nil {
			t.Fatal(err)
		}
		// 02 possui comando legítimo para garantir que o fluxo continuou com sucesso após o continue
		if err := os.WriteFile(filepath.Join(tmpDir, "02_valid.sql"), []byte("CREATE TABLE logs (id INT);"), 0644); err != nil {
			t.Fatal(err)
		}

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: tmpDir,
		}

		err := cfg.RunMigrations()
		if err != nil {
			t.Fatalf("nao esperava erro ao processar migracao vazia com continue, recebeu: %v", err)
		}

		// Garante que o 02_valid rodou após o 01_empty ter sido ignorado pelo continue
		_, err = db.Exec("INSERT INTO logs (id) VALUES (1);")
		if err != nil {
			t.Errorf("tabela logs deveria ter sido criada, indicando que o continue funcionou: %v", err)
		}
	})

	t.Run("Erro de sintaxe SQL na execucao da migracao", func(t *testing.T) {
		db := newMockDB(t)
		defer db.Close()

		tmpDir := t.TempDir()

		// Uma instrução truncada e sem sentido faz o parser do SQLite falhar de imediato
		sqlInvalido := "INSERT INTO;"

		if err := os.WriteFile(filepath.Join(tmpDir, "01_broken.sql"), []byte(sqlInvalido), 0644); err != nil {
			t.Fatal(err)
		}

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: tmpDir,
		}

		err := cfg.RunMigrations()
		if err == nil || !strings.Contains(err.Error(), "failed to execute migration file") {
			t.Errorf("esperava erro de execucao de script SQL invalido, recebeu: %v", err)
		}
	})

	t.Run("Erro se o diretorio de migracoes nao existir em disco", func(t *testing.T) {
		db := newMockDB(t)
		defer db.Close()

		// Criamos um caminho teoricamente valido dentro do TempDir, mas sem criar a pasta fisicamente
		tmpDir := t.TempDir()
		caminhoInexistente := filepath.Join(tmpDir, "pasta_totalmente_fantasma_123")

		cfg := xdb.DBConfig{
			DB:                db,
			MigrationsDirPath: caminhoInexistente, // Dispara o os.IsNotExist(err)
		}

		err := cfg.RunMigrations()
		if err == nil || !strings.Contains(err.Error(), "migrations directory does not exist") {
			t.Errorf("esperava erro de diretorio inexistente, recebeu: %v", err)
		}
	})
}
