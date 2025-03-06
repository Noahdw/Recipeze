package repo

import "github.com/jackc/pgx/v5/pgtype"

func StringPG(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}
