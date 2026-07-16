package database

import (
	"database/sql"
	"fmt"

	"github.com/yasseraitnasser/omni-association/src/utils"
)

var DB *sql.DB

func InitDB() error {
	dsn := fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=disable",
		utils.DB_NAME,
		utils.DB_HOST,
		utils.DB_PORT,
		utils.DB_USER,
		utils.DB_PASS,
	)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	return DB.Ping()
}
