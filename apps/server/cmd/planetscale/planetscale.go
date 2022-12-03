package ps

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"teniditter-server/cmd/global/console"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type databaseManager struct {
}

var DBManager = &databaseManager{}

var sqlConn *sql.DB

func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DSN")+"&parseTime=true")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

var counts int64

func connectToDB() *sql.DB {
	for {
		connection, err := openDB()
		if err != nil {
			console.Warn(fmt.Sprintf("Planetscale not yet ready. Attempt n°%02d\n", counts))
			counts++
		} else {
			console.Success("Connected to Planetscale!", true)
			return connection
		}

		if counts > 10 {
			log.Println(err)
			console.Error("Failed to connect to Planetscale", true)
			return nil
		}

		console.Log("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

// Create a new
func (*databaseManager) NewDB() *sql.DB {
	if sqlConn != nil {
		return sqlConn
	}

	sqlConn = connectToDB()
	return sqlConn
}

// Connect to an already existing DB
func (*databaseManager) Connect() *sql.DB {
	return sqlConn
}

func (*databaseManager) Disconnect() bool {
	if sqlConn == nil {
		return false
	}

	err := sqlConn.Close()
	if err == nil {
		sqlConn = nil
		console.Warn("Planetscale Disconnected")
	} else {
		console.Error("Failed to disconnect planetscale")
	}

	return err == nil
}

var ErrDbNotFound = errors.New("no db connected")
