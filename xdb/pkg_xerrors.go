package xdb

import (
	"github.com/AeonDigital/Go-Core/xerrors"
)

const (
	XERR_NONE   xerrors.ErrorCode = ""
	XERR_PKGCTX xerrors.ErrorCode = "ERR_XDB"

	XERR_REPO_INSERT_HAS_NUMERICAL_PK     xerrors.ErrorCode = "E0001"
	XERR_REPO_INSERT_HAS_STRING_PK        xerrors.ErrorCode = "E0002"
	XERR_REPO_INSERT_NATURAL_PK_EMPTY     xerrors.ErrorCode = "E0003"
	XERR_REPO_INSERT_NATURAL_PK_NIL       xerrors.ErrorCode = "E0004"
	XERR_REPO_INSERT_EXEC_FAILED          xerrors.ErrorCode = "E0005"
	XERR_REPO_INSERT_FETCH_ID_FAILED      xerrors.ErrorCode = "E0006"
	XERR_REPO_UPDATE_INVALID_NUMERICAL_PK xerrors.ErrorCode = "E0007"
	XERR_REPO_UPDATE_EMPTY_STRING_PK      xerrors.ErrorCode = "E0008"
	XERR_REPO_UPDATE_PK_NIL               xerrors.ErrorCode = "E0009"
	XERR_REPO_UPDATE_UNKNOWN_PK_TYPE      xerrors.ErrorCode = "E0010"
	XERR_REPO_UPDATE_NO_COLUMNS_DEFINED   xerrors.ErrorCode = "E0011"
	XERR_REPO_UPDATE_EXEC_FAILED          xerrors.ErrorCode = "E0012"
	XERR_REPO_UPDATE_VERIFY_ROWS_FAILED   xerrors.ErrorCode = "E0013"
	XERR_REPO_UPDATE_RECORD_NOT_FOUND     xerrors.ErrorCode = "E0014"
	XERR_REPO_DELETE_EXEC_FAILED          xerrors.ErrorCode = "E0015"
	XERR_REPO_DELETE_VERIFY_ROWS_FAILED   xerrors.ErrorCode = "E0016"
	XERR_REPO_DELETE_RECORD_NOT_FOUND     xerrors.ErrorCode = "E0017"
	XERR_REPO_INTERFACE_ASSERTION_FAILED  xerrors.ErrorCode = "E0018"
	XERR_REPO_GET_BY_ID_EXEC_FAILED       xerrors.ErrorCode = "E0019"
	XERR_REPO_GET_BY_ID_RECORD_NOT_FOUND  xerrors.ErrorCode = "E0020"
	XERR_REPO_GET_BY_ID_SCAN_FAILED       xerrors.ErrorCode = "E0021"
	XERR_REPO_GET_ALL_EXEC_FAILED         xerrors.ErrorCode = "E0022"
	XERR_REPO_GET_ALL_SCAN_FAILED         xerrors.ErrorCode = "E0023"
	XERR_REPO_GET_ALL_ITERATION_FAILED    xerrors.ErrorCode = "E0024"
	XERR_REPO_GET_BY_FIELD_EXEC_FAILED    xerrors.ErrorCode = "E0025"
	XERR_REPO_GET_BY_FIELD_SCAN_FAILED    xerrors.ErrorCode = "E0026"
	XERR_REPO_GET_WHERE_ARGS_MISMATCH     xerrors.ErrorCode = "E0027"
	XERR_REPO_GET_WHERE_EXEC_FAILED       xerrors.ErrorCode = "E0028"
	XERR_REPO_GET_WHERE_SCAN_FAILED       xerrors.ErrorCode = "E0029"
	XERR_REPO_QUERY_RAW_EXEC_FAILED       xerrors.ErrorCode = "E0030"
	XERR_REPO_QUERY_RAW_SCAN_FAILED       xerrors.ErrorCode = "E0031"
)

// ErrorMessages acts as a centralized English translation catalog for core repository errors.
var ErrorMessages = map[xerrors.ErrorCode]string{
	XERR_REPO_INSERT_HAS_NUMERICAL_PK:     "cannot insert entity with an existing numerical ID",
	XERR_REPO_INSERT_HAS_STRING_PK:        "cannot insert entity with an existing primary key string",
	XERR_REPO_INSERT_NATURAL_PK_EMPTY:     "cannot insert entity: natural primary key string cannot be empty",
	XERR_REPO_INSERT_NATURAL_PK_NIL:       "cannot insert entity: natural primary key cannot be nil",
	XERR_REPO_INSERT_EXEC_FAILED:          "failed to execute insert statement in the database",
	XERR_REPO_INSERT_FETCH_ID_FAILED:      "failed to retrieve last inserted numerical ID from database",
	XERR_REPO_UPDATE_INVALID_NUMERICAL_PK: "cannot update entity with invalid numerical ID (must be > 0)",
	XERR_REPO_UPDATE_EMPTY_STRING_PK:      "cannot update entity with empty string key",
	XERR_REPO_UPDATE_PK_NIL:               "cannot update entity without a primary key",
	XERR_REPO_UPDATE_UNKNOWN_PK_TYPE:      "unknown primary key type format",
	XERR_REPO_UPDATE_NO_COLUMNS_DEFINED:   "no columns defined for update operation in this table",
	XERR_REPO_UPDATE_EXEC_FAILED:          "failed to execute update statement in the database",
	XERR_REPO_UPDATE_VERIFY_ROWS_FAILED:   "failed to verify affected rows during update operation",
	XERR_REPO_UPDATE_RECORD_NOT_FOUND:     "target record not found in database for update operation",
	XERR_REPO_DELETE_EXEC_FAILED:          "failed to execute delete statement in the database",
	XERR_REPO_DELETE_VERIFY_ROWS_FAILED:   "failed to verify affected rows during delete operation",
	XERR_REPO_DELETE_RECORD_NOT_FOUND:     "target record not found in database for delete operation",
	XERR_REPO_INTERFACE_ASSERTION_FAILED:  "internal error: entity type does not implement database dao interface",
	XERR_REPO_GET_BY_ID_EXEC_FAILED:       "failed to execute select single record query",
	XERR_REPO_GET_BY_ID_RECORD_NOT_FOUND:  "requested record not found in the database table",
	XERR_REPO_GET_BY_ID_SCAN_FAILED:       "failed to scan database columns into entity memory pointers",
	XERR_REPO_GET_ALL_EXEC_FAILED:         "failed to execute select collection list query",
	XERR_REPO_GET_ALL_SCAN_FAILED:         "failed to scan row iteration into entity memory pointers",
	XERR_REPO_GET_ALL_ITERATION_FAILED:    "database cursor failure during rows iteration process",
	XERR_REPO_GET_BY_FIELD_EXEC_FAILED:    "failed to execute field query statement in the database",
	XERR_REPO_GET_BY_FIELD_SCAN_FAILED:    "failed to scan database row into the entity fields",
	XERR_REPO_GET_WHERE_ARGS_MISMATCH:     "query parameters count does not match the provided arguments length",
	XERR_REPO_GET_WHERE_EXEC_FAILED:       "failed to execute conditional query statement in the database",
	XERR_REPO_GET_WHERE_SCAN_FAILED:       "failed to scan conditional database row into the entity fields",
	XERR_REPO_QUERY_RAW_EXEC_FAILED:       "failed to execute raw sql query statement in the database",
	XERR_REPO_QUERY_RAW_SCAN_FAILED:       "failed to scan database row into the raw query destination structure",
}
