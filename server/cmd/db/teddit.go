package db

import (
	"errors"
)

// Query Subteddit by its name from DB; if the subteddit is not yet in the db this function will insert it.
func GetSubteddit(subname string, depth ...int) (*SubtedditModel, error) {
	if len(depth) == 1 && depth[0] > 3 {
		return nil, errors.New("recursion emergency stop triggered")
	}

	db := DBManager.Connect()
	if db == nil {
		return nil, ErrDbNotFound
	}

	var result SubtedditModel
	err := db.QueryRow("SELECT * FROM Subteddits WHERE subname=?", subname).Scan(&result.SubID, &result.Subname)
	if err != nil || result.Subname != subname {
		// Not in db --> insert it
		if _, err = db.Exec("INSERT INTO Subteddits (subname) VALUES (?);", subname); err == nil {
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
