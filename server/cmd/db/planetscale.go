package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"teniditter-server/cmd/global/console"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type dbManager struct {
}

var DBManager = &dbManager{}

var sqlConn *sql.DB

func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, nil
}

var counts int64

func connectToDB() *sql.DB {
	for {
		connection, err := openDB()
		if err != nil {
			console.Log(fmt.Sprintf("Planetscale not yet ready. Attempt n°%02d\n", counts), console.Warning)
			counts++
		} else {
			console.Log("Connected to Planetscale!", console.Success, true)
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		console.Log("Backing off for two seconds", console.Info)
		time.Sleep(2 * time.Second)
		continue
	}
}

// Create a new
func (*dbManager) NewDB() *sql.DB {
	if sqlConn != nil {
		return sqlConn
	}

	db := connectToDB()
	if db != nil {
		sqlConn = db
	}

	return db
}

// Connect to an already existing DB
func (*dbManager) Connect() *sql.DB {
	return sqlConn
}

func (*dbManager) Disconnect() bool {
	if sqlConn == nil {
		return false
	}

	err := sqlConn.Close()
	if err == nil {
		sqlConn = nil
		console.Log("Planetscale Disconnected", console.Warning)
	} else {
		console.Log("Failed to disconnect planetscale", console.Error)
	}

	return err == nil
}

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DATABASE_NAME"))
}
