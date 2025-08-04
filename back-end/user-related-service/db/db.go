package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"user-related-service/config"
)

var DB *sqlx.DB

func Init() error {
	var err error
	DB, err = sqlx.Connect("postgres", config.DB_URL)
	return err
}
