package utils

import "database/sql"

func PtrFromNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func NullStringFromPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func PtrFromNullInt64(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func NullInt64FromPtr(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{
		Int64: *i,
		Valid: true,
	}
}
