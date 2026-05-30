package xdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// sqlExecutor unifies common database operations available on both *sql.DB and *sql.Tx connections.
type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// DBGeneric implements a type-safe, generic repository pattern dedicated to a specific domain entity model.
type DBGeneric[T any, PT interface {
	*T
	Entity
}] struct {
	db                     *sql.DB
	tx                     *sql.Tx
	executor               sqlExecutor
	idempotentUpdateActive bool
	idempotentDeleteActive bool
}

// NewDBGeneric initializes and yields a new operational instance of the generic repository interface.
func NewDBGeneric[T any, PT interface {
	*T
	Entity
}](db *sql.DB) *DBGeneric[T, PT] {
	return &DBGeneric[T, PT]{
		db:       db,
		executor: db,
	}
}

// WithTx returns a contextual shallow clone of the repository bound to an active database transaction lifecycle.
func (r *DBGeneric[T, PT]) WithTx(tx *sql.Tx) *DBGeneric[T, PT] {
	return &DBGeneric[T, PT]{
		db:                     r.db,
		tx:                     tx,
		executor:               tx,
		idempotentUpdateActive: r.idempotentUpdateActive,
		idempotentDeleteActive: r.idempotentDeleteActive,
	}
}

// getExecutor evaluates whether to pipeline execution states through an active transaction isolation or the global connection pool.
func (r *DBGeneric[T, PT]) getExecutor() sqlExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// ContextWithForcedIdempotency wraps the provided context to guarantee that down-stream update and delete actions run idempotently.
func ContextWithForcedIdempotency(ctx context.Context) context.Context {
	return context.WithValue(ctx, forceIdempotencyKey, true)
}

// ContextWithProhibitedIdempotency wraps the provided context to strictly block idempotent behavior on down-stream data modifications.
func ContextWithProhibitedIdempotency(ctx context.Context) context.Context {
	return context.WithValue(ctx, prohibitIdempotencyKey, true)
}

// QueryRaw coordinates the isolation, manual execution, and custom collection scan mapping of arbitrary database commands.
func QueryRaw[R any](
	ctx context.Context,
	db *sql.DB,
	cq CustomQuery[R],
) ([]R, xerrors.ErrorCode) {
	rows, err := db.QueryContext(ctx, cq.SQL, cq.Args...)
	if err != nil {
		logRepoError(ctx, db, ErrRepoQueryRawExecFailed, err, cq.SQL, cq.Args)
		return nil, ErrRepoQueryRawExecFailed
	}
	defer rows.Close()

	var result []R
	for rows.Next() {
		item, err := cq.Scanner(rows)
		if err != nil {
			logRepoError(ctx, db, ErrRepoQueryRawScanFailed, err, cq.SQL, cq.Args)
			return nil, ErrRepoQueryRawScanFailed
		}
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		logRepoError(ctx, db, ErrRepoGetAllIterationFailed, err, cq.SQL, cq.Args)
		return nil, ErrRepoGetAllIterationFailed
	}

	return result, ErrNone
}

// SetIdempotentUpdate overrides instance update settings to prevent throwing missing record validation errors on missing datasets.
func (r *DBGeneric[T, PT]) SetIdempotentUpdate(enabled bool) *DBGeneric[T, PT] {
	r.idempotentUpdateActive = enabled
	return r
}

// SetIdempotentDelete overrides instance delete settings to silently tolerate non-existent entities during removal attempts.
func (r *DBGeneric[T, PT]) SetIdempotentDelete(enabled bool) *DBGeneric[T, PT] {
	r.idempotentDeleteActive = enabled
	return r
}

// shouldBeIdempotent calculates the situational hierarchy between active context variables and historical instance flags.
func (r *DBGeneric[T, PT]) shouldBeIdempotent(ctx context.Context, isUpdate bool) bool {
	if prohibit, _ := ctx.Value(prohibitIdempotencyKey).(bool); prohibit {
		return false
	}

	if force, _ := ctx.Value(forceIdempotencyKey).(bool); force {
		return true
	}

	if isUpdate {
		return r.idempotentUpdateActive
	}
	return r.idempotentDeleteActive
}

