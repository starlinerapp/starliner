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
