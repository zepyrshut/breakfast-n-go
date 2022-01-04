package dbrepo

import (
	"database/sql"

	"github.com/zepyrshut/breakfast-n-go/internal/config"
	"github.com/zepyrshut/breakfast-n-go/internal/repository"
)

type postgreDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgreDBRepo{
		App: a,
		DB:  conn,
	}
}