// Insert validates model states, produces identifiers when applicable, and stores records securely.
func (r *DBGeneric[T, PT]) Insert(ctx context.Context, entity PT) xerrors.ErrorCode {
	currentPK := entity.PKValue()

	switch v := currentPK.(type) {
	case int64:
		if v > 0 {
			logRepoError(ctx, r.db, ErrRepoInsertHasNumericalPK, nil, "", nil)
			return ErrRepoInsertHasNumericalPK
		}
	case string:
		if !entity.IsNaturalPK() && strings.TrimSpace(v) != "" {
			logRepoError(ctx, r.db, ErrRepoInsertHasStringPK, nil, "", nil)
			return ErrRepoInsertHasStringPK
		}
		if entity.IsNaturalPK() && strings.TrimSpace(v) == "" {
			logRepoError(ctx, r.db, ErrRepoInsertNaturalPKEmpty, nil, "", nil)
			return ErrRepoInsertNaturalPKEmpty
		}
	case nil:
		if entity.IsNaturalPK() {
			logRepoError(ctx, r.db, ErrRepoInsertNaturalPKNil, nil, "", nil)
			return ErrRepoInsertNaturalPKNil
		}
	}

	generated := entity.GeneratePK()
	if generated != nil {
		entity.BindPK(generated)
	}

	entity.Normalize()
	if valid, errCode := entity.Validate(); !valid {
		logRepoError(ctx, r.db, errCode, nil, "", nil)
		return errCode
	}

	cols := entity.Columns()
	values := entity.Values()

	placeholders := make([]string, len(cols))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s);",
		entity.TableName(),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	executor := r.getExecutor()
	result, err := executor.ExecContext(ctx, query, values...)
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoInsertExecFailed, err, query, values)
		return ErrRepoInsertExecFailed
	}

	if entity.TablePK() == "id" {
		currentVal := entity.PKValue()

		isUnset := currentVal == nil
		if !isUnset {
			pkInt, ok := currentVal.(int64)
			if pkInt == 0 && ok {
				isUnset = true
			}
		}

		if isUnset {
			id, err := result.LastInsertId()
			if id > 0 && err == nil {
				entity.BindPK(id)
			}
		}
	}

	return ErrNone
}

// Update coordinates column updates while enforcing primary key validations and idempotency rules.
func (r *DBGeneric[T, PT]) Update(ctx context.Context, entity PT) xerrors.ErrorCode {
	currentPK := entity.PKValue()

	switch v := currentPK.(type) {
	case int64:
		if v <= 0 {
			logRepoError(ctx, r.db, ErrRepoUpdateInvalidNumericalPK, nil, "", nil)
			return ErrRepoUpdateInvalidNumericalPK
		}
	case string:
		if strings.TrimSpace(v) == "" {
			logRepoError(ctx, r.db, ErrRepoUpdateEmptyStringPK, nil, "", nil)
			return ErrRepoUpdateEmptyStringPK
		}
	case nil:
		logRepoError(ctx, r.db, ErrRepoUpdatePKNil, nil, "", nil)
		return ErrRepoUpdatePKNil
	default:
		logRepoError(ctx, r.db, ErrRepoUpdateUnknownPKType, nil, "", nil)
		return ErrRepoUpdateUnknownPKType
	}

	entity.Normalize()

	if valid, errCode := entity.Validate(); !valid {
		logRepoError(ctx, r.db, errCode, nil, "", nil)
		return errCode
	}

	cols := entity.Columns()
	values := entity.Values()

	if len(cols) == 0 {
		logRepoError(ctx, r.db, ErrRepoUpdateNoColumnsDefined, nil, "", nil)
		return ErrRepoUpdateNoColumnsDefined
	}

	setFragments := make([]string, len(cols))
	for i, col := range cols {
		setFragments[i] = fmt.Sprintf("%s = ?", col)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = ?;",
		entity.TableName(),
		strings.Join(setFragments, ", "),
		entity.TablePK(),
	)

	args := append(values, entity.PKValue())

	result, err := r.executor.ExecContext(ctx, query, args...)
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoUpdateExecFailed, err, query, args)
		return ErrRepoUpdateExecFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoUpdateVerifyRowsFailed, err, query, args)
		return ErrRepoUpdateVerifyRowsFailed
	}

	if rowsAffected == 0 {
		if r.shouldBeIdempotent(ctx, true) {
			return ErrNone
		}

		logRepoError(ctx, r.db, ErrRepoUpdateRecordNotFound, nil, query, args)
		return ErrRepoUpdateRecordNotFound
	}

	return ErrNone
}

// Delete drops target records based on explicit key mapping evaluations.
func (r *DBGeneric[T, PT]) Delete(ctx context.Context, entity PT) xerrors.ErrorCode {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = ?;",
		entity.TableName(),
		entity.TablePK(),
	)

	pkValue := entity.PKValue()
	args := []any{pkValue}

	result, err := r.executor.ExecContext(ctx, query, pkValue)
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoDeleteExecFailed, err, query, args)
		return ErrRepoDeleteExecFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoDeleteVerifyRowsFailed, err, query, args)
		return ErrRepoDeleteVerifyRowsFailed
	}

	if rowsAffected == 0 {
		if r.shouldBeIdempotent(ctx, false) {
			return ErrNone
		}

		logRepoError(ctx, r.db, ErrRepoDeleteRecordNotFound, nil, query, args)
		return ErrRepoDeleteRecordNotFound
	}

	return ErrNone
}

