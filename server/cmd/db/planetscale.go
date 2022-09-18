package db

import (
	"database/sql"
	"fmt"
	"os"
)

func generateDbConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DATABASE_NAME"))
}

// Just Testing, not official
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", generateDbConnString())
	if err != nil {
		return nil, err
	}
	return db, nil
}
