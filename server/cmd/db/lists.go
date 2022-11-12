package db

import (
	"encoding/json"
	"errors"
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/services/nitter"
)

func CreateNitterList(user *AccountModel, listname string) bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	_, err := db.Exec("INSERT INTO NitterLists (account_id, title) VALUES (?, ?);", user.AccountId, listname)
	if err != nil {
		return false
	}

	go user.ListHasChange()
	return true
}

func GetListById(listId uint) (*NitterListModel, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	var list NitterListModel
	err := db.QueryRow("SELECT * FROM NitterLists WHERE list_id=?", listId).Scan(&list.ListID, &list.AccountID, &list.ListName)
	if err != nil {
		return nil, errors.New("cannot fetch user")
	}

	return &list, nil
}

func GetListContentById(listId uint) ([]nitter.NeetComment, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	sqlQuery := "SELECT neet_data FROM ListToNeet INNER JOIN Neets ON ListToNeet.neet_id = Neets.neet_id WHERE list_id=?"
	rows, err := db.Query(sqlQuery, listId)
	if err != nil {
		return nil, errors.New("error when neets")
	}
	defer rows.Close()

	type sqlResult struct {
		Neet string `json:"neet_data"`
	}

	var neets []nitter.NeetComment
	for rows.Next() {
		var rowResp sqlResult
		if err := rows.Scan(&rowResp.Neet); err != nil {
			return nil, err
		}

		var neet nitter.NeetComment
		if err := json.Unmarshal([]byte(rowResp.Neet), &neet); err != nil {
			return nil, err
		}

		neets = append(neets, neet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return neets, nil
}

func (list NitterListModel) AddNeet(neetId string) bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM ListToNeet WHERE list_id=? AND neet_id=?;", list.ListID, neetId).Scan(&count)
	if err == nil && count != 0 {
		return false
	}

	_, err = db.Exec("INSERT INTO ListToNeet (list_id, neet_id) VALUES (?, ?);", list.ListID, neetId)
	return err == nil
}
func (list NitterListModel) RemoveNeet(neetId string) bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	_, err := db.Exec("DELETE FROM ListToNeet WHERE list_id=? AND neet_id=?;", list.ListID, neetId)
	return err == nil
}

func DeleteNitterListByID(listId uint) bool {
	list := NitterListModel{ListID: listId}
	return list.Delete()
}
func (list NitterListModel) Delete() bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	_, err := db.Exec("DELETE FROM ListToNeet WHERE list_id=?", list.ListID)
	if err != nil {
		return false
	}

	_, err = db.Exec("DELETE FROM NitterLists WHERE list_id=?", list.ListID)
	return err == nil
}

// **NOT IMPLEMENTED**
func (list NitterListModel) Update() {
}
