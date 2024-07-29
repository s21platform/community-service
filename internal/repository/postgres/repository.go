package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/s21platform/community-service/internal/config"
	"log"
)

type Repository struct {
	conn *sql.DB
}

func New(cfg *config.Config) *Repository {

	connectCmd := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sql.Open("postgres", connectCmd)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &Repository{db}
}

func (r *Repository) Close() {
	r.conn.Close()
}
