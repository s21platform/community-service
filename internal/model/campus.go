package model

type Campus struct {
	Id        int64  `json:"id"`
	Uuid      string `db:"campus_uuid"`
	ShortName string `db:"short_name"`
	FullName  string `db:"full_name"`
}
