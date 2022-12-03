package db

import (
	"errors"
	ps "teniditter-server/cmd/planetscale"
)

// Query Subteddit by its name from DB; if the subteddit is not yet in the db this function will insert it.
func GetSubteddit(subname string, depth ...int) (*SubtedditModel, error) {
	if len(depth) == 1 && depth[0] > 3 {
		return nil, errors.New("recursion emergency stop triggered")
	}

	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	var result SubtedditModel
	err := db.QueryRow("SELECT * FROM Subteddits WHERE subname=?", subname).Scan(&result.SubID, &result.Subname)
	if err != nil || result.Subname != subname {
		// Not in db --> insert it
		if ok := AddSubteddit(subname); ok {
			var depthVal int
			if len(depth) == 1 {
				depthVal = depth[0]
			}

			return GetSubteddit(subname, depthVal+1) // refetch by controlled recursion
		}

		return nil, err
	}

	return &result, nil
}

func AddSubteddit(subname string) bool {
	db := ps.DBManager.Connect()
	if db == nil {
		return false
	}

	_, err := db.Exec("INSERT INTO Subteddits (subname) VALUES (?);", subname)
	return err == nil
}

func SearchSubteddit(subname string) ([]SubtedditModel, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	rows, err := db.Query("SELECT * FROM Subteddits WHERE subname LIKE ?", "%"+subname+"%")
	if err != nil {
		return nil, errors.New("error when fetching Subteddits")
	}
	defer rows.Close()

	var AllSubs []SubtedditModel
	for rows.Next() {
		var sub SubtedditModel
		if err := rows.Scan(&sub.SubID, &sub.Subname); err != nil {
			return AllSubs, nil
		}
		AllSubs = append(AllSubs, sub)
	}
	if err = rows.Err(); err != nil {
		return AllSubs, err
	}

	return AllSubs, nil
}
