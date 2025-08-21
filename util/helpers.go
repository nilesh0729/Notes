package util

import "database/sql"

func ToNullInt32(v int32) sql.NullInt32 {
	return sql.NullInt32{Int32: v, Valid: true}
}
