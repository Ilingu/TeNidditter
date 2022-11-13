package db

import (
	"encoding/json"
	"errors"
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/services/nitter"
)

func IsNeetAlreadyExist(neetId string) bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Neets WHERE neet_id=?;", neetId).Scan(&count)
	if err != nil {
		return false
	}

	return count != 0
}

func InsertNewNeet(neet nitter.NeetComment) bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	jsonBlob, err := json.Marshal(neet)
	if err != nil {
		return false
	}

	_, err = db.Exec("INSERT INTO Neets (neet_id, neet_data) VALUES (?, ?);", neet.Id, jsonBlob)
	return err == nil
}

func GetNeetById(neetId string) (*DBNeets, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	var neet DBNeets
	err := db.QueryRow("SELECT * FROM Neets WHERE neet_id=?;", neetId).Scan(&neet.NeetId, &neet.Neet)
	if err != nil {
		return nil, errors.New("cannot fetch neet")
	}

	return &neet, nil
}