// GetByID performs a target row execution based on primary key mappings to fetch a singular type-safe entry instance.
func (r *DBGeneric[T, PT]) GetByID(ctx context.Context, id any) (*T, xerrors.ErrorCode) {
	var instance PT = new(T)

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s = ? LIMIT 1;",
		instance.TableName(),
		instance.TablePK(),
	)

	args := []any{id}

	rows, err := r.executor.QueryContext(ctx, query, id)
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoGetByIDExecFailed, err, query, args)
		return nil, ErrRepoGetByIDExecFailed
	}
	defer rows.Close()

	if !rows.Next() {
		logRepoError(ctx, r.db, ErrRepoGetByIDRecordNotFound, nil, query, args)
		return nil, ErrRepoGetByIDRecordNotFound
	}

	if err := instance.ScanRow(rows); err != nil {
		logRepoError(ctx, r.db, ErrRepoGetByIDScanFailed, err, query, args)
		return nil, ErrRepoGetByIDScanFailed
	}

	return instance, ErrNone
}

// GetAll extracts every existing collection sequence context from the entity schema targets.
func (r *DBGeneric[T, PT]) GetAll(ctx context.Context) ([]*T, xerrors.ErrorCode) {
	var meta PT = new(T)
	query := fmt.Sprintf(
		"SELECT * FROM %s;",
		meta.TableName(),
	)

	rows, err := r.executor.QueryContext(ctx, query)
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoGetAllExecFailed, err, query, nil)
		return nil, ErrRepoGetAllExecFailed
	}
	defer rows.Close()

	var list []*T
	for rows.Next() {
		var item PT = new(T)

		if err := item.ScanRow(rows); err != nil {
			logRepoError(ctx, r.db, ErrRepoGetAllScanFailed, err, query, nil)
			return nil, ErrRepoGetAllScanFailed
		}

		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		logRepoError(ctx, r.db, ErrRepoGetAllIterationFailed, err, query, nil)
		return nil, ErrRepoGetAllIterationFailed
	}

	return list, ErrNone
}

// GetByField searches for matching dataset groups filtered by a specific column variable signature.
func (r *DBGeneric[T, PT]) GetByField(ctx context.Context, field string, value any) ([]*T, xerrors.ErrorCode) {
	var meta PT = new(T)

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?;", meta.TableName(), field)
	args := []any{value}

	rows, err := r.executor.QueryContext(ctx, query, value)
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoGetByFieldExecFailed, err, query, args)
		return nil, ErrRepoGetByFieldExecFailed
	}
	defer rows.Close()

	var list []*T
	for rows.Next() {
		var item PT = new(T)
		if err := item.ScanRow(rows); err != nil {
			logRepoError(ctx, r.db, ErrRepoGetByFieldScanFailed, err, query, args)
			return nil, ErrRepoGetByFieldScanFailed
		}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		logRepoError(ctx, r.db, ErrRepoGetAllIterationFailed, err, query, args)
		return nil, ErrRepoGetAllIterationFailed
	}

	return list, ErrNone
}

// GetWhere parses complex conditional dynamic parameters to retrieve subset target collections.
func (r *DBGeneric[T, PT]) GetWhere(ctx context.Context, queryFragment string, args ...any) ([]*T, xerrors.ErrorCode) {
	placeholderCount := strings.Count(queryFragment, "?")
	argCount := len(args)

	if placeholderCount != argCount {
		logRepoError(ctx, r.db, ErrRepoGetWhereArgsMismatch, nil, queryFragment, args)
		return nil, ErrRepoGetWhereArgsMismatch
	}

	var meta PT = new(T)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s;", meta.TableName(), queryFragment)

	rows, err := r.executor.QueryContext(ctx, query, args...)
	if err != nil {
		logRepoError(ctx, r.db, ErrRepoGetWhereExecFailed, err, query, args)
		return nil, ErrRepoGetWhereExecFailed
	}
	defer rows.Close()

	var list []*T
	for rows.Next() {
		var item PT = new(T)
		if err := item.ScanRow(rows); err != nil {
			logRepoError(ctx, r.db, ErrRepoGetWhereScanFailed, err, query, args)
			return nil, ErrRepoGetWhereScanFailed
		}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		logRepoError(ctx, r.db, ErrRepoGetAllIterationFailed, err, query, args)
		return nil, ErrRepoGetAllIterationFailed
	}

	return list, ErrNone
}
