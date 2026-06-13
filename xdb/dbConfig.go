package xdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// DBConfig manages database driver options, connection pool life cycles, and engine-specific options.
type DBConfig struct {
	Driver                string        `json:"driver"`
	DSN                   string        `json:"dsn"`
	MigrationsDirPath     string        `json:"migrationsDirPath"`
	MaxOpenConnections    int           `json:"maxOpenConnections"`
	MaxIdleConnections    int           `json:"maxIdleConnections"`
	ConnectionMaxLifetime time.Duration `json:"connectionMaxLifetime"`
	SQLite                SQLiteConfig  `json:"sqlite"`

	DB *sql.DB
}

// NewDBConfig instantiates a baseline operational configuration setup with optimized defaults for embedded deployments.
func NewDBConfig(
	driver string,
	dsn string,
	migrationsDirPath string,
	sqliteMode string,
	sqliteDir string,
	sqliteFileName string,
	sqliteQuerystring string,
) DBConfig {
	sqlite := SQLiteConfig{
		Mode:        sqliteMode,
		Dir:         sqliteDir,
		FileName:    sqliteFileName,
		QueryString: sqliteQuerystring,

		Pragma: map[string]string{
			"journal_mode": "WAL",
			"synchronous":  "NORMAL",
			"busy_timeout": "5000",
			"foreign_keys": "ON",
		},
	}

	dbCfg := DBConfig{
		Driver:                driver,
		DSN:                   dsn,
		MigrationsDirPath:     migrationsDirPath,
		MaxOpenConnections:    1,
		MaxIdleConnections:    1,
		ConnectionMaxLifetime: 30 * time.Minute,
		SQLite:                sqlite,
	}

	dbCfg.buildDSN()

	return dbCfg
}

// buildDSN synthesizes target variables and environments to compile a clean connection URI sequence.
func (o *DBConfig) buildDSN() {
	if o.DSN != "" {
		o.SQLite.Mode = ""
		o.SQLite.Dir = ""
		o.SQLite.FileName = ""
		o.SQLite.QueryString = ""
		return
	}

	o.Driver = "sqlite"

	mode := strings.TrimSpace(o.SQLite.Mode)
	if mode == "" {
		mode = "file:"
	}

	if mode == ":memory:" || mode == "file::memory:" {
		if o.SQLite.QueryString != "" {
			if mode == "file::memory:" {
				o.DSN = mode + "?" + o.SQLite.QueryString
				return
			}
			o.DSN = "file::memory:?" + o.SQLite.QueryString
			return
		}
		o.DSN = mode
		return
	}

	if o.SQLite.Dir == "" || o.SQLite.Dir == "." || o.SQLite.FileName == "" {
		return
	}

	fullPath := filepath.Join(o.SQLite.Dir, o.SQLite.FileName)

	dsn := fmt.Sprintf("file:%s", filepath.ToSlash(fullPath))

	if o.SQLite.QueryString != "" {
		dsn = fmt.Sprintf("%s?%s", dsn, o.SQLite.QueryString)
	}

	o.DSN = dsn
}

// CheckConfiguration verifies structural parameters to maintain setup layout consistency before initializing adapters.
func (o *DBConfig) CheckConfiguration() error {
	o.buildDSN()
	return nil
}

// InitDataBaseConnection activates database interface structures and applies custom engine tuning options safely.
func (o *DBConfig) InitDataBaseConnection() error {
	db, err := sql.Open(o.Driver, o.DSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(o.MaxOpenConnections)
	db.SetMaxIdleConns(o.MaxIdleConnections)
	db.SetConnMaxLifetime(o.ConnectionMaxLifetime)

	if strings.Contains(o.Driver, "sqlite") {
		requiredDefaults := map[string]string{
			"journal_mode": "WAL",
			"synchronous":  "NORMAL",
			"busy_timeout": "5000",
			"foreign_keys": "ON",
		}

		if o.SQLite.Pragma == nil {
			o.SQLite.Pragma = make(map[string]string)
		}

		for key, defaultValue := range requiredDefaults {
			if _, exists := o.SQLite.Pragma[key]; !exists {
				o.SQLite.Pragma[key] = defaultValue
			}
		}

		for key, value := range o.SQLite.Pragma {
			query := fmt.Sprintf("PRAGMA %s = %s;", key, value)
			if _, err := db.Exec(query); err != nil {
				db.Close()
				return err
			}
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(pingCtx); err != nil {
			db.Close()
			return err
		}
	}

	o.DB = db
	return err
}

// RunMigrations processes organized script collections chronologically against target database environments.
func (o *DBConfig) RunMigrations() error {
	if o.DB == nil {
		return fmt.Errorf("database instance is nil")
	}

	dirPath := strings.TrimSpace(o.MigrationsDirPath)
	if dirPath == "" {
		return fmt.Errorf("migrations directory path is empty")
	}

	dirInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", dirPath)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("provided path is a file, not a directory: %s", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	sort.Strings(migrationFiles)

	if len(migrationFiles) == 0 {
		return nil
	}

	for _, fileName := range migrationFiles {
		fullFilePath := filepath.Join(dirPath, fileName)

		scriptBytes, err := os.ReadFile(fullFilePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file [%s]: %w", fileName, err)
		}

		scriptSQL := string(scriptBytes)
		if strings.TrimSpace(scriptSQL) == "" {
			continue
		}

		_, err = o.DB.ExecContext(context.Background(), scriptSQL)
		if err != nil {
			return fmt.Errorf("failed to execute migration file [%s]: %w", fileName, err)
		}
	}

	return nil
}
