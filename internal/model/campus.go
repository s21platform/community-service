package model

type Campus struct {
	Uuid      string `db:"campus_uuid"`
	ShortName string `db:"short_name"`
	FullName  string `db:"full_name"`
}
