package xdb

import (
	"github.com/AeonDigital/Go-Core/xerrors"
)

const (
	ErrNone xerrors.ErrorCode = ""
	ErrXDB  xerrors.ErrorCode = "ERR-XDB"

	ErrRepoInsertHasNumericalPK     xerrors.ErrorCode = "E0001"
	ErrRepoInsertHasStringPK        xerrors.ErrorCode = "E0002"
	ErrRepoInsertNaturalPKEmpty     xerrors.ErrorCode = "E0003"
	ErrRepoInsertNaturalPKNil       xerrors.ErrorCode = "E0004"
	ErrRepoInsertExecFailed         xerrors.ErrorCode = "E0005"
	ErrRepoInsertFetchIDFailed      xerrors.ErrorCode = "E0006"
	ErrRepoUpdateInvalidNumericalPK xerrors.ErrorCode = "E0007"
	ErrRepoUpdateEmptyStringPK      xerrors.ErrorCode = "E0008"
	ErrRepoUpdatePKNil              xerrors.ErrorCode = "E0009"
	ErrRepoUpdateUnknownPKType      xerrors.ErrorCode = "E0010"
	ErrRepoUpdateNoColumnsDefined   xerrors.ErrorCode = "E0011"
	ErrRepoUpdateExecFailed         xerrors.ErrorCode = "E0012"
	ErrRepoUpdateVerifyRowsFailed   xerrors.ErrorCode = "E0013"
	ErrRepoUpdateRecordNotFound     xerrors.ErrorCode = "E0014"
	ErrRepoDeleteExecFailed         xerrors.ErrorCode = "E0015"
	ErrRepoDeleteVerifyRowsFailed   xerrors.ErrorCode = "E0016"
	ErrRepoDeleteRecordNotFound     xerrors.ErrorCode = "E0017"
	ErrRepoInterfaceAssertionFailed xerrors.ErrorCode = "E0018"
	ErrRepoGetByIDExecFailed        xerrors.ErrorCode = "E0019"
	ErrRepoGetByIDRecordNotFound    xerrors.ErrorCode = "E0020"
	ErrRepoGetByIDScanFailed        xerrors.ErrorCode = "E0021"
	ErrRepoGetAllExecFailed         xerrors.ErrorCode = "E0022"
	ErrRepoGetAllScanFailed         xerrors.ErrorCode = "E0023"
	ErrRepoGetAllIterationFailed    xerrors.ErrorCode = "E0024"
	ErrRepoGetByFieldExecFailed     xerrors.ErrorCode = "E0025"
	ErrRepoGetByFieldScanFailed     xerrors.ErrorCode = "E0026"
	ErrRepoGetWhereArgsMismatch     xerrors.ErrorCode = "E0027"
	ErrRepoGetWhereExecFailed       xerrors.ErrorCode = "E0028"
	ErrRepoGetWhereScanFailed       xerrors.ErrorCode = "E0029"
	ErrRepoQueryRawExecFailed       xerrors.ErrorCode = "E0030"
	ErrRepoQueryRawScanFailed       xerrors.ErrorCode = "E0031"
)

// ErrorMessages acts as a centralized English translation catalog for core repository errors.
var ErrorMessages = map[xerrors.ErrorCode]string{
	ErrRepoInsertHasNumericalPK:     "cannot insert entity with an existing numerical ID",
	ErrRepoInsertHasStringPK:        "cannot insert entity with an existing primary key string",
	ErrRepoInsertNaturalPKEmpty:     "cannot insert entity: natural primary key string cannot be empty",
	ErrRepoInsertNaturalPKNil:       "cannot insert entity: natural primary key cannot be nil",
	ErrRepoInsertExecFailed:         "failed to execute insert statement in the database",
	ErrRepoInsertFetchIDFailed:      "failed to retrieve last inserted numerical ID from database",
	ErrRepoUpdateInvalidNumericalPK: "cannot update entity with invalid numerical ID (must be > 0)",
	ErrRepoUpdateEmptyStringPK:      "cannot update entity with empty string key",
	ErrRepoUpdatePKNil:              "cannot update entity without a primary key",
	ErrRepoUpdateUnknownPKType:      "unknown primary key type format",
	ErrRepoUpdateNoColumnsDefined:   "no columns defined for update operation in this table",
	ErrRepoUpdateExecFailed:         "failed to execute update statement in the database",
	ErrRepoUpdateVerifyRowsFailed:   "failed to verify affected rows during update operation",
	ErrRepoUpdateRecordNotFound:     "target record not found in database for update operation",
	ErrRepoDeleteExecFailed:         "failed to execute delete statement in the database",
	ErrRepoDeleteVerifyRowsFailed:   "failed to verify affected rows during delete operation",
	ErrRepoDeleteRecordNotFound:     "target record not found in database for delete operation",
	ErrRepoInterfaceAssertionFailed: "internal error: entity type does not implement database dao interface",
	ErrRepoGetByIDExecFailed:        "failed to execute select single record query",
	ErrRepoGetByIDRecordNotFound:    "requested record not found in the database table",
	ErrRepoGetByIDScanFailed:        "failed to scan database columns into entity memory pointers",
	ErrRepoGetAllExecFailed:         "failed to execute select collection list query",
	ErrRepoGetAllScanFailed:         "failed to scan row iteration into entity memory pointers",
	ErrRepoGetAllIterationFailed:    "database cursor failure during rows iteration process",
	ErrRepoGetByFieldExecFailed:     "failed to execute field query statement in the database",
	ErrRepoGetByFieldScanFailed:     "failed to scan database row into the entity fields",
	ErrRepoGetWhereArgsMismatch:     "query parameters count does not match the provided arguments length",
	ErrRepoGetWhereExecFailed:       "failed to execute conditional query statement in the database",
	ErrRepoGetWhereScanFailed:       "failed to scan conditional database row into the entity fields",
	ErrRepoQueryRawExecFailed:       "failed to execute raw sql query statement in the database",
	ErrRepoQueryRawScanFailed:       "failed to scan database row into the raw query destination structure",
}
