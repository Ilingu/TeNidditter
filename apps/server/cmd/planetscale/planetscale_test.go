package ps_test

import (
	"os"
	utils_env "teniditter-server/cmd/global/utils/env"
	ps "teniditter-server/cmd/planetscale"
	"testing"
)

func init() {
	os.Setenv("TEST", "1")
	utils_env.LoadEnv()
}

func TestNewDB(t *testing.T) {
	db := ps.DBManager.NewDB()
	if db == nil {
		t.Fatal("cannot connect to db")
	}
}

func TestConnectDB(t *testing.T) {
	db := ps.DBManager.Connect()
	if db == nil {
		t.Fatal("connection not opened")
	}
}

var dataArgs = []any{-6, "0eb4d2844e98a*6d568e0e8a507fc9f976aea10b97a91db4d46f58dfbce33d30c"}

type NitterListModel struct {
	ListID    uint   `json:"list_id"`
	AccountID int    `json:"-"`
	ListName  string `json:"title"`
}

func TestCRUD(t *testing.T) {
	db := ps.DBManager.Connect()
	if db == nil {
		t.Fatal("connection not opened")
	}

	var list NitterListModel

	// WRITE
	_, err := db.Exec("INSERT INTO NitterLists (account_id, title) VALUES (?, ?);", dataArgs...)
	if err != nil {
		t.Fatalf("cannot write, %s", err)
	}

	// READ
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM NitterLists WHERE account_id=? AND title=?;", dataArgs...).Scan(&count)
	if err != nil || count <= 0 {
		t.Fatalf("cannot read, %s", err)
	}

	row := db.QueryRow("SELECT * FROM NitterLists WHERE account_id=? AND title=?;", dataArgs...)
	err = row.Scan(&list.ListID, &list.AccountID, &list.ListName)
	if err != nil || list == (NitterListModel{}) {
		t.Fatalf("cannot read, %s", err)
	}
	if list.AccountID != dataArgs[0] || list.ListName != dataArgs[1] {
		t.Fatal("invalid read")
	}

	// UPDATE
	newTitle := "8a59c769904caa3233926*19a9434ee5cb7ea3f17ce4779c320fcc59d60299ad5"
	_, err = db.Exec("UPDATE NitterLists SET title=? WHERE list_id=?;", newTitle, list.ListID)
	if err != nil {
		t.Fatalf("cannot update, %s", err)
	}

	row = db.QueryRow("SELECT * FROM NitterLists WHERE list_id=?;", list.ListID)
	err = row.Scan(&list.ListID, &list.AccountID, &list.ListName)
	if err != nil || list == (NitterListModel{}) {
		t.Fatalf("cannot read, %s", err)
	}
	if list.ListName != newTitle {
		t.Fatal("invalid update")
	}

	// Delete
	_, err = db.Exec("DELETE FROM NitterLists WHERE list_id=?;", list.ListID)
	if err != nil {
		t.Fatalf("cannot delete, %s", err)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM NitterLists WHERE account_id=? AND title=?;", dataArgs...).Scan(&count)
	if err != nil && count == 0 {
		t.Fatalf("cannot read, %s", err)
	}
}

func TestDisconnectDB(t *testing.T) {
	if ok := ps.DBManager.Disconnect(); !ok {
		t.Fatal("cannot disconnect db")
	}
}
